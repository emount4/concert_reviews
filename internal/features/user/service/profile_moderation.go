package user_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *UserService) GetMyProfileModerationRequests(
	ctx context.Context,
	userID uuid.UUID,
	status *domain.ModerationStatus,
) ([]domain.ProfileModerationRequest, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user id is empty: %w", core_errors.ErrUnauthorized)
	}

	if status != nil {
		switch *status {
		case domain.StatusPending, domain.StatusApproved, domain.StatusRejected:
		default:
			return nil, fmt.Errorf("invalid moderation status: %w", core_errors.ErrInvalidArgument)
		}
	}

	requests, err := s.userRepository.GetProfileModerationRequests(ctx, userID, status)
	if err != nil {
		return nil, fmt.Errorf("get profile moderation requests: %w", err)
	}

	return requests, nil
}
