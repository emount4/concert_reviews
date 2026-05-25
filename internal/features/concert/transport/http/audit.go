package concert_transport_http

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	"github.com/google/uuid"
)

func (h *ConcertHTTPHandler) logAdminAction(ctx context.Context, action string, concertID uuid.UUID, details map[string]any) {
	if h.adminLogger == nil {
		return
	}

	moderatorID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		return
	}

	h.adminLogger.Log(ctx, domain.AdminLog{
		ModeratorID: moderatorID,
		Action:      action,
		TargetType:  domain.LogTargetConcert,
		TargetID:    concertID.String(),
		Details:     details,
	})
}
