package city_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *CityHTTPHandler) GetCities(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'limit/offset' query param")
		return
	}

	citiesDomains, err := h.cityService.GetCities(ctx, limit, offset)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get cities")
		return
	}

	citiesDTO := GetCitiesResponse(CitiesDTOFromDomain(citiesDomains))

	rw.Header().Set("Content-Type", "application/json")
	responseHandler.JSONResponse(citiesDTO, http.StatusOK)
}

func (h *CityHTTPHandler) GetCity(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get city id path value")
		return
	}

	cityDomain, err := h.cityService.GetCityByID(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get city")
		return
	}

	response := CityDTOFromDomain(cityDomain)

	responseHandler.JSONResponse(response, http.StatusOK)
}
