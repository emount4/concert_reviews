package auth_transport_http

import (
	"context"
	"net/http"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, user core_domain.User, password string) (core_domain.AuthResponse, error)
	Login(ctx context.Context, email, password string) (core_domain.AuthResponse, error)
	Logout(ctx context.Context, token string) error
	Refresh(ctx context.Context, oldToken string) (core_domain.AuthResponse, error)

	LinkTG(ctx context.Context, user core_domain.User, initData string) error
	Link(ctx context.Context, id uuid.UUID, initData string) error
	LoginTG(ctx context.Context, initData string) (core_domain.AuthResponse, error)
}

type AuthHTTPHandler struct {
	authService Service
}

func NewAuthHTTPHandler(authService Service) *AuthHTTPHandler {
	return &AuthHTTPHandler{authService: authService}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/tg-login",
			Handler: h.LoginTG,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.Refresh,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: h.Logout,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/tg-link",
			Handler: h.Link,
			Access:  core_http_server.AccessAuthOnly,
		},
	}
}
