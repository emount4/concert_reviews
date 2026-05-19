package favorites_postgres_repository

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

type favoriteRecord struct {
	FavoriteID int
	UserID     uuid.UUID
	TargetType string
	TargetID   string
	TargetName string
	ImageURL   *string
	CreatedAt  time.Time
}

func (r favoriteRecord) MapToDomain() domain.Favorite {
	targetType, _ := domain.ParseFavoriteTargetType(r.TargetType)
	return domain.Favorite{
		FavoriteID: r.FavoriteID,
		UserID:     r.UserID,
		Target: domain.FavoriteTarget{
			Type:     targetType,
			ID:       r.TargetID,
			Name:     r.TargetName,
			ImageURL: r.ImageURL,
		},
		CreatedAt: r.CreatedAt,
	}
}
