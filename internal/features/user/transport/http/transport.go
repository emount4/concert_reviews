package user_transport_http

import (
	"context"
	"net/http"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type UserHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	GetMe(ctx context.Context, userID uuid.UUID) (domain.User, error)
	GetProfileByUsername(ctx context.Context, username string) (domain.User, error)
	GetUserReviews(ctx context.Context, userID uuid.UUID, includeStatuses []string) ([]domain.Review, error)
}

func NewUsersHTTPHandler(userService UsersService) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService,
	}
}

func (h *UserHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users/me",
			Handler: h.GetMe,
			Access:  core_http_server.AccessAuthOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{username}",
			Handler: h.GetProfile,
		},
	}
}
