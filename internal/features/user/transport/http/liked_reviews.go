package user_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	review_transport_http "github.com/emount4/concert_reviews/internal/features/reviews/transport/http"
	"github.com/google/uuid"
)

func (h *UserHTTPHandler) GetLikedReviews(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	userID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		res.ErrorResponse(err, "failed to get pagination params")
		return
	}

	reviews, total, err := h.usersService.GetLikedReviews(ctx, userID, &userID, limit, offset)
	if err != nil {
		res.ErrorResponse(err, "failed to get liked reviews")
		return
	}

	response := review_transport_http.MapDomainListToReviewResponse(reviews)
	response.PageCount = core_http_request.GetPageCount(total, limit)
	res.JSONResponse(response, http.StatusOK)
}

func (h *UserHTTPHandler) GetUserLikedReviews(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	username, err := core_http_request.GetStringPathValue(r, "username")
	if err != nil {
		res.ErrorResponse(err, "missing username")
		return
	}

	user, err := h.usersService.GetProfileByUsername(ctx, username)
	if err != nil {
		res.ErrorResponse(err, "user not found")
		return
	}

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		res.ErrorResponse(err, "failed to get pagination params")
		return
	}

	var viewerID *uuid.UUID
	if id, err := core_http_middleware.GetUserID(ctx); err == nil {
		viewerID = &id
	}

	reviews, total, err := h.usersService.GetLikedReviews(ctx, user.ID, viewerID, limit, offset)
	if err != nil {
		res.ErrorResponse(err, "failed to get user liked reviews")
		return
	}

	response := review_transport_http.MapDomainListToReviewResponse(reviews)
	response.PageCount = core_http_request.GetPageCount(total, limit)
	res.JSONResponse(response, http.StatusOK)
}
