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

func (h *UserHTTPHandler) GetProfile(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	username, err := core_http_request.GetStringPathValue(r, "username")
	if err != nil {
		res.ErrorResponse(err, "invalid username")
		return
	}

	user, err := h.usersService.GetProfileByUsername(ctx, username)
	if err != nil {
		res.ErrorResponse(err, "user not found")
		return
	}

	profileResp := MapDomainToPublicResponse(user)

	var viewerID *uuid.UUID
	if id, err := core_http_middleware.GetUserID(ctx); err == nil {
		viewerID = &id
	}

	reviews, err := h.usersService.GetUserReviews(ctx, user.ID, viewerID, []string{"approved"})
	if err != nil {
		res.ErrorResponse(err, "failed to fetch user reviews")
		return
	}

	if len(reviews) > 0 {
		profileResp.Reviews = review_transport_http.MapDomainListToReviewResponse(reviews).Items
	}

	res.JSONResponse(profileResp, http.StatusOK)
}
