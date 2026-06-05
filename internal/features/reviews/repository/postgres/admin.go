package review_postgres_repository

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/google/uuid"
)

func (r *ReviewRepository) ApproveReview(
	ctx context.Context,
	id uuid.UUID,
	moderatorID uuid.UUID,
	finalTitle string,
	finalText string,
	allowedMediaIDs []uuid.UUID,
	rev domain.Review,
) error {
	// Используем твой WithinTx для обеспечения атомарности
	return r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		// 1. ОБНОВЛЯЕМ РЕЦЕНЗИЮ
		// Сохраняем правки модератора и переносим старый текст в историю
		queryReview := `
			UPDATE reviews 
			SET status = 'approved', 
			    title = $1, 
			    text = $2, 
			    original_text = $3, 
			    moderated_by_user_id = $4,
			    rejection_reason = NULL,
			    is_visible = true
			WHERE review_id = $5
		`
		if _, err := exec.Exec(txCtx, queryReview, finalTitle, finalText, rev.Text, moderatorID, id); err != nil {
			return fmt.Errorf("update review: %w", err)
		}

		// 2. МОДЕРАЦИЯ МЕДИА
		// Сначала сбрасываем всё в rejected
		if _, err := exec.Exec(txCtx, `UPDATE review_media SET status = 'rejected' WHERE review_id = $1`, id); err != nil {
			return fmt.Errorf("reset media status: %w", err)
		}
		// Одобряем только те, что прислал админ
		if len(allowedMediaIDs) > 0 {
			queryApproveMedia := `UPDATE review_media SET status = 'approved' WHERE review_id = $1 AND media_id = ANY($2)`
			if _, err := exec.Exec(txCtx, queryApproveMedia, id, allowedMediaIDs); err != nil {
				return fmt.Errorf("approve selected media: %w", err)
			}
		}

		// 3. СТАТИСТИКА КОНЦЕРТА
		queryConcert := `
			UPDATE concert_stats SET 
				sum_p1 = sum_p1 + $1, sum_p2 = sum_p2 + $2, sum_p3 = sum_p3 + $3, 
				sum_p4 = sum_p4 + $4, sum_p5 = sum_p5 + $5, 
				sum_rating_total = sum_rating_total + $6,
				reviews_count = reviews_count + 1,
				updated_at = NOW()
			WHERE concert_id = $7
		`
		if _, err := exec.Exec(txCtx, queryConcert, rev.P1, rev.P2, rev.P3, rev.P4, rev.P5, rev.RatingTotal, rev.ConcertID); err != nil {
			return fmt.Errorf("update concert stats: %w", err)
		}

		// 4. СТАТИСТИКА АРТИСТОВ (Всех основных артистов концерта)
		queryArtists := `
			UPDATE artist_stats SET 
				sum_p1 = sum_p1 + $1,
				sum_p2 = sum_p2 + $2,
				sum_p3 = sum_p3 + $3,
				sum_p4 = sum_p4 + $4,
				sum_p5 = sum_p5 + $5,
				sum_rating_total = sum_rating_total + $6,
				reviews_count = reviews_count + 1,
				updated_at = NOW()
			WHERE artist_id IN (
				SELECT artist_id FROM concert_artists 
				WHERE concert_id = $7 AND is_main = TRUE
			)
		`
		if _, err := exec.Exec(txCtx, queryArtists, rev.P1, rev.P2, rev.P3, rev.P4, rev.P5, rev.RatingTotal, rev.ConcertID); err != nil {
			return fmt.Errorf("update artists stats: %w", err)
		}

		// 5. СТАТИСТИКА ПЛОЩАДКИ
		queryVenue := `
			UPDATE venue_stats SET 
				sum_p1 = sum_p1 + $1,
				sum_p2 = sum_p2 + $2,
				sum_p3 = sum_p3 + $3,
				sum_p4 = sum_p4 + $4,
				sum_p5 = sum_p5 + $5,
				sum_rating_total = sum_rating_total + $6,
				reviews_count = reviews_count + 1,
				updated_at = NOW()
			WHERE venue_id = (SELECT venue_id FROM concerts WHERE concert_id = $7)
		`
		if _, err := exec.Exec(txCtx, queryVenue, rev.P1, rev.P2, rev.P3, rev.P4, rev.P5, rev.RatingTotal, rev.ConcertID); err != nil {
			return fmt.Errorf("update venue stats: %w", err)
		}

		// 6. СТАТИСТИКА ПОЛЬЗОВАТЕЛЯ (Автора рецензии)
		queryUser := `UPDATE user_stats SET reviews_count = reviews_count + 1 WHERE user_id = $1`
		if _, err := exec.Exec(txCtx, queryUser, rev.UserID); err != nil {
			return fmt.Errorf("update user stats: %w", err)
		}

		return nil
	})
}

