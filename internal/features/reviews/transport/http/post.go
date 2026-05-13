package review_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *ReviewHTTPHandler) CreateReview(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	res := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	var req CreateReviewRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		res.ErrorResponse(err, "invalid review data")
		return
	}

	reviewDomain := MapCreateRequestToDomain(req, userID)

	created, err := h.reviewService.CreateReview(ctx, reviewDomain)
	if err != nil {
		res.ErrorResponse(err, "failed to create review")
		return
	}

	res.JSONResponse(MapDomainToReviewResponse(created), http.StatusCreated)
}

func (h *ReviewHTTPHandler) LikeReview(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	res := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, _ := core_http_middleware.GetUserID(ctx)

	idStr, _ := core_http_request.GetStringPathValue(r, "id")
	reviewID, err := uuid.Parse(idStr)
	if err != nil {
		res.ErrorResponse(err, "invalid review uuid")
		return
	}

	if err := h.reviewService.LikeReview(ctx, reviewID, userID); err != nil {
		res.ErrorResponse(err, "failed to toggle like")
		return
	}

	res.JSONResponse(nil, http.StatusNoContent)
}
