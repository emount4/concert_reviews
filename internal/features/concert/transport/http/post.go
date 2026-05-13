package concert_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *ConcertHTTPHandler) CreateConcert(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var req CreateConcertRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate create concert request")
		return
	}

	concertDomain, artistsDomain := MapCreateConcertDTOToDomain(req)

	createdConcert, err := h.concertService.CreateConcert(ctx, concertDomain, artistsDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create concert")
		return
	}

	responseHandler.JSONResponse(MapDomainToConcertResponse(createdConcert), http.StatusCreated)
}

func (h *ConcertHTTPHandler) SuggestConcert(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot get user id from ctx")
		return
	}

	var req CreateSuggestionRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate suggestion request")
		return
	}

	suggestionDomain := MapCreateSuggestionDTOToDomain(req, userID)

	createdSuggestion, err := h.concertService.SuggestConcert(ctx, suggestionDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to submit concert suggestion")
		return
	}

	responseHandler.JSONResponse(MapDomainToSuggestionResponse(createdSuggestion), http.StatusCreated)
}

func (h *ConcertHTTPHandler) RestoreConcert(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	idStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid path value")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid concert uuid format")
		return
	}

	restoredConcert, err := h.concertService.RestoreConcert(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to restore concert")
		return
	}

	responseHandler.JSONResponse(MapDomainToConcertResponse(restoredConcert), http.StatusOK)
}
