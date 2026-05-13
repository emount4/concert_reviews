package review_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ReviewService) GetReviews(
	ctx context.Context,
	userID *uuid.UUID,
	concertID *uuid.UUID,
	artistID, venueID *int,
	sort string,
	direction string,
	limit, offset *int,
) ([]domain.Review, int, error) {

	// Здесь можно добавить валидацию параметров фильтрации
	if limit != nil && (*limit <= 0 || *limit > 100) {
		return nil, 0, fmt.Errorf("invalid limit: %w", core_errors.ErrInvalidArgument)
	}

	reviews, total, err := s.reviewRepository.GetReviews(ctx, userID, concertID, artistID, venueID, sort, direction, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fetch reviews from repository: %w", err)
	}

	return reviews, total, nil
}

func (s *ReviewService) GetReviewByID(ctx context.Context, id uuid.UUID) (domain.Review, error) {
	if id == uuid.Nil {
		return domain.Review{}, fmt.Errorf("invalid review id: %w", core_errors.ErrInvalidArgument)
	}

	review, err := s.reviewRepository.GetReviewByID(ctx, id)
	if err != nil {
		return domain.Review{}, fmt.Errorf("get review %s: %w", id, err)
	}

	return review, nil
}

func (s *ReviewService) GetPendingReviews(ctx context.Context, limit, offset *int) ([]domain.Review, int, error) {

	reviews, total, err := s.reviewRepository.GetPendingReviews(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fetch pending queue: %w", err)
	}

	return reviews, total, nil
}

func (s *ReviewService) GetLikers(ctx context.Context, reviewID uuid.UUID) ([]domain.User, error) {
	if s.reviewRepository == nil {
		return nil, core_errors.ErrRepositoryNotConfigured
	}

	if reviewID == uuid.Nil {
		return nil, fmt.Errorf("review id is nil: %w", core_errors.ErrInvalidArgument)
	}

	_, err := s.reviewRepository.GetReviewByID(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("verify review %s existence: %w", reviewID, err)
	}

	users, err := s.reviewRepository.GetLikers(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("fetch likers for review %s from repository: %w", reviewID, err)
	}

	return users, nil
}
