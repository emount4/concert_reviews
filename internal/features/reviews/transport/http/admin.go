package review_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *ReviewHTTPHandler) ApproveReview(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	res := core_http_response.NewHTTPResponseHandler(log, rw)

	// 1. Получаем ID модератора из контекста
	moderatorID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	// 2. Получаем ID рецензии из URL
	reviewIDStr, err := core_http_request.GetStringPathValue(r, "id") // проверь имя параметра в роуте
	if err != nil {
		res.ErrorResponse(err, "missing review id")
		return
	}

	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		res.ErrorResponse(err, "invalid review uuid format")
		return
	}

	// 3. Декодируем тело запроса (правки текста и список медиа)
	var req ApproveReviewRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		res.ErrorResponse(err, "failed to decode approval request")
		return
	}

	// 4. Вызываем сервис. Он выполнит транзакцию в БД и сбросит Redis
	approvedReview, err := h.reviewService.ApproveReview(
		ctx,
		reviewID,
		moderatorID,
		req.FinalTitle,
		req.FinalText,
		req.AllowedMediaIDs,
	)

	if err != nil {
		res.ErrorResponse(err, "failed to approve review")
		return
	}

	// 5. Маппим и отдаем обновленную рецензию
	res.JSONResponse(MapDomainToReviewResponse(approvedReview), http.StatusOK)
}
