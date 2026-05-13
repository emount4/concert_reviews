package review_postgres_repository

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
)

func (r *ReviewRepository) CreateReview(
	ctx context.Context,
	rev domain.Review,
) (domain.Review, error) {
	err := r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		// 1. Вставляем основную запись рецензии
		// Считаем RatingTotal методом домена перед вызовом этого метода в сервисе
		queryReview := `
			INSERT INTO reviews (
				user_id, concert_id, title, text, 
				p1, p2, p3, p4, p5, rating_total, 
				status, is_visible
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			RETURNING review_id, created_at
		`

		err := exec.QueryRow(txCtx, queryReview,
			rev.UserID, rev.ConcertID, rev.Title, rev.Text,
			rev.P1, rev.P2, rev.P3, rev.P4, rev.P5, rev.RatingTotal,
			rev.Status, rev.IsVisible,
		).Scan(&rev.ReviewID, &rev.CreatedAt)

		if err != nil {
			return fmt.Errorf("insert review: %w", err)
		}

		if len(rev.Media) > 0 {
			queryMedia := `
				INSERT INTO review_media (review_id, media_url, media_type, file_size, status)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING media_id, created_at
			`

			for i := range rev.Media {
				m := &rev.Media[i]
				err := exec.QueryRow(txCtx, queryMedia,
					rev.ReviewID, m.MediaURL, m.MediaType, m.FileSize, domain.StatusPending,
				).Scan(&m.MediaID, &m.CreatedAt)

				if err != nil {
					return fmt.Errorf("insert review media at index %d: %w", i, err)
				}
				m.ReviewID = rev.ReviewID
			}
		}

		return nil
	})

	if err != nil {
		return domain.Review{}, err
	}

	return rev, nil
}
