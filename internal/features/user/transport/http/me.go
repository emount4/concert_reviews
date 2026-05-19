package user_transport_http

import (
	"net/http"
	"strings"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	review_transport_http "github.com/emount4/concert_reviews/internal/features/reviews/transport/http"
)

type MeResponse struct {
	Message string `json:"message"`
}

func (h *UserHTTPHandler) GetMe(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	userID, _ := core_http_middleware.GetUserID(ctx)

	user, err := h.usersService.GetMe(ctx, userID)
	if err != nil {
		res.ErrorResponse(err, "failed to get profile")
		return
	}

	meResp := MapDomainToMeResponse(user)

	includeStatuses := r.URL.Query().Get("include_statuses")
	statuses := []string{"approved"}
	if includeStatuses != "" {
		statuses = nil
		for _, s := range splitTrimmed(includeStatuses, ",") {
			statuses = append(statuses, s)
		}
	}

	reviews, err := h.usersService.GetUserReviews(ctx, userID, &userID, statuses)
	if err != nil {
		res.ErrorResponse(err, "failed to fetch user reviews")
		return
	}

	if len(reviews) > 0 {
		meResp.Reviews = review_transport_http.MapDomainListToReviewResponse(reviews).Items
	}

	res.JSONResponse(meResp, http.StatusOK)
}

func splitTrimmed(s, sep string) []string {
	var result []string
	parts := strings.Split(s, sep)
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
