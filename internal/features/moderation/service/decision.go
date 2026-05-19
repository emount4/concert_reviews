package moderation_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ModerationService) ApproveProfileRequest(
	ctx context.Context,
	id int,
	moderatorID uuid.UUID,
) (domain.ProfileModerationRequest, error) {
	if id <= 0 {
		return domain.ProfileModerationRequest{}, fmt.Errorf("moderation request id must be positive: %w", core_errors.ErrInvalidArgument)
	}
	if moderatorID == uuid.Nil {
		return domain.ProfileModerationRequest{}, fmt.Errorf("moderator id is empty: %w", core_errors.ErrUnauthorized)
	}

	req, err := s.moderationRepository.ApproveProfileRequest(ctx, id, moderatorID)
	if err != nil {
		return domain.ProfileModerationRequest{}, fmt.Errorf("approve profile request: %w", err)
	}
	return req, nil
}

func (s *ModerationService) RejectProfileRequest(
	ctx context.Context,
	id int,
	moderatorID uuid.UUID,
) (domain.ProfileModerationRequest, error) {
	if id <= 0 {
		return domain.ProfileModerationRequest{}, fmt.Errorf("moderation request id must be positive: %w", core_errors.ErrInvalidArgument)
	}
	if moderatorID == uuid.Nil {
		return domain.ProfileModerationRequest{}, fmt.Errorf("moderator id is empty: %w", core_errors.ErrUnauthorized)
	}

	req, err := s.moderationRepository.RejectProfileRequest(ctx, id, moderatorID)
	if err != nil {
		return domain.ProfileModerationRequest{}, fmt.Errorf("reject profile request: %w", err)
	}
	return req, nil
}
