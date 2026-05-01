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
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

	var req RefreshRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	resp, err := h.authService.Refresh(ctx, req.RefreshToken)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to refresh token")
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	response := RefreshResponse{
		UserID:       resp.User.ID.String(),
		Username:     resp.User.Username,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}
	responseHandler.JSONResponse(response, http.StatusOK)
}
