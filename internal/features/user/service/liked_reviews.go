package user_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *UserService) GetLikedReviews(
	ctx context.Context,
	userID uuid.UUID,
	viewerID *uuid.UUID,
	limit, offset *int,
) ([]domain.Review, int, error) {
	if userID == uuid.Nil {
		return nil, 0, fmt.Errorf("user id is empty: %w", core_errors.ErrInvalidArgument)
	}
	if viewerID != nil && *viewerID == uuid.Nil {
		viewerID = nil
	}
	if limit != nil && *limit < 0 {
		return nil, 0, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, 0, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("fetch user from repository: %w", err)
	}
	if !user.IsActive || user.IsBanned {
		return nil, 0, fmt.Errorf("user is not available: %w", core_errors.ErrNotFound)
	}

	reviews, total, err := s.reviewRepository.GetLikedReviews(ctx, userID, viewerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fetch liked reviews from repository: %w", err)
	}
	return reviews, total, nil
}
