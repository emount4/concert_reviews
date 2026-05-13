package review_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ReviewService) CreateReview(ctx context.Context, review domain.Review) (domain.Review, error) {
	if err := review.Validate(); err != nil {
		return domain.Review{}, fmt.Errorf("validate review: %w", err)
	}

	review.RatingTotal = review.CalculateRating()

	for _, m := range review.Media {
		if _, err := s.s3.FileExists(ctx, m.MediaURL); err != nil {
			return domain.Review{}, fmt.Errorf("media file %s not found in S3: %w", m.MediaURL, err)
		}
	}

	review.ReviewID = uuid.New()
	review.Status = domain.StatusPending
	review.IsVisible = true

	created, err := s.reviewRepository.CreateReview(ctx, review)
	if err != nil {
		return domain.Review{}, fmt.Errorf("save review to repository: %w", err)
	}

	return created, nil
}

func (s *ReviewService) LikeReview(ctx context.Context, reviewID, userID uuid.UUID) error {
	review, err := s.reviewRepository.GetReviewByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("review not found: %w", err)
	}
	if review.UserID == userID {
		return fmt.Errorf("self-liking is forbidden: %w", core_errors.ErrConflict)
	}

	if review.Status != domain.StatusApproved {
		return fmt.Errorf("cannot like unapproved review: %w", core_errors.ErrConflict)
	}

	count, err := s.reviewRepository.GetUserReviewCount(ctx, userID)
	if err != nil {
		return fmt.Errorf("check user review count: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("only authors with at least one approved review can like: %w", core_errors.ErrForbidden)
	}

	like, err := s.reviewRepository.GetLike(ctx, reviewID, userID)

	if errors.Is(err, core_errors.ErrNotFound) {
		return s.reviewRepository.CreateLike(ctx, reviewID, userID)
	}

	if err != nil {
		return err
	}

	if time.Since(like.CreatedAt) > 5*time.Second {
		return fmt.Errorf("unlike period expired: %w", core_errors.ErrConflict)
	}

	return s.reviewRepository.DeleteLike(ctx, reviewID, userID)
}
