package user_transport_http

import (
	"net/http"

	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
)

type UserHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
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
			Handler: h.Me,
		},
	}
}
