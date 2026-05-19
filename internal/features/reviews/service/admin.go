package review_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ReviewService) ApproveReview(
	ctx context.Context,
	id uuid.UUID,
	moderatorID uuid.UUID,
	finalTitle string,
	finalText string,
	allowedMediaIDs []uuid.UUID,
) (domain.Review, error) {

	review, err := s.reviewRepository.GetReviewByID(ctx, id)
	if err != nil {
		return domain.Review{}, fmt.Errorf("review not found: %w", err)
	}

	if review.Status == domain.StatusApproved {
		return domain.Review{}, fmt.Errorf("review is already approved: %w", core_errors.ErrConflict)
	}

	err = s.reviewRepository.ApproveReview(ctx, id, moderatorID, finalTitle, finalText, allowedMediaIDs, review)
	if err != nil {
		return domain.Review{}, fmt.Errorf("failed to execute approval transaction: %w", err)
	}

	_ = s.statsRedis.InvalidateGlobalStats(ctx)

	return s.reviewRepository.GetReviewByID(ctx, id)
}

func (s *ReviewService) RejectReview(ctx context.Context, id, moderatorID uuid.UUID, reason string) error {

	err := s.reviewRepository.RejectReview(ctx, id, moderatorID, reason)
	if err != nil {
		return fmt.Errorf("repository rejection failed: %w", err)
	}

	return nil
}

func (s *ReviewService) ReturnReviewToPending(ctx context.Context, id, moderatorID uuid.UUID) error {
	err := s.reviewRepository.ReturnReviewToPending(ctx, id, moderatorID)
	if err != nil {
		return fmt.Errorf("repository return to pending failed: %w", err)
	}

	return nil
}
