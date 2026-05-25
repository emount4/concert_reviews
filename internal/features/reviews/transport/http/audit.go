package review_transport_http

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

func (h *ReviewHTTPHandler) logAdminAction(ctx context.Context, moderatorID uuid.UUID, action string, reviewID uuid.UUID, details map[string]any) {
	if h.adminLogger == nil {
		return
	}

	h.adminLogger.Log(ctx, domain.AdminLog{
		ModeratorID: moderatorID,
		Action:      action,
		TargetType:  domain.LogTargetReview,
		TargetID:    reviewID.String(),
		Details:     details,
	})
}
