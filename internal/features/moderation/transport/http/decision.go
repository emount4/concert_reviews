package moderation_transport_http

import (
	"net/http"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *ModerationHTTPHandler) ApproveProfileRequest(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	moderatorID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		res.ErrorResponse(err, "invalid profile moderation request id")
		return
	}

	req, err := h.moderationService.ApproveProfileRequest(ctx, id, moderatorID)
	if err != nil {
		res.ErrorResponse(err, "failed to approve profile moderation request")
		return
	}

	h.logAdminAction(ctx, moderatorID, "profile_request_approved", domain.LogTargetUser, req.UserID.String(), map[string]any{
		"moderation_id": id,
		"field_name":    req.FieldName,
	})

	res.JSONResponse(MapProfileModerationRequestToResponse(req), http.StatusOK)
}

func (h *ModerationHTTPHandler) RejectProfileRequest(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	moderatorID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		res.ErrorResponse(err, "invalid profile moderation request id")
		return
	}

	req, err := h.moderationService.RejectProfileRequest(ctx, id, moderatorID)
	if err != nil {
		res.ErrorResponse(err, "failed to reject profile moderation request")
		return
	}

	h.logAdminAction(ctx, moderatorID, "profile_request_rejected", domain.LogTargetUser, req.UserID.String(), map[string]any{
		"moderation_id": id,
		"field_name":    req.FieldName,
	})

	res.JSONResponse(MapProfileModerationRequestToResponse(req), http.StatusOK)
}
