package moderation_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

func (s *ModerationService) GetActiveProfileRequests(
	ctx context.Context,
	limit, offset *int,
) ([]domain.ProfileModerationRequest, int, error) {
	if limit != nil && *limit < 0 {
		return nil, 0, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, 0, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	requests, total, err := s.moderationRepository.GetActiveProfileRequests(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get active profile requests: %w", err)
	}

	return requests, total, nil
}
