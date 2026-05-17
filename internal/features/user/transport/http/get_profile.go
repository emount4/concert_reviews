package user_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
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

	// Получаем одобренные рецензии пользователя
	reviews, err := h.usersService.GetUserReviews(ctx, user.ID, []string{"approved"})
	if err != nil {
		res.ErrorResponse(err, "failed to fetch user reviews")
		return
	}

	// Маппируем рецензии для включения в ответ (только нужные поля, без причин отказа)
	if len(reviews) > 0 {
		profileResp.Reviews = make([]interface{}, len(reviews))
		for i, rev := range reviews {
			profileResp.Reviews[i] = map[string]interface{}{
				"review_id":     rev.ReviewID,
				"title":         rev.Title,
				"text":          rev.Text,
				"rating":        rev.RatingTotal,
				"p1":            rev.P1,
				"p2":            rev.P2,
				"p3":            rev.P3,
				"p4":            rev.P4,
				"p5":            rev.P5,
				"created_at":    rev.CreatedAt,
				"concert_title": rev.ConcertTitle,
				"likes_count":   rev.LikesCount,
			}
		}
	}

	res.JSONResponse(profileResp, http.StatusOK)
}
