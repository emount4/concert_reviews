package favorites_service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *FavoritesService) AddFavorite(
	ctx context.Context,
	userID uuid.UUID,
	targetTypeValue string,
	targetID string,
) (domain.Favorite, error) {
	if userID == uuid.Nil {
		return domain.Favorite{}, fmt.Errorf("user id is empty: %w", core_errors.ErrUnauthorized)
	}

	targetType, targetID, err := parseFavoriteTarget(targetTypeValue, targetID)
	if err != nil {
		return domain.Favorite{}, err
	}

	target, err := s.favoritesRepository.GetFavoriteTarget(ctx, targetType, targetID)
	if err != nil {
		return domain.Favorite{}, fmt.Errorf("fetch favorite target: %w", err)
	}

	count, err := s.favoritesRepository.CountFavoritesByType(ctx, userID, targetType)
	if err != nil {
		return domain.Favorite{}, fmt.Errorf("count favorites: %w", err)
	}
	if count >= MaxFavoritesPerType {
		return domain.Favorite{}, fmt.Errorf("max %d favorites of type %s exceeded: %w", MaxFavoritesPerType, targetType, core_errors.ErrConflict)
	}

	favorite, err := s.favoritesRepository.AddFavorite(ctx, userID, target)
	if err != nil {
		return domain.Favorite{}, fmt.Errorf("save favorite: %w", err)
	}

	return favorite, nil
}

func (s *FavoritesService) DeleteFavorite(
	ctx context.Context,
	userID uuid.UUID,
	targetTypeValue string,
	targetID string,
) error {
	if userID == uuid.Nil {
		return fmt.Errorf("user id is empty: %w", core_errors.ErrUnauthorized)
	}

	targetType, targetID, err := parseFavoriteTarget(targetTypeValue, targetID)
	if err != nil {
		return err
	}

	if err := s.favoritesRepository.DeleteFavorite(ctx, userID, domain.FavoriteTarget{Type: targetType, ID: targetID}); err != nil {
		return fmt.Errorf("delete favorite: %w", err)
	}

	return nil
}

func (s *FavoritesService) GetFavoritesByUsername(
	ctx context.Context,
	username string,
	targetTypeValue *string,
) ([]domain.Favorite, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("username is required: %w", core_errors.ErrInvalidArgument)
	}

	var targetType *domain.FavoriteTargetType
	if targetTypeValue != nil && strings.TrimSpace(*targetTypeValue) != "" {
		parsedType, err := domain.ParseFavoriteTargetType(strings.TrimSpace(*targetTypeValue))
		if err != nil {
			return nil, err
		}
		targetType = &parsedType
	}

	favorites, err := s.favoritesRepository.GetFavoritesByUsername(ctx, username, targetType)
	if err != nil {
		return nil, fmt.Errorf("fetch favorites by username: %w", err)
	}

	return favorites, nil
}

func parseFavoriteTarget(targetTypeValue string, targetIDValue string) (domain.FavoriteTargetType, string, error) {
	targetType, err := domain.ParseFavoriteTargetType(strings.TrimSpace(targetTypeValue))
	if err != nil {
		return "", "", err
	}

	targetID := strings.TrimSpace(targetIDValue)
	if targetID == "" {
		return "", "", fmt.Errorf("target id is required: %w", core_errors.ErrInvalidArgument)
	}
	switch targetType {
	case domain.FavoriteTargetArtist:
		if _, err := strconv.Atoi(targetID); err != nil {
			return "", "", fmt.Errorf("artist id must be integer: %w", core_errors.ErrInvalidArgument)
		}
	case domain.FavoriteTargetVenue:
		if _, err := strconv.Atoi(targetID); err != nil {
			return "", "", fmt.Errorf("venue id must be integer: %w", core_errors.ErrInvalidArgument)
		}
	case domain.FavoriteTargetConcert:
		if _, err := uuid.Parse(targetID); err != nil {
			return "", "", fmt.Errorf("concert id must be uuid: %w", core_errors.ErrInvalidArgument)
		}
	}

	return targetType, targetID, nil
}
