package moderation_transport_http

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

func (h *ModerationHTTPHandler) logAdminAction(ctx context.Context, moderatorID uuid.UUID, action, targetType, targetID string, details map[string]any) {
	if h.adminLogger == nil {
		return
	}

	h.adminLogger.Log(ctx, domain.AdminLog{
		ModeratorID: moderatorID,
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		Details:     details,
	})
}
