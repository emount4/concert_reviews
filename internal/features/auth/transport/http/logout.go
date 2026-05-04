package auth_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"go.uber.org/zap"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *AuthHTTPHandler) Logout(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler :=
		core_http_response.NewHTTPResponseHandler(log, rw)

	var req LogoutRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	if err := h.authService.Logout(ctx, req.RefreshToken); err != nil {
		log.Warn("logout failed or token not found", zap.Error(err))
	}

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}
