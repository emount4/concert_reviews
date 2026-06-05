package review_transport_http

import (
	"net/http"
	"strings"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"github.com/google/uuid"
)

// GetReviews — Получение списка рецензий (Лента / Поиск)
func (h *ReviewHTTPHandler) GetReviews(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit/offset")
		return
	}

	// Допустимые поля для сортировки в этой таблице
	allowedSortFields := map[string]bool{
		"date":   true, // -> created_at
		"rating": true, // -> rating_total
		"count":  true, // -> likes_count
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

	artistID, err := core_http_request.GetIntQueryParam(r, "artist_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid artist id")
		return
	}
	venueID, err := core_http_request.GetIntQueryParam(r, "venue_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid venue id")
		return
	}

	// Получаем текущего пользователя для is_liked_by_me (может быть nil если не авторизован)
	userID, _ := core_http_middleware.GetUserID(ctx)

	reviews, total, err := h.reviewService.GetReviews(ctx, &userID, concertIDPtr, artistID, venueID, sort, direction, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to fetch reviews")
		return
	}

	response := MapDomainListToReviewResponse(reviews)
	response.PageCount = core_http_request.GetPageCount(total, limit)
	responseHandler.JSONResponse(response, http.StatusOK)
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

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		res.ErrorResponse(err, "failed to get limit/offset")
		return
	}
	status := strings.TrimSpace(r.URL.Query().Get("status"))

	reviews, total, err := h.reviewService.GetPendingReviews(ctx, status, limit, offset)
	if err != nil {
		res.ErrorResponse(err, "failed to fetch admin reviews")
		return
	}

	resp := MapDomainListToReviewResponse(reviews)
	resp.PageCount = core_http_request.GetPageCount(total, limit)
	res.JSONResponse(resp, http.StatusOK)
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
