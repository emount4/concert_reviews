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

	refreshToken, err := logoutTokenFromRequest(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get refresh token")
		return
	}

	if err := h.authService.Logout(ctx, refreshToken); err != nil {
		log.Warn("logout failed or token not found", zap.Error(err))
	}

	clearRefreshTokenCookie(rw, r)
	responseHandler.JSONResponse(nil, http.StatusNoContent)
}

func logoutTokenFromRequest(r *http.Request) (string, error) {
	if cookie, err := r.Cookie(refreshTokenCookieName); err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	var req LogoutRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		return "", err
	}

	return req.RefreshToken, nil
}
