package review_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const baseReviewSelect = `
	SELECT 
		r.*,
		u.username as author_name, u.avatar_url as author_avatar,
		c.title as concert_title,
		c.poster_url as concert_poster_url,
		(SELECT COUNT(*) FROM review_likes WHERE review_id = r.review_id) as likes_count,
		COUNT(*) OVER() as total_count,
		(
			SELECT jsonb_agg(jsonb_build_object(
				'MediaID', rm.media_id,
				'MediaURL', rm.media_url,
				'MediaType', rm.media_type,
				'Status', rm.status
			))
			FROM review_media rm
			WHERE rm.review_id = r.review_id AND rm.status = 'approved'
		) as media_json,
		(
			SELECT COALESCE(jsonb_agg(jsonb_build_object(
				'ArtistID', a.artist_id,
				'Name', a.name,
				'IsMain', ca.is_main
			)), '[]'::jsonb)
			FROM concert_artists ca
			JOIN artists a ON a.artist_id = ca.artist_id
			WHERE ca.concert_id = r.concert_id AND ca.is_main = true
		) as concert_artists_json
	FROM reviews r
	JOIN users u ON r.user_id = u.user_id
	JOIN concerts c ON r.concert_id = c.concert_id
`

func (r *ReviewRepository) GetReviews(
	ctx context.Context,
	userID *uuid.UUID,
	concertID *uuid.UUID,
	artistID *int,
	venueID *int,
	sort string,
	direction string,
	limit, offset *int,
) ([]domain.Review, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var conditions []string
	args := []any{}
	argIndex := 1 // Для отслеживания номера параметра

	// Добавляем userID для проверки is_liked_by_me (может быть nil)
	if userID != nil {
		args = append(args, *userID)
		argIndex++
	}

	// Только одобренные и видимые
	conditions = append(conditions, "r.status = 'approved'", "r.is_visible = true", "r.deleted_at IS NULL")

	if concertID != nil {
		args = append(args, *concertID)
		conditions = append(conditions, fmt.Sprintf("r.concert_id = $%d", argIndex))
		argIndex++
	}
	if venueID != nil {
		args = append(args, *venueID)
		conditions = append(conditions, fmt.Sprintf("c.venue_id = $%d", argIndex))
		argIndex++
	}
	if artistID != nil {
		args = append(args, *artistID)
		conditions = append(conditions, fmt.Sprintf("r.concert_id IN (SELECT concert_id FROM concert_artists WHERE artist_id = $%d)", argIndex))
		argIndex++
	}

	// Сортировка
	sortMap := map[string]string{"date": "r.created_at", "rating": "r.rating_total", "likes": "likes_count", "count": "likes_count"}
	orderBy, ok := sortMap[sort]
	if !ok {
		orderBy = "r.created_at"
	}

	// Динамически добавляем проверку is_liked_by_me в SELECT (вставляем перед FROM)
	isLikedByMeExpr := "FALSE as is_liked_by_me" // по умолчанию false
	if userID != nil {
		isLikedByMeExpr = fmt.Sprintf("EXISTS(SELECT 1 FROM review_likes WHERE review_id = r.review_id AND user_id = $1) as is_liked_by_me")
	}

	// Вставляем выражение в конец списка SELECT перед 'FROM reviews r'
	selectWithLike := strings.Replace(baseReviewSelect, "\n\tFROM reviews r", ", "+isLikedByMeExpr+"\n\tFROM reviews r", 1)

	query := selectWithLike + " WHERE " + strings.Join(conditions, " AND ")
	query += fmt.Sprintf(" ORDER BY %s %s", orderBy, direction)

	if limit != nil && offset != nil {
		args = append(args, *limit, *offset)
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query reviews: %w", err)
	}
	defer rows.Close()

	var reviews []domain.Review
	total := 0
	for rows.Next() {
		var rec reviewRecord
		err := rows.Scan(
			&rec.ReviewID, &rec.UserID, &rec.ConcertID, &rec.Title, &rec.Text, &rec.OriginalText,
			&rec.P1, &rec.P2, &rec.P3, &rec.P4, &rec.P5, &rec.RatingTotal,
			&rec.Status, &rec.RejectionReason, &rec.ModeratedByUserID, &rec.IsVisible,
			&rec.CreatedAt, &rec.DeletedAt,
			&rec.AuthorName, &rec.AuthorAvatar,
			&rec.ConcertTitle, &rec.ConcertPosterURL, &rec.LikesCount, &rec.TotalCount, &rec.MediaJSON, &rec.ConcertArtistsJSON,
			&rec.IsLikedByMe,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan review: %w", err)
		}

		reviews = append(reviews, rec.MapToDomain())
		total = rec.TotalCount
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate reviews rows: %w", err)
	}
	return reviews, total, nil
}

