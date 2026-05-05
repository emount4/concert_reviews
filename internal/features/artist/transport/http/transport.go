package artist_transport_http

import (
	"context"
	"net/http"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
)

type Service interface {
	CreateArtist(ctx context.Context, city domain.Artist) (domain.Artist, error)
	GetArtistByID(ctx context.Context, id int) (domain.Artist, error)
	GetArtists(ctx context.Context, search string, limit, offset *int) ([]domain.Artist, error)
}

type ArtistHTTPHandler struct {
	artistService Service
}

func NewArtistHTTPHandler(
	artistService Service,
) *ArtistHTTPHandler {
	return &ArtistHTTPHandler{
		artistService: artistService,
	}
}

func (h *ArtistHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/artists",
			Handler: h.GetArtists,
		},
		{
			Method:  http.MethodGet,
			Path:    "/artists/{id}",
			Handler: h.GetArtist,
		},
		{
			Method:  http.MethodPost,
			Path:    "/artists",
			Handler: h.Create,
			Access:  core_http_server.AccessAdminOnly,
		},
	}
}
