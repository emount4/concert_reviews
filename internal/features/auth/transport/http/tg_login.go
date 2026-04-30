package auth_transport_http

import (
	"net/http"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

type TGLoginRequest struct {
	InitData string `json:"init_data" validate:"required"`
}

type TGLoginResponse struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHTTPHandler) LoginTG(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var req TGLoginRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	authResponse, err := h.authService.LoginTG(ctx, req.InitData)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login via telegram")
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	response := tgLoginDtoFromDomain(authResponse)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func tgLoginDtoFromDomain(response core_domain.AuthResponse) TGLoginResponse {
	return TGLoginResponse{
		UserID:       response.User.ID.String(),
		Username:     response.User.Username,
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}
}
