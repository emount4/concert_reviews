package moderation_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *ModerationHTTPHandler) GetActiveProfileRequests(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		res.ErrorResponse(err, "failed to get pagination params")
		return
	}

	requests, total, err := h.moderationService.GetActiveProfileRequests(ctx, limit, offset)
	if err != nil {
		res.ErrorResponse(err, "failed to get active profile moderation requests")
		return
	}

	response := MapProfileModerationRequestsToResponse(requests)
	response.PageCount = core_http_request.GetPageCount(total, limit)
	res.JSONResponse(response, http.StatusOK)
}
