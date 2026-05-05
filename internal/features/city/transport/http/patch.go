package city_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *CityHTTPHandler) Update(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get city path param")
		return
	}

	var req UpdateCityRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode request")
		return
	}

	updatedCity, err := h.cityService.Update(r.Context(), id, req.Name, req.Slug, req.Timezone)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to update city")
		return
	}

	responseHandler.JSONResponse(CityDTOFromDomain(updatedCity), http.StatusOK)
}
