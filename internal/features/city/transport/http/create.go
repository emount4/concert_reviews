package city_transport_http

import (
	"net/http"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *CityHTTPHandler) Create(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var req NewCityRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r,
		&req,
	); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode city request")
		return
	}

	city := domainFromDto(req)

	createdCity, err := h.cityService.Create(ctx, city)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create city")
		return
	}

	h.logAdminAction(ctx, "city_created", createdCity.CityID, map[string]any{
		"name": createdCity.Name,
		"slug": createdCity.Slug,
	})

	rw.Header().Set("Content-Type", "application/json")
	response := dtoFromDomain(createdCity)

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDto(cityDTO NewCityRequest) core_models.City {
	return core_models.City{
		Name:     cityDTO.Name,
		Slug:     cityDTO.Slug,
		Timezone: cityDTO.Timezone,
	}
}

func dtoFromDomain(cityDomain core_models.City) NewCityResponse {
	return NewCityResponse{
		CityID:    cityDomain.CityID,
		Name:      cityDomain.Name,
		Slug:      cityDomain.Slug,
		Timezone:  cityDomain.Timezone,
		CreatedAt: cityDomain.CreatedAt,
	}
}
