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

	moderatorID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

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

	var req ApproveReviewRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		res.ErrorResponse(err, "failed to decode approval request")
		return
	}

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

	h.logAdminAction(ctx, moderatorID, "review_approved", reviewID, map[string]any{
		"allowed_media_count": len(req.AllowedMediaIDs),
	})

	res.JSONResponse(MapDomainToReviewResponse(approvedReview), http.StatusOK)
}

func (h *ReviewHTTPHandler) RejectReview(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	moderatorID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	idStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		res.ErrorResponse(err, "missing review id")
		return
	}

	reviewID, err := uuid.Parse(idStr)
	if err != nil {
		res.ErrorResponse(err, "invalid review uuid format")
		return
	}

	var req RejectReviewRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		res.ErrorResponse(err, "invalid rejection data")
		return
	}

	if err := h.reviewService.RejectReview(ctx, reviewID, moderatorID, req.Reason); err != nil {
		res.ErrorResponse(err, "failed to reject review")
		return
	}

	h.logAdminAction(ctx, moderatorID, "review_rejected", reviewID, map[string]any{
		"reason": req.Reason,
	})

	res.JSONResponse(nil, http.StatusNoContent)
}

func (h *ReviewHTTPHandler) ReturnReviewToPending(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	moderatorID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	idStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		res.ErrorResponse(err, "missing review id")
		return
	}

	reviewID, err := uuid.Parse(idStr)
	if err != nil {
		res.ErrorResponse(err, "invalid review uuid format")
		return
	}

	if err := h.reviewService.ReturnReviewToPending(ctx, reviewID, moderatorID); err != nil {
		res.ErrorResponse(err, "failed to return review to pending")
		return
	}

	h.logAdminAction(ctx, moderatorID, "review_returned_to_pending", reviewID, nil)

	res.JSONResponse(nil, http.StatusNoContent)
}
