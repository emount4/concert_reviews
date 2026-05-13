package venue_transport_http

import (
	"context"
	"net/http"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
)

type Service interface {
	CreateVenue(ctx context.Context, venue domain.Venue) (domain.Venue, error)
	GetVenueByID(ctx context.Context, id int) (domain.Venue, error)
	GetVenues(ctx context.Context, cityID *int, search string, sort string, direction string, capacityFrom, capacityTo *int, limit, offset *int) ([]domain.Venue, int, error)
	UpdateVenue(ctx context.Context, id int, patch domain.VenuePatch) (domain.Venue, error)
	DeleteVenueHard(ctx context.Context, id int) error
	DeleteVenueSoft(ctx context.Context, id int) error
	RestoreVenue(ctx context.Context, id int) (domain.Venue, error)
	GetVenuesAdmin(ctx context.Context, cityID *int, search string, sort string, direction string, capacityFrom, capacityTo *int, limit, offset *int, includeDeleted bool, status string) ([]domain.Venue, int, error)
}

type VenueHTTPHandler struct {
	venueService Service
}

func NewVenueHTTPHandler(venueService Service) *VenueHTTPHandler {
	return &VenueHTTPHandler{
		venueService: venueService,
	}
}

func (h *VenueHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/venues",
			Handler: h.GetVenues,
		},
		{
			Method:  http.MethodGet,
			Path:    "/venues/{id}",
			Handler: h.GetVenue,
		},
		{
			Method:  http.MethodPost,
			Path:    "/venues",
			Handler: h.CreateVenue,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/venues/{id}",
			Handler: h.UpdateVenue,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/venues-hard/{id}",
			Handler: h.DeleteVenueHard,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/venues-soft/{id}",
			Handler: h.DeleteVenueSoft,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodPost,
			Path:    "/venues/{id}/restore",
			Handler: h.RestoreVenue,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/admin/venues",
			Handler: h.GetVenuesAdmin,
			Access:  core_http_server.AccessAdminOnly,
		},
	}
}
