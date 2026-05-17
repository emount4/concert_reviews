package user_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *UserService) GetMe(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	if userID == uuid.Nil {
		return domain.User{}, fmt.Errorf("user id is empty: %w", core_errors.ErrUnauthorized)
	}

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("fetch me from repository: %w", err)
	}

	if !user.IsActive {
		return domain.User{}, fmt.Errorf("account is deactivated: %w", core_errors.ErrNotFound)
	}

	return user, nil
}

func (s *UserService) GetProfileByUsername(ctx context.Context, username string) (domain.User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return domain.User{}, fmt.Errorf("username is required: %w", core_errors.ErrInvalidArgument)
	}

	user, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return domain.User{}, fmt.Errorf("find user %s: %w", username, err)
	}

	if !user.IsActive || user.IsBanned {
		return domain.User{}, fmt.Errorf("user %s is not available: %w", username, core_errors.ErrNotFound)
	}

	return user, nil
}

func (s *UserService) GetUserReviews(ctx context.Context, userID uuid.UUID, includeStatuses []string) ([]domain.Review, error) {
	reviews, err := s.reviewRepository.GetUserReviews(ctx, userID, includeStatuses)
	if err != nil {
		return nil, fmt.Errorf("fetch user reviews from repository: %w", err)
	}
	return reviews, nil
}
