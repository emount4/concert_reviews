package domain

import (
	"fmt"
	"time"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

type FavoriteTargetType string

const (
	FavoriteTargetArtist  FavoriteTargetType = "artist"
	FavoriteTargetVenue   FavoriteTargetType = "venue"
	FavoriteTargetConcert FavoriteTargetType = "concert"
)

type FavoriteTarget struct {
	Type     FavoriteTargetType
	ID       string
	Name     string
	ImageURL *string
}

type Favorite struct {
	FavoriteID int
	UserID     uuid.UUID
	Target     FavoriteTarget
	CreatedAt  time.Time
}

func ParseFavoriteTargetType(value string) (FavoriteTargetType, error) {
	targetType := FavoriteTargetType(value)
	switch targetType {
	case FavoriteTargetArtist, FavoriteTargetVenue, FavoriteTargetConcert:
		return targetType, nil
	default:
		return "", fmt.Errorf("unsupported favorite target type %q: %w", value, core_errors.ErrInvalidArgument)
	}
}
