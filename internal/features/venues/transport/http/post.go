package venue_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *VenueHTTPHandler) CreateVenue(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var req CreateVenueRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode venue creation request")
		return
	}

	venue := MapCreateVenueDTOToDomain(req)

	createdVenue, err := h.venueService.CreateVenue(ctx, venue)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create venue")
		return
	}

	h.logAdminAction(ctx, "venue_created", createdVenue.VenueID, map[string]any{
		"name": createdVenue.Name,
	})

	rw.Header().Set("Content-Type", "application/json")
	response := MapDomainToVenueResponse(createdVenue)

	responseHandler.JSONResponse(response, http.StatusCreated)
}