func (r *ReviewRepository) GetPendingReviews(
	ctx context.Context,
	limit, offset *int,
) ([]domain.Review, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	const adminPendingSelect = `
		SELECT 
			r.*,
			u.username as author_name, u.avatar_url as author_avatar,
			c.title as concert_title,
			c.poster_url as concert_poster_url,
			(SELECT COUNT(*) FROM review_likes WHERE review_id = r.review_id) as likes_count,
			COUNT(*) OVER() as total_count,
			(
				SELECT jsonb_agg(jsonb_build_object(
					'MediaID', rm.media_id,
					'MediaURL', rm.media_url,
					'MediaType', rm.media_type,
					'Status', rm.status
				))
				FROM review_media rm
				WHERE rm.review_id = r.review_id -- Убрали фильтр по status = 'approved'
			) as media_json,
			(
				SELECT COALESCE(jsonb_agg(jsonb_build_object(
					'ArtistID', a.artist_id,
					'Name', a.name,
					'IsMain', ca.is_main
				)), '[]'::jsonb)
				FROM concert_artists ca
				JOIN artists a ON a.artist_id = ca.artist_id
				WHERE ca.concert_id = r.concert_id AND ca.is_main = true
			) as concert_artists_json
		FROM reviews r
		JOIN users u ON r.user_id = u.user_id
		JOIN concerts c ON r.concert_id = c.concert_id
		WHERE r.status = 'pending' AND r.deleted_at IS NULL
		ORDER BY r.created_at ASC -- Сначала самые старые (очередь)
	`

	var args []any
	query := adminPendingSelect

	if limit != nil && offset != nil {
		query += " LIMIT $1 OFFSET $2"
		args = append(args, *limit, *offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query pending reviews: %w", err)
	}
	defer rows.Close()

	var reviews []domain.Review
	total := 0

	for rows.Next() {
		var rec reviewRecord
		err := rows.Scan(
			&rec.ReviewID, &rec.UserID, &rec.ConcertID, &rec.Title, &rec.Text, &rec.OriginalText,
			&rec.P1, &rec.P2, &rec.P3, &rec.P4, &rec.P5, &rec.RatingTotal,
			&rec.Status, &rec.RejectionReason, &rec.ModeratedByUserID, &rec.IsVisible,
			&rec.CreatedAt, &rec.DeletedAt,
			&rec.AuthorName, &rec.AuthorAvatar,
			&rec.ConcertTitle, &rec.ConcertPosterURL, &rec.LikesCount, &rec.TotalCount, &rec.MediaJSON, &rec.ConcertArtistsJSON,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan pending review: %w", err)
		}

		reviews = append(reviews, rec.MapToDomain())
		total = rec.TotalCount
	}

	return reviews, total, nil
}

func (r *ReviewRepository) GetUserReviewCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	query := `SELECT reviews_count FROM user_stats WHERE user_id = $1`
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func (r *ReviewRepository) GetReviewByID(ctx context.Context, id uuid.UUID) (domain.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	// Используем агрегацию jsonb_agg, чтобы вытащить все медиа сразу.
	// COALESCE нужен, чтобы если медиа нет, вернулся пустой массив '[]', а не NULL.
	query := `
		SELECT 
			r.review_id, r.user_id, r.concert_id, r.title, r.text, r.original_text,
			r.p1, r.p2, r.p3, r.p4, r.p5, r.rating_total,
			r.status, r.rejection_reason, r.moderated_by_user_id, r.is_visible,
			r.created_at, r.deleted_at,
			u.username as author_name, u.avatar_url as author_avatar,
			c.title as concert_title,
			c.poster_url as concert_poster_url,
			(SELECT COUNT(*) FROM review_likes WHERE review_id = r.review_id) as likes_count,
			(
				SELECT COALESCE(jsonb_agg(jsonb_build_object(
					'MediaID', rm.media_id,
					'ReviewID', rm.review_id,
					'MediaURL', rm.media_url,
					'MediaType', rm.media_type,
					'FileSize', rm.file_size,
					'Status', rm.status,
					'CreatedAt', rm.created_at
				)), '[]'::jsonb)
				FROM review_media rm
				WHERE rm.review_id = r.review_id
			) as media_json,
			(
				SELECT COALESCE(jsonb_agg(jsonb_build_object(
					'ArtistID', a.artist_id,
					'Name', a.name,
					'IsMain', ca.is_main
				)), '[]'::jsonb)
				FROM concert_artists ca
				JOIN artists a ON a.artist_id = ca.artist_id
				WHERE ca.concert_id = r.concert_id AND ca.is_main = true
			) as concert_artists_json
		FROM reviews r
		JOIN users u ON r.user_id = u.user_id
		JOIN concerts c ON r.concert_id = c.concert_id
		WHERE r.review_id = $1 AND r.deleted_at IS NULL
	`

	var rec reviewRecord // Используем структуру из твоего файла model.go
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&rec.ReviewID, &rec.UserID, &rec.ConcertID, &rec.Title, &rec.Text, &rec.OriginalText,
		&rec.P1, &rec.P2, &rec.P3, &rec.P4, &rec.P5, &rec.RatingTotal,
		&rec.Status, &rec.RejectionReason, &rec.ModeratedByUserID, &rec.IsVisible,
		&rec.CreatedAt, &rec.DeletedAt,
		&rec.AuthorName, &rec.AuthorAvatar,
		&rec.ConcertTitle, &rec.ConcertPosterURL, &rec.LikesCount,
		&rec.MediaJSON, &rec.ConcertArtistsJSON,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Review{}, core_errors.ErrNotFound
		}
		return domain.Review{}, fmt.Errorf("query review by id: %w", err)
	}

	// Вызываем твой маппер из record в domain
	return rec.MapToDomain(), nil
}

