package review_postgres_repository

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
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
			    moderated_by_user_id = $4
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
			queryApproveMedia := `UPDATE review_media SET status = 'approved' WHERE media_id = ANY($1)`
			if _, err := exec.Exec(txCtx, queryApproveMedia, allowedMediaIDs); err != nil {
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
				sum_rating_total = sum_rating_total + $1,
				reviews_count = reviews_count + 1,
				updated_at = NOW()
			WHERE artist_id IN (
				SELECT artist_id FROM concert_artists 
				WHERE concert_id = $2 AND is_main = TRUE
			)
		`
		if _, err := exec.Exec(txCtx, queryArtists, rev.RatingTotal, rev.ConcertID); err != nil {
			return fmt.Errorf("update artists stats: %w", err)
		}

		// 5. СТАТИСТИКА ПЛОЩАДКИ
		queryVenue := `
			UPDATE venue_stats SET 
				sum_rating_total = sum_rating_total + $1,
				reviews_count = reviews_count + 1,
				updated_at = NOW()
			WHERE venue_id = (SELECT venue_id FROM concerts WHERE concert_id = $2)
		`
		if _, err := exec.Exec(txCtx, queryVenue, rev.RatingTotal, rev.ConcertID); err != nil {
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
	query := `
		UPDATE reviews 
		SET status = 'rejected', 
		    rejection_reason = $1, 
		    moderated_by_user_id = $2,
            is_visible = false,
		    updated_at = NOW()
		WHERE review_id = $3
	`
	_, err := r.pool.Exec(ctx, query, reason, moderatorID, id)
	return err
}
