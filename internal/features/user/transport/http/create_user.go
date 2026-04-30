package user_transport_http

import "net/http"

type CreateUserRequest struct {
	UserName string
	Email    string
	RoleID   int

	PasswordHash string

	TelegramID       *int64
	TelegramUsername *string
}

type CreateUserResponse struct {
}

func (h *UserHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {

}
