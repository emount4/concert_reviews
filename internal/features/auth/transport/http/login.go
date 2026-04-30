package auth_transport_http

import (
	"net/http"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=4,max=50"`
}

type LoginResponse struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHTTPHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var req LoginRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	authResponse, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login")
		return
	}

	response := loginDtoFromDomain(authResponse)
	rw.Header().Set("Content-Type", "application/json")
	responseHandler.JSONResponse(response, http.StatusOK)
}

func loginDtoFromDomain(response core_domain.AuthResponse) LoginResponse {
	return LoginResponse{
		UserID:       response.User.ID.String(),
		Username:     response.User.Username,
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}
}