func (r *ReviewRepository) RejectReview(ctx context.Context, id, moderatorID uuid.UUID, reason string) error {
	return r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		rev, err := lockReviewForModeration(txCtx, exec, id)
		if err != nil {
			return err
		}
		if rev.Status == domain.StatusRejected {
			return fmt.Errorf("review is already rejected: %w", core_errors.ErrConflict)
		}

		query := `
			UPDATE reviews
			SET status = 'rejected',
			    rejection_reason = $1,
			    moderated_by_user_id = $2,
			    is_visible = false
			WHERE review_id = $3
		`
		if _, err := exec.Exec(txCtx, query, reason, moderatorID, id); err != nil {
			return fmt.Errorf("reject review: %w", err)
		}

		if _, err := exec.Exec(txCtx, `UPDATE review_media SET status = 'rejected' WHERE review_id = $1`, id); err != nil {
			return fmt.Errorf("reject review media: %w", err)
		}

		if rev.Status == domain.StatusApproved {
			if err := subtractApprovedReviewStats(txCtx, exec, rev); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *ReviewRepository) ReturnReviewToPending(ctx context.Context, id, _ uuid.UUID) error {
	return r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		rev, err := lockReviewForModeration(txCtx, exec, id)
		if err != nil {
			return err
		}
		if rev.Status == domain.StatusPending {
			return fmt.Errorf("review is already pending: %w", core_errors.ErrConflict)
		}

		query := `
			UPDATE reviews
			SET status = 'pending',
			    rejection_reason = NULL,
			    moderated_by_user_id = NULL,
			    is_visible = true
			WHERE review_id = $1
		`
		if _, err := exec.Exec(txCtx, query, id); err != nil {
			return fmt.Errorf("return review to pending: %w", err)
		}

		if _, err := exec.Exec(txCtx, `UPDATE review_media SET status = 'pending' WHERE review_id = $1`, id); err != nil {
			return fmt.Errorf("return review media to pending: %w", err)
		}

		if rev.Status == domain.StatusApproved {
			if err := subtractApprovedReviewStats(txCtx, exec, rev); err != nil {
				return err
			}
		}

		return nil
	})
}

func lockReviewForModeration(
	ctx context.Context,
	exec core_postgres_tx.Executor,
	id uuid.UUID,
) (domain.Review, error) {
	var rev domain.Review
	var status string
	query := `
		SELECT user_id, concert_id, p1, p2, p3, p4, p5, rating_total, status
		FROM reviews
		WHERE review_id = $1 AND deleted_at IS NULL
		FOR UPDATE
	`
	if err := exec.QueryRow(ctx, query, id).Scan(
		&rev.UserID,
		&rev.ConcertID,
		&rev.P1,
		&rev.P2,
		&rev.P3,
		&rev.P4,
		&rev.P5,
		&rev.RatingTotal,
		&status,
	); err != nil {
		return domain.Review{}, fmt.Errorf("lock review for moderation: %w", err)
	}
	rev.ReviewID = id
	rev.Status = domain.ModerationStatus(status)
	return rev, nil
}

