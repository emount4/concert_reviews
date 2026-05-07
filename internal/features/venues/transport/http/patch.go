package venue_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *VenueHTTPHandler) UpdateVenue(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot get venue id path value")
		return
	}

	var req UpdateVenueRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil { // ← новая функция, см. ниже
		responseHandler.ErrorResponse(err, "failed to decode patch venue request")
		return
	}

	venuePatch := MapPatchVenueReqToDomain(req)

	if err := venuePatch.Validate(); err != nil {
		responseHandler.ErrorResponse(err, "venue patch validation failed")
		return
	}

	updatedVenue, err := h.venueService.UpdateVenue(ctx, id, venuePatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot patch venue")
		return
	}

	resp := MapDomainToVenueResponse(updatedVenue)
	responseHandler.JSONResponse(resp, http.StatusOK)
}
