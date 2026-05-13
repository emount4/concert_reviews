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
	GetArtists(ctx context.Context, search string, sort string, direction string, hasReviews *bool, limit, offset *int) ([]domain.Artist, int, error)
	PatchArtist(ctx context.Context, id int, patch domain.ArtistPatch) (domain.Artist, error)
	DeleteArtistHard(ctx context.Context, id int) error
	DeleteArtistSoft(ctx context.Context, id int) error
	RestoreArtist(ctx context.Context, id int) (domain.Artist, error)
	GetArtistsAdmin(
		ctx context.Context,
		search string,
		sort string,
		direction string,
		hasReviews *bool,
		limit, offset *int,
		includeDeleted bool,
		status string,
	) ([]domain.Artist, int, error)
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
		{
			Method:  http.MethodPatch,
			Path:    "/artists/{id}",
			Handler: h.PatchArtist,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/artists-hard/{id}",
			Handler: h.DeleteArtistHard,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/artists-soft/{id}",
			Handler: h.DeleteArtistSoft,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodPost,
			Path:    "/artists/{id}/restore",
			Handler: h.RestoreArtist,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/admin/artists",
			Handler: h.GetArtistsAdmin,
			Access:  core_http_server.AccessAdminOnly,
		},
	}
}
