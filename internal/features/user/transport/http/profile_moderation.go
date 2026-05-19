package user_transport_http

import (
	"net/http"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *UserHTTPHandler) GetMyProfileModerationRequests(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	userID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	status, err := parseModerationStatusQuery(r.URL.Query().Get("status"))
	if err != nil {
		res.ErrorResponse(err, "invalid status query param")
		return
	}

	requests, err := h.usersService.GetMyProfileModerationRequests(ctx, userID, status)
	if err != nil {
		res.ErrorResponse(err, "failed to get profile moderation requests")
		return
	}

	res.JSONResponse(MapProfileModerationRequestsToResponse(requests), http.StatusOK)
}

func parseModerationStatusQuery(raw string) (*domain.ModerationStatus, error) {
	if raw == "" {
		return nil, nil
	}

	status := domain.ModerationStatus(raw)
	switch status {
	case domain.StatusPending, domain.StatusApproved, domain.StatusRejected:
		return &status, nil
	default:
		return nil, core_errors.ErrInvalidArgument
	}
}
