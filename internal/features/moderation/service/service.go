package moderation_service

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

type ModerationService struct {
	moderationRepository ModerationRepository
	statsCache           GlobalStatsCache
}

type ModerationRepository interface {
	GetActiveProfileRequests(ctx context.Context, limit, offset *int) ([]domain.ProfileModerationRequest, int, error)
	ApproveProfileRequest(ctx context.Context, id int, moderatorID uuid.UUID) (domain.ProfileModerationRequest, error)
	RejectProfileRequest(ctx context.Context, id int, moderatorID uuid.UUID) (domain.ProfileModerationRequest, error)
	GetUsers(ctx context.Context, search string, limit, offset *int) ([]domain.User, int, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	SetUserBan(ctx context.Context, moderatorID uuid.UUID, targetID uuid.UUID, isBanned bool) (domain.User, error)
	SetUserRole(ctx context.Context, moderatorID uuid.UUID, targetID uuid.UUID, roleID int) (domain.User, error)
	AnonymizeUser(ctx context.Context, moderatorID uuid.UUID, targetID uuid.UUID) (domain.User, error)
	GetAdminLogs(ctx context.Context, moderatorID *uuid.UUID, targetType *string, action *string, limit, offset *int) ([]domain.AdminLog, int, error)
}

type GlobalStatsCache interface {
	InvalidateGlobalStats(ctx context.Context) error
}

func NewModerationService(moderationRepository ModerationRepository, statsCache ...GlobalStatsCache) *ModerationService {
	var cache GlobalStatsCache
	if len(statsCache) > 0 {
		cache = statsCache[0]
	}

	return &ModerationService{
		moderationRepository: moderationRepository,
		statsCache:           cache,
	}
}
