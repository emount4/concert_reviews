package moderation_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *ModerationHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		res.ErrorResponse(err, "failed to get pagination params")
		return
	}

	search := ""
	if value := core_http_request.GetStringQueryParam(r, "search"); value != nil {
		search = *value
	}

	users, total, err := h.moderationService.GetUsers(ctx, search, limit, offset)
	if err != nil {
		res.ErrorResponse(err, "failed to get users")
		return
	}

	response := MapAdminUsersToResponse(users)
	response.PageCount = core_http_request.GetPageCount(total, limit)
	res.JSONResponse(response, http.StatusOK)
}

func (h *ModerationHTTPHandler) SetUserBan(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	moderatorID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}
	moderatorRoleID, err := core_http_middleware.GetRoleID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	targetID, err := getUserIDPathValue(r)
	if err != nil {
		res.ErrorResponse(err, "invalid user id")
		return
	}

	var req SetUserBanRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		res.ErrorResponse(err, "invalid ban request")
		return
	}

	user, err := h.moderationService.SetUserBan(ctx, moderatorID, moderatorRoleID, targetID, req.IsBanned)
	if err != nil {
		res.ErrorResponse(err, "failed to set user ban")
		return
	}

	res.JSONResponse(MapAdminUserToResponse(user), http.StatusOK)
}

func (h *ModerationHTTPHandler) SetUserRole(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	moderatorID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	targetID, err := getUserIDPathValue(r)
	if err != nil {
		res.ErrorResponse(err, "invalid user id")
		return
	}

	var req SetUserRoleRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		res.ErrorResponse(err, "invalid role request")
		return
	}

	user, err := h.moderationService.SetUserRole(ctx, moderatorID, targetID, req.RoleID)
	if err != nil {
		res.ErrorResponse(err, "failed to set user role")
		return
	}

	res.JSONResponse(MapAdminUserToResponse(user), http.StatusOK)
}

func (h *ModerationHTTPHandler) AnonymizeUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	moderatorID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	targetID, err := getUserIDPathValue(r)
	if err != nil {
		res.ErrorResponse(err, "invalid user id")
		return
	}

	user, err := h.moderationService.AnonymizeUser(ctx, moderatorID, targetID)
	if err != nil {
		res.ErrorResponse(err, "failed to anonymize user")
		return
	}

	res.JSONResponse(MapAdminUserToResponse(user), http.StatusOK)
}

func getUserIDPathValue(r *http.Request) (uuid.UUID, error) {
	idStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(idStr)
}
