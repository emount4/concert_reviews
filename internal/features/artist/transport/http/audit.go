package artist_transport_http

import (
	"context"
	"strconv"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
)

func (h *ArtistHTTPHandler) logAdminAction(ctx context.Context, action string, artistID int, details map[string]any) {
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
		TargetType:  domain.LogTargetArtist,
		TargetID:    strconv.Itoa(artistID),
		Details:     details,
	})
}
