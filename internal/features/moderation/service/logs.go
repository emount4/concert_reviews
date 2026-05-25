package moderation_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ModerationService) GetAdminLogs(
	ctx context.Context,
	moderatorID *uuid.UUID,
	targetType *string,
	action *string,
	limit, offset *int,
) ([]domain.AdminLog, int, error) {
	if limit != nil && *limit < 0 {
		return nil, 0, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, 0, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if moderatorID != nil && *moderatorID == uuid.Nil {
		return nil, 0, fmt.Errorf("moderator_id is empty: %w", core_errors.ErrInvalidArgument)
	}
	if targetType != nil {
		trimmed := strings.TrimSpace(*targetType)
		if err := validateAdminLogTargetType(trimmed); err != nil {
			return nil, 0, err
		}
		targetType = &trimmed
	}
	if action != nil {
		trimmed := strings.TrimSpace(*action)
		action = &trimmed
	}

	logs, total, err := s.moderationRepository.GetAdminLogs(ctx, moderatorID, targetType, action, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get admin logs: %w", err)
	}

	return logs, total, nil
}

func validateAdminLogTargetType(targetType string) error {
	switch targetType {
	case "", domain.LogTargetUser, domain.LogTargetReview, domain.LogTargetArtist, domain.LogTargetVenue, domain.LogTargetConcert, domain.LogTargetCity:
		return nil
	default:
		return fmt.Errorf("invalid log target_type: %w", core_errors.ErrInvalidArgument)
	}
}
