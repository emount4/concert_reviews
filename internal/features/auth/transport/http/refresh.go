package auth_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshResponse struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

func (h *AuthHTTPHandler) Refresh(
	rw http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(
		log,
		rw,
	)

	refreshToken, err := refreshTokenFromRequest(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get refresh token")
		return
	}

	resp, err := h.authService.Refresh(ctx, refreshToken)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to refresh token")
		return
	}

	setRefreshTokenCookie(rw, r, resp.RefreshToken, resp.ExpiresAt)
	rw.Header().Set("Content-Type", "application/json")
	response := RefreshResponse{
		UserID:      resp.User.ID.String(),
		Username:    resp.User.Username,
		AccessToken: resp.AccessToken,
	}
	responseHandler.JSONResponse(response, http.StatusOK)
}

func refreshTokenFromRequest(r *http.Request) (string, error) {
	if cookie, err := r.Cookie(refreshTokenCookieName); err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	var req RefreshRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		return "", err
	}

	return req.RefreshToken, nil
}
