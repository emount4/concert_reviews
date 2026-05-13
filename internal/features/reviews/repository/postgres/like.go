package review_postgres_repository

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *ReviewRepository) CreateLike(ctx context.Context, reviewID, userID uuid.UUID) error {
	return r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		if _, err := exec.Exec(txCtx, `INSERT INTO review_likes (user_id, review_id) VALUES ($1, $2)`, userID, reviewID); err != nil {
			return err
		}

		if _, err := exec.Exec(txCtx, `UPDATE user_stats SET likes_given_count = likes_given_count + 1 WHERE user_id = $1`, userID); err != nil {
			return err
		}

		queryReceived := `
			UPDATE user_stats SET likes_received_count = likes_received_count + 1 
			WHERE user_id = (SELECT user_id FROM reviews WHERE review_id = $1)
		`
		if _, err := exec.Exec(txCtx, queryReceived, reviewID); err != nil {
			return err
		}

		return nil
	})
}

func (r *ReviewRepository) DeleteLike(ctx context.Context, reviewID, userID uuid.UUID) error {
	return r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		if _, err := exec.Exec(txCtx, `DELETE FROM review_likes WHERE user_id = $1 AND review_id = $2`, userID, reviewID); err != nil {
			return err
		}

		if _, err := exec.Exec(txCtx, `UPDATE user_stats SET likes_given_count = GREATEST(0, likes_given_count - 1) WHERE user_id = $1`, userID); err != nil {
			return err
		}

		queryReceived := `
			UPDATE user_stats SET likes_received_count = GREATEST(0, likes_received_count - 1) 
			WHERE user_id = (SELECT user_id FROM reviews WHERE review_id = $1)
		`
		if _, err := exec.Exec(txCtx, queryReceived, reviewID); err != nil {
			return err
		}

		return nil
	})
}

func (r *ReviewRepository) GetLike(ctx context.Context, reviewID, userID uuid.UUID) (domain.ReviewLike, error) {
	var l domain.ReviewLike
	query := `SELECT like_id, user_id, review_id, created_at FROM review_likes WHERE user_id = $1 AND review_id = $2`
	err := r.pool.QueryRow(ctx, query, userID, reviewID).Scan(&l.LikeID, &l.UserID, &l.ReviewID, &l.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ReviewLike{}, core_errors.ErrNotFound
		}
		return domain.ReviewLike{}, err
	}
	return l, nil
}

func (r *ReviewRepository) GetLikers(ctx context.Context, reviewID uuid.UUID) ([]domain.User, error) {
	query := `
		SELECT u.user_id, u.username, u.avatar_url 
		FROM users u
		JOIN review_likes rl ON u.user_id = rl.user_id
		WHERE rl.review_id = $1
		ORDER BY rl.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, reviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate likers rows: %w", err)
	}

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
