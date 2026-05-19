package favorites_transport_http

import (
	"context"
	"net/http"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type FavoritesService interface {
	AddFavorite(ctx context.Context, userID uuid.UUID, targetType string, targetID string) (domain.Favorite, error)
	DeleteFavorite(ctx context.Context, userID uuid.UUID, targetType string, targetID string) error
	GetFavoritesByUsername(ctx context.Context, username string, targetType *string) ([]domain.Favorite, error)
}

type FavoritesHTTPHandler struct {
	favoritesService FavoritesService
}

func NewFavoritesHTTPHandler(favoritesService FavoritesService) *FavoritesHTTPHandler {
	return &FavoritesHTTPHandler{
		favoritesService: favoritesService,
	}
}

func (h *FavoritesHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/favorites",
			Handler: h.AddFavorite,
			Access:  core_http_server.AccessAuthOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/favorites/{target_type}/{target_id}",
			Handler: h.DeleteFavorite,
			Access:  core_http_server.AccessAuthOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{username}/favorites",
			Handler: h.GetFavoritesByUsername,
		},
	}
}
