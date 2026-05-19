package user_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *UserHTTPHandler) PatchMe(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	userID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	var req UpdateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		res.ErrorResponse(err, "invalid profile patch")
		return
	}

	updated, err := h.usersService.UpdateMe(ctx, userID, MapUpdateUserRequestToDomain(req))
	if err != nil {
		res.ErrorResponse(err, "failed to update profile")
		return
	}

	res.JSONResponse(MapDomainToMeResponse(updated), http.StatusOK)
}
