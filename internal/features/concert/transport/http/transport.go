package concert_transport_http

import (
	"context"
	"net/http"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type Service interface {
	// Основные концерты
	CreateConcert(ctx context.Context, concert domain.Concert, artists []domain.ConcertArtist) (domain.Concert, error)
	GetConcertByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (domain.Concert, error)
	GetConcerts(ctx context.Context, cityID *int, artistID *int, search string, sort string, direction string, limit, offset *int) ([]domain.Concert, int, error)
	UpdateConcert(ctx context.Context, id uuid.UUID, patch domain.ConcertPatch) (domain.Concert, error)
	DeleteConcertHard(ctx context.Context, id uuid.UUID) error
	DeleteConcertSoft(ctx context.Context, id uuid.UUID) error
	RestoreConcert(ctx context.Context, id uuid.UUID) (domain.Concert, error)
	GetConcertsAdmin(ctx context.Context, cityID *int, artistID *int, search string, sort string, direction string, limit, offset *int, includeDeleted bool) ([]domain.Concert, int, error)

	// Предложения (Suggestions)
	SuggestConcert(ctx context.Context, suggestion domain.ConcertSuggestion) (domain.ConcertSuggestion, error)
	GetSuggestionsAdmin(ctx context.Context, limit, offset *int, status string) ([]domain.ConcertSuggestion, error)
	GetSuggestionByIDAdmin(ctx context.Context, id uuid.UUID) (domain.ConcertSuggestion, error)
	DeleteSuggestionAdmin(ctx context.Context, id uuid.UUID) error

	// Управление артистами в концертах
	AddArtistToConcert(ctx context.Context, concertID uuid.UUID, artistID int, isMain bool) error
	RemoveArtistFromConcert(ctx context.Context, concertID uuid.UUID, artistID int) error
	UpdateArtistMainStatus(ctx context.Context, concertID uuid.UUID, artistID int, isMain bool) error
}

type ConcertHTTPHandler struct {
	concertService Service
}

func NewConcertHTTPHandler(concertService Service) *ConcertHTTPHandler {
	return &ConcertHTTPHandler{
		concertService: concertService,
	}
}

func (h *ConcertHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/concerts",
			Handler: h.GetConcerts,
		},
		{
			Method:  http.MethodGet,
			Path:    "/concerts/{id}",
			Handler: h.GetConcert,
		},
		{
			Method:  http.MethodPost,
			Path:    "/concerts/suggest",
			Handler: h.SuggestConcert,
			Access:  core_http_server.AccessAuthOnly,
		},

		{
			Method:  http.MethodPost,
			Path:    "/admin/concerts",
			Handler: h.CreateConcert,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/admin/concerts/{id}",
			Handler: h.UpdateConcert,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/admin/concerts-hard/{id}",
			Handler: h.DeleteConcertHard,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/admin/concerts-soft/{id}",
			Handler: h.DeleteConcertSoft,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodPost,
			Path:    "/admin/concerts/{id}/restore",
			Handler: h.RestoreConcert,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/admin/concerts",
			Handler: h.GetConcertsAdmin,
			Access:  core_http_server.AccessAdminOnly,
		},

		{
			Method:  http.MethodGet,
			Path:    "/admin/concerts/suggestions",
			Handler: h.GetSuggestionsAdmin,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/admin/concerts/suggestions/{id}",
			Handler: h.GetSuggestionByIDAdmin,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/admin/concerts/suggestions/{id}",
			Handler: h.DeleteSuggestionAdmin,
			Access:  core_http_server.AccessAdminOnly,
		},

		// Artist management routes
		{
			Method:  http.MethodPost,
			Path:    "/admin/concerts/{id}/artists",
			Handler: h.AddConcertArtist,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/admin/concerts/{id}/artists/{artist_id}",
			Handler: h.RemoveConcertArtist,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/admin/concerts/{id}/artists/{artist_id}",
			Handler: h.UpdateConcertArtistIsMain,
			Access:  core_http_server.AccessAdminOnly,
		},
	}
}
