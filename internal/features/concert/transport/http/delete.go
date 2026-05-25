package concert_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *ConcertHTTPHandler) DeleteConcertSoft(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	idStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot get concert id from path")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid concert uuid format")
		return
	}

	if err := h.concertService.DeleteConcertSoft(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, "failed to soft delete concert")
		return
	}

	h.logAdminAction(ctx, "concert_soft_deleted", id, nil)

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}

func (h *ConcertHTTPHandler) DeleteConcertHard(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	idStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot get concert id from path")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid concert uuid format")
		return
	}

	if err := h.concertService.DeleteConcertHard(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, "failed to hard delete concert")
		return
	}

	h.logAdminAction(ctx, "concert_hard_deleted", id, nil)

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}

func (h *ConcertHTTPHandler) DeleteSuggestionAdmin(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	idStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot get suggestion id from path")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid suggestion uuid format")
		return
	}

	if err := h.concertService.DeleteSuggestionAdmin(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete concert suggestion")
		return
	}

	h.logAdminAction(ctx, "concert_suggestion_deleted", id, nil)

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}