func (r *ReviewRepository) GetUserReviews(ctx context.Context, userID uuid.UUID, includeStatuses []string) ([]domain.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if len(includeStatuses) == 0 {
		includeStatuses = []string{"approved"}
	}

	// Строим список статусов для фильтра
	statusCondition := "r.status = ANY($2)"

	query := `
		SELECT 
			r.*,
			u.username as author_name, u.avatar_url as author_avatar,
			c.title as concert_title,
			c.poster_url as concert_poster_url,
			(SELECT COUNT(*) FROM review_likes WHERE review_id = r.review_id) as likes_count,
			(
				SELECT jsonb_agg(jsonb_build_object(
					'MediaID', rm.media_id,
					'MediaURL', rm.media_url,
					'MediaType', rm.media_type,
					'Status', rm.status
				))
				FROM review_media rm
				WHERE rm.review_id = r.review_id AND rm.status = 'approved'
			) as media_json,
			(
				SELECT COALESCE(jsonb_agg(jsonb_build_object(
					'ArtistID', a.artist_id,
					'Name', a.name,
					'IsMain', ca.is_main
				)), '[]'::jsonb)
				FROM concert_artists ca
				JOIN artists a ON a.artist_id = ca.artist_id
				WHERE ca.concert_id = r.concert_id AND ca.is_main = true
			) as concert_artists_json
		FROM reviews r
		JOIN users u ON r.user_id = u.user_id
		JOIN concerts c ON r.concert_id = c.concert_id
		WHERE r.user_id = $1 AND ` + statusCondition + ` AND r.deleted_at IS NULL
		ORDER BY r.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID, includeStatuses)
	if err != nil {
		return nil, fmt.Errorf("query user reviews: %w", err)
	}
	defer rows.Close()

	var reviews []domain.Review
	for rows.Next() {
		var rec reviewRecord
		err := rows.Scan(
			&rec.ReviewID, &rec.UserID, &rec.ConcertID, &rec.Title, &rec.Text, &rec.OriginalText,
			&rec.P1, &rec.P2, &rec.P3, &rec.P4, &rec.P5, &rec.RatingTotal,
			&rec.Status, &rec.RejectionReason, &rec.ModeratedByUserID, &rec.IsVisible,
			&rec.CreatedAt, &rec.DeletedAt,
			&rec.AuthorName, &rec.AuthorAvatar,
			&rec.ConcertTitle, &rec.ConcertPosterURL, &rec.LikesCount,
			&rec.MediaJSON, &rec.ConcertArtistsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user review: %w", err)
		}

		reviews = append(reviews, rec.MapToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user reviews rows: %w", err)
	}

	return reviews, nil
}
