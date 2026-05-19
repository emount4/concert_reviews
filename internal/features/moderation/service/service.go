package moderation_service

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

type ModerationService struct {
	moderationRepository ModerationRepository
}

type ModerationRepository interface {
	GetActiveProfileRequests(ctx context.Context, limit, offset *int) ([]domain.ProfileModerationRequest, int, error)
	ApproveProfileRequest(ctx context.Context, id int, moderatorID uuid.UUID) (domain.ProfileModerationRequest, error)
	RejectProfileRequest(ctx context.Context, id int, moderatorID uuid.UUID) (domain.ProfileModerationRequest, error)
}

func NewModerationService(moderationRepository ModerationRepository) *ModerationService {
	return &ModerationService{
		moderationRepository: moderationRepository,
	}
}
