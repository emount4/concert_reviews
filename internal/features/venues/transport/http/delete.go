package venue_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

// DeleteVenueHard — физическое удаление
func (h *VenueHTTPHandler) DeleteVenueHard(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid venue id in url")
		return
	}

	if err := h.venueService.DeleteVenueHard(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, "failed to hard delete venue")
		return
	}

	h.logAdminAction(ctx, "venue_hard_deleted", id, nil)

	rw.WriteHeader(http.StatusNoContent)
}

// DeleteVenueSoft — логическое удаление (soft delete)
func (h *VenueHTTPHandler) DeleteVenueSoft(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid venue id in url")
		return
	}

	if err := h.venueService.DeleteVenueSoft(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, "failed to soft delete venue")
		return
	}

	h.logAdminAction(ctx, "venue_soft_deleted", id, nil)

	rw.WriteHeader(http.StatusNoContent)
}

// RestoreVenue — восстановление после soft delete
func (h *VenueHTTPHandler) RestoreVenue(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid venue id in url")
		return
	}

	restoredVenue, err := h.venueService.RestoreVenue(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to restore venue")
		return
	}

	h.logAdminAction(ctx, "venue_restored", id, map[string]any{
		"name": restoredVenue.Name,
	})

	rw.Header().Set("Content-Type", "application/json")
	response := MapDomainToVenueResponse(restoredVenue)
	responseHandler.JSONResponse(response, http.StatusOK)
}
