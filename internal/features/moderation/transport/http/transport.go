package moderation_transport_http

import (
	"context"
	"net/http"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type ModerationService interface {
	GetActiveProfileRequests(ctx context.Context, limit, offset *int) ([]domain.ProfileModerationRequest, int, error)
	ApproveProfileRequest(ctx context.Context, id int, moderatorID uuid.UUID) (domain.ProfileModerationRequest, error)
	RejectProfileRequest(ctx context.Context, id int, moderatorID uuid.UUID) (domain.ProfileModerationRequest, error)
}

type ModerationHTTPHandler struct {
	moderationService ModerationService
}

func NewModerationHTTPHandler(moderationService ModerationService) *ModerationHTTPHandler {
	return &ModerationHTTPHandler{
		moderationService: moderationService,
	}
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
	}
}
