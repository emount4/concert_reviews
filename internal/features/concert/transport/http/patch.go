package concert_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *ConcertHTTPHandler) UpdateConcert(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	idStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot get concert id path value")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid concert uuid format")
		return
	}

	var req UpdateConcertRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode patch concert request")
		return
	}

	concertPatch := MapPatchConcertReqToDomain(req)

	if err := concertPatch.Validate(); err != nil {
		responseHandler.ErrorResponse(err, "concert patch validation failed")
		return
	}

	updatedConcert, err := h.concertService.UpdateConcert(ctx, id, concertPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot patch concert")
		return
	}

	h.logAdminAction(ctx, "concert_updated", id, map[string]any{
		"title": updatedConcert.Title,
	})

	resp := MapDomainToConcertResponse(updatedConcert)
	responseHandler.JSONResponse(resp, http.StatusOK)
}
