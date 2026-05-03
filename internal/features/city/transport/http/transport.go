package city_transport_http

import (
	"context"
	"net/http"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
)

type Service interface {
	Create(ctx context.Context, city core_models.City) (core_models.City, error)
	GetCities(ctx context.Context, limit, offset *int) ([]core_models.City, error)
	GetCityByID(ctx context.Context, id int) (core_models.City, error)

	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, name, slug, tz *string) (core_models.City, error)
}

type CityHTTPHandler struct {
	cityService Service
}

func NewCityHTTPHandler(cityService Service) *CityHTTPHandler {
	return &CityHTTPHandler{
		cityService: cityService,
	}
}

func (h *CityHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/cities",
			Handler: h.Create,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/cities/{id}",
			Handler: h.Delete,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/cities/{id}",
			Handler: h.Update,
		},
		{
			Method:  http.MethodGet,
			Path:    "/cities",
			Handler: h.GetCities,
		},
		{
			Method:  http.MethodGet,
			Path:    "/cities/{id}",
			Handler: h.GetCity,
		},
	}
}
