package user_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const baseUserSelect = `
	SELECT 
		u.user_id, u.email, u.password_hash, u.tg_id, u.tg_username, u.role_id, 
		u.username, u.bio, u.avatar_url, u.banner_url, u.is_email_verified, 
		u.is_active, u.is_banned, u.banned_by_user_id, u.created_at, u.updated_at,
		COALESCE(s.reviews_count, 0) as reviews_count,
		COALESCE(s.likes_given_count, 0) as likes_given_count,
		COALESCE(s.likes_received_count, 0) as likes_received_count,
		s.updated_at as stats_updated_at
	FROM users u
	LEFT JOIN user_stats s ON u.user_id = s.user_id
`

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := baseUserSelect + ` WHERE u.user_id = $1`

	var rec userRecord
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&rec.ID, &rec.Email, &rec.PasswordHash, &rec.TelegramID, &rec.TelegramUsername, &rec.RoleID,
		&rec.Username, &rec.Bio, &rec.AvatarURL, &rec.BannerURL, &rec.IsEmailVerified,
		&rec.IsActive, &rec.IsBanned, &rec.BannedByUserID, &rec.CreatedAt, &rec.UpdatedAt,
		&rec.ReviewsCount, &rec.LikesGivenCount, &rec.LikesReceivedCount, &rec.StatsUpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, core_errors.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("query user by id: %w", err)
	}

	return rec.MapToDomain(), nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := baseUserSelect + ` WHERE u.username = $1`

	var rec userRecord
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&rec.ID, &rec.Email, &rec.PasswordHash, &rec.TelegramID, &rec.TelegramUsername, &rec.RoleID,
		&rec.Username, &rec.Bio, &rec.AvatarURL, &rec.BannerURL, &rec.IsEmailVerified,
		&rec.IsActive, &rec.IsBanned, &rec.BannedByUserID, &rec.CreatedAt, &rec.UpdatedAt,
		&rec.ReviewsCount, &rec.LikesGivenCount, &rec.LikesReceivedCount, &rec.StatsUpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, core_errors.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("query user by username: %w", err)
	}

	return rec.MapToDomain(), nil
}
