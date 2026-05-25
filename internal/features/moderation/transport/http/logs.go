package moderation_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *ModerationHTTPHandler) GetAdminLogs(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		res.ErrorResponse(err, "failed to get pagination params")
		return
	}

	var moderatorID *uuid.UUID
	if value := core_http_request.GetStringQueryParam(r, "moderator_id"); value != nil {
		parsedID, err := uuid.Parse(*value)
		if err != nil {
			res.ErrorResponse(err, "invalid moderator id")
			return
		}
		moderatorID = &parsedID
	}

	targetType := core_http_request.GetStringQueryParam(r, "target_type")
	action := core_http_request.GetStringQueryParam(r, "action")

	logs, total, err := h.moderationService.GetAdminLogs(ctx, moderatorID, targetType, action, limit, offset)
	if err != nil {
		res.ErrorResponse(err, "failed to get admin logs")
		return
	}

	response := MapAdminLogsToResponse(logs)
	response.PageCount = core_http_request.GetPageCount(total, limit)
	res.JSONResponse(response, http.StatusOK)
}
