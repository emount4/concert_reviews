package moderation_transport_http

import (
	"context"
	"net/http"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_ports "github.com/emount4/concert_reviews/internal/core/domain/ports"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type ModerationService interface {
	GetActiveProfileRequests(ctx context.Context, limit, offset *int) ([]domain.ProfileModerationRequest, int, error)
	ApproveProfileRequest(ctx context.Context, id int, moderatorID uuid.UUID) (domain.ProfileModerationRequest, error)
	RejectProfileRequest(ctx context.Context, id int, moderatorID uuid.UUID) (domain.ProfileModerationRequest, error)
	GetAdminLogs(ctx context.Context, moderatorID *uuid.UUID, targetType *string, action *string, limit, offset *int) ([]domain.AdminLog, int, error)
	GetUsers(ctx context.Context, search string, limit, offset *int) ([]domain.User, int, error)
	SetUserBan(ctx context.Context, moderatorID uuid.UUID, moderatorRoleID int, targetID uuid.UUID, isBanned bool) (domain.User, error)
	SetUserRole(ctx context.Context, moderatorID uuid.UUID, targetID uuid.UUID, roleID int) (domain.User, error)
	AnonymizeUser(ctx context.Context, moderatorID uuid.UUID, targetID uuid.UUID) (domain.User, error)
}

type ModerationHTTPHandler struct {
	moderationService ModerationService
	adminLogger       core_ports.AdminLogger
}

func NewModerationHTTPHandler(moderationService ModerationService, adminLogger ...core_ports.AdminLogger) *ModerationHTTPHandler {
	h := &ModerationHTTPHandler{
		moderationService: moderationService,
	}
	if len(adminLogger) > 0 {
		h.adminLogger = adminLogger[0]
	}
	return h
}

func (h *ModerationHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/admin/moderation/profile-requests",
			Handler: h.GetActiveProfileRequests,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodPost,
			Path:    "/admin/moderation/profile-requests/{id}/approve",
			Handler: h.ApproveProfileRequest,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodPost,
			Path:    "/admin/moderation/profile-requests/{id}/reject",
			Handler: h.RejectProfileRequest,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/admin/users",
			Handler: h.GetUsers,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodPost,
			Path:    "/admin/users/{id}/ban",
			Handler: h.SetUserBan,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/admin/users/{id}/role",
			Handler: h.SetUserRole,
			Access:  core_http_server.AccessSuperAdminOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/admin/users/{id}",
			Handler: h.AnonymizeUser,
			Access:  core_http_server.AccessSuperAdminOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/admin/logs",
			Handler: h.GetAdminLogs,
			Access:  core_http_server.AccessSuperAdminOnly,
		},
	}
}
