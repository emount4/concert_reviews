package user_transport_http

import (
	"net/http"
	"strings"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
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

	// Получаем рецензии пользователя (approved по умолчанию)
	includeStatuses := r.URL.Query().Get("include_statuses")
	statuses := []string{"approved"}
	if includeStatuses != "" {
		// Парсим статусы из query параметра (например: pending,rejected,approved)
		statuses = nil
		for _, s := range splitTrimmed(includeStatuses, ",") {
			statuses = append(statuses, s)
		}
	}

	reviews, err := h.usersService.GetUserReviews(ctx, userID, statuses)
	if err != nil {
		res.ErrorResponse(err, "failed to fetch user reviews")
		return
	}

	// Маппируем рецензии для включения в ответ
	if len(reviews) > 0 {
		meResp.Reviews = make([]interface{}, len(reviews))
		for i, rev := range reviews {
			meResp.Reviews[i] = map[string]interface{}{
				"review_id":        rev.ReviewID,
				"title":            rev.Title,
				"text":             rev.Text,
				"status":           rev.Status,
				"rating":           rev.RatingTotal,
				"p1":               rev.P1,
				"p2":               rev.P2,
				"p3":               rev.P3,
				"p4":               rev.P4,
				"p5":               rev.P5,
				"rejection_reason": rev.RejectionReason,
				"created_at":       rev.CreatedAt,
				"concert_title":    rev.ConcertTitle,
				"likes_count":      rev.LikesCount,
			}
		}
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
