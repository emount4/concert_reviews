package user_postgres_repository

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

type userRecord struct {
	// Поля таблицы users
	ID               uuid.UUID  `db:"user_id"`
	Email            string     `db:"email"`
	PasswordHash     string     `db:"password_hash"`
	TelegramID       *int64     `db:"tg_id"`
	TelegramUsername *string    `db:"tg_username"`
	RoleID           int        `db:"role_id"`
	Username         string     `db:"username"`
	Bio              *string    `db:"bio"`
	AvatarURL        *string    `db:"avatar_url"`
	BannerURL        *string    `db:"banner_url"`
	IsEmailVerified  bool       `db:"is_email_verified"`
	IsActive         bool       `db:"is_active"`
	IsBanned         bool       `db:"is_banned"`
	BannedByUserID   *uuid.UUID `db:"banned_by_user_id"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`

	// Поля таблицы user_stats (через LEFT JOIN)
	ReviewsCount       int        `db:"reviews_count"`
	LikesGivenCount    int        `db:"likes_given_count"`
	LikesReceivedCount int        `db:"likes_received_count"`
	StatsUpdatedAt     *time.Time `db:"stats_updated_at"`
}

func (r userRecord) MapToDomain() domain.User {
	u := domain.User{
		ID:               r.ID,
		Email:            r.Email,
		PasswordHash:     r.PasswordHash,
		TelegramID:       r.TelegramID,
		TelegramUsername: r.TelegramUsername,
		RoleID:           r.RoleID,
		Username:         r.Username,
		Bio:              r.Bio,
		AvatarURL:        r.AvatarURL,
		BannerURL:        r.BannerURL,
		IsEmailVerified:  r.IsEmailVerified,
		IsActive:         r.IsActive,
		IsBanned:         r.IsBanned,
		BannedByUserID:   r.BannedByUserID,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}

	// Всегда инициализируем структуру Stats, так как мы использовали COALESCE в SQL
	u.Stats = &domain.UserStats{
		ReviewsCount:       r.ReviewsCount,
		LikesGivenCount:    r.LikesGivenCount,
		LikesReceivedCount: r.LikesReceivedCount,
	}
	if r.StatsUpdatedAt != nil {
		u.Stats.UpdatedAt = *r.StatsUpdatedAt
	}

	return u
}
