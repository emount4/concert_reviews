package favorites_service

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

const MaxFavoritesPerType = 5

type FavoritesService struct {
	favoritesRepository FavoritesRepository
}

type FavoritesRepository interface {
	AddFavorite(ctx context.Context, userID uuid.UUID, target domain.FavoriteTarget) (domain.Favorite, error)
	DeleteFavorite(ctx context.Context, userID uuid.UUID, target domain.FavoriteTarget) error
	CountFavoritesByType(ctx context.Context, userID uuid.UUID, targetType domain.FavoriteTargetType) (int, error)
	GetFavoriteTarget(ctx context.Context, targetType domain.FavoriteTargetType, targetID string) (domain.FavoriteTarget, error)
	GetFavoritesByUsername(ctx context.Context, username string, targetType *domain.FavoriteTargetType) ([]domain.Favorite, error)
}

func NewFavoritesService(favoritesRepository FavoritesRepository) *FavoritesService {
	return &FavoritesService{
		favoritesRepository: favoritesRepository,
	}
}
