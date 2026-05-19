package favorites_transport_http

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
)

func MapDomainToResponse(f domain.Favorite) FavoriteResponse {
	return FavoriteResponse{
		ID:         f.FavoriteID,
		TargetType: string(f.Target.Type),
		TargetID:   f.Target.ID,
		Name:       f.Target.Name,
		ImageURL:   f.Target.ImageURL,
		CreatedAt:  f.CreatedAt.Format(time.RFC3339),
	}
}

func MapDomainListToResponse(favorites []domain.Favorite) ListFavoritesResponse {
	items := make([]FavoriteResponse, len(favorites))
	for i, favorite := range favorites {
		items[i] = MapDomainToResponse(favorite)
	}
	return ListFavoritesResponse{Items: items}
}
