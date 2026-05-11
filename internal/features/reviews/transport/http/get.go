package review_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"github.com/google/uuid"
)

// GetReviews — Получение списка рецензий (Лента / Поиск)
func (h *ReviewHTTPHandler) GetReviews(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	res := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		res.ErrorResponse(err, "failed to get limit/offset")
		return
	}

	// Допустимые поля для сортировки в этой таблице
	allowedSortFields := map[string]bool{
		"date":   true, // -> created_at
		"rating": true, // -> rating_total
		"likes":  true, // -> likes_count
	}

	sort, direction := core_http_request.GetSortParams(r, "date")
	if !allowedSortFields[sort] {
		sort = "date"
	}

	var concertIDPtr *uuid.UUID
	if cID := r.URL.Query().Get("concert_id"); cID != "" {
		if parsed, err := uuid.Parse(cID); err == nil {
			concertIDPtr = &parsed
		}
	}

	artistID, _ := core_http_request.GetIntQueryParam(r, "artist_id")

	reviews, err := h.reviewService.GetReviews(ctx, concertIDPtr, artistID, sort, direction, limit, offset)
	if err != nil {
		res.ErrorResponse(err, "failed to fetch reviews")
		return
	}

	// 5. Маппинг и Ответ
	res.JSONResponse(MapDomainListToReviewResponse(reviews), http.StatusOK)
}

// GetReview — Детальная информация об одной рецензии
func (h *ReviewHTTPHandler) GetReview(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	idStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		res.ErrorResponse(err, "missing review id")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		res.ErrorResponse(err, "invalid review uuid format")
		return
	}

	review, err := h.reviewService.GetReviewByID(ctx, id)
	if err != nil {
		res.ErrorResponse(err, "review not found")
		return
	}

	res.JSONResponse(MapDomainToReviewResponse(review), http.StatusOK)
}

// GetPendingReviews — Очередь для модератора
func (h *ReviewHTTPHandler) GetPendingReviews(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	limit, offset, _ := core_http_request.GetLimitOffsetByQueryParam(r)

	reviews, err := h.reviewService.GetPendingReviews(ctx, limit, offset)
	if err != nil {
		res.ErrorResponse(err, "failed to fetch pending queue")
		return
	}

	res.JSONResponse(MapDomainListToReviewResponse(reviews), http.StatusOK)
}

// GetLikers — Список пользователей, лайкнувших рецензию
func (h *ReviewHTTPHandler) GetLikers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	idStr, _ := core_http_request.GetStringPathValue(r, "id")
	reviewID, _ := uuid.Parse(idStr)

	users, err := h.reviewService.GetLikers(ctx, reviewID)
	if err != nil {
		res.ErrorResponse(err, "failed to get likers list")
		return
	}

	resp := MapUsersToLikersResponse(users)
	res.JSONResponse(resp, http.StatusOK)
}
