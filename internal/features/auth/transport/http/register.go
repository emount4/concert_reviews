package auth_transport_http

import (
	"net/http"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"go.uber.org/zap"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=4,max=50"`

	InitData string `json:"init_data"`
}

type RegisterResponse struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHTTPHandler) Register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var req RegisterRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	user := domainFromDto(req)
	authResponse, err := h.authService.Register(ctx, user, req.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to register user")
		return
	}

	if req.InitData != "" {
		err := h.authService.LinkTG(ctx, authResponse.User, req.InitData)

		if err != nil {
			log.Debug("cannot link tg", zap.Error(err))
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	response := dtoFromDomain(authResponse)

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDto(dto RegisterRequest) core_domain.User {
	return core_domain.NewUser(dto.Username, dto.Email)
}

func dtoFromDomain(response core_domain.AuthResponse) RegisterResponse {
	return RegisterResponse{
		UserID:       response.User.ID.String(),
		Username:     response.User.Username,
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}
}