func subtractApprovedReviewStats(ctx context.Context, exec core_postgres_tx.Executor, rev domain.Review) error {
	queryConcert := `
		UPDATE concert_stats SET
			sum_p1 = GREATEST(0, COALESCE(sum_p1, 0) - $1),
			sum_p2 = GREATEST(0, COALESCE(sum_p2, 0) - $2),
			sum_p3 = GREATEST(0, COALESCE(sum_p3, 0) - $3),
			sum_p4 = GREATEST(0, COALESCE(sum_p4, 0) - $4),
			sum_p5 = GREATEST(0, COALESCE(sum_p5, 0) - $5),
			sum_rating_total = GREATEST(0, COALESCE(sum_rating_total, 0) - $6),
			reviews_count = GREATEST(0, COALESCE(reviews_count, 0) - 1),
			updated_at = NOW()
		WHERE concert_id = $7
	`
	if _, err := exec.Exec(ctx, queryConcert, rev.P1, rev.P2, rev.P3, rev.P4, rev.P5, rev.RatingTotal, rev.ConcertID); err != nil {
		return fmt.Errorf("subtract concert stats: %w", err)
	}

	queryArtists := `
		UPDATE artist_stats SET
			sum_p1 = GREATEST(0, COALESCE(sum_p1, 0) - $1),
			sum_p2 = GREATEST(0, COALESCE(sum_p2, 0) - $2),
			sum_p3 = GREATEST(0, COALESCE(sum_p3, 0) - $3),
			sum_p4 = GREATEST(0, COALESCE(sum_p4, 0) - $4),
			sum_p5 = GREATEST(0, COALESCE(sum_p5, 0) - $5),
			sum_rating_total = GREATEST(0, COALESCE(sum_rating_total, 0) - $6),
			reviews_count = GREATEST(0, COALESCE(reviews_count, 0) - 1),
			updated_at = NOW()
		WHERE artist_id IN (
			SELECT artist_id FROM concert_artists
			WHERE concert_id = $7 AND is_main = TRUE
		)
	`
	if _, err := exec.Exec(ctx, queryArtists, rev.P1, rev.P2, rev.P3, rev.P4, rev.P5, rev.RatingTotal, rev.ConcertID); err != nil {
		return fmt.Errorf("subtract artist stats: %w", err)
	}

	queryVenue := `
		UPDATE venue_stats SET
			sum_p1 = GREATEST(0, COALESCE(sum_p1, 0) - $1),
			sum_p2 = GREATEST(0, COALESCE(sum_p2, 0) - $2),
			sum_p3 = GREATEST(0, COALESCE(sum_p3, 0) - $3),
			sum_p4 = GREATEST(0, COALESCE(sum_p4, 0) - $4),
			sum_p5 = GREATEST(0, COALESCE(sum_p5, 0) - $5),
			sum_rating_total = GREATEST(0, COALESCE(sum_rating_total, 0) - $6),
			reviews_count = GREATEST(0, COALESCE(reviews_count, 0) - 1),
			updated_at = NOW()
		WHERE venue_id = (SELECT venue_id FROM concerts WHERE concert_id = $7)
	`
	if _, err := exec.Exec(ctx, queryVenue, rev.P1, rev.P2, rev.P3, rev.P4, rev.P5, rev.RatingTotal, rev.ConcertID); err != nil {
		return fmt.Errorf("subtract venue stats: %w", err)
	}

	queryUser := `
		UPDATE user_stats
		SET reviews_count = GREATEST(0, COALESCE(reviews_count, 0) - 1),
		    updated_at = NOW()
		WHERE user_id = $1
	`
	if _, err := exec.Exec(ctx, queryUser, rev.UserID); err != nil {
		return fmt.Errorf("subtract user stats: %w", err)
	}

	return nil
}
