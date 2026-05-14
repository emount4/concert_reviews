package review_service

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_ports "github.com/emount4/concert_reviews/internal/core/domain/ports"
	"github.com/google/uuid"
)

type ReviewService struct {
	s3               core_ports.S3Provider
	reviewRepository ReviewRepository
	statsRedis       ReviewRedisRepository
}

type ReviewRepository interface {
	CreateReview(ctx context.Context, review domain.Review) (domain.Review, error)
	GetReviewByID(ctx context.Context, id uuid.UUID) (domain.Review, error)

	GetLike(ctx context.Context, reviewID, userID uuid.UUID) (domain.ReviewLike, error)
	GetReviews(
		ctx context.Context,
		userID *uuid.UUID,
		concertID *uuid.UUID,
		artistID *int,
		venueID *int,
		sort string,
		direction string,
		limit, offset *int,
	) (reviews []domain.Review, totalCount int, err error)
	GetLikers(ctx context.Context, reviewID uuid.UUID) ([]domain.User, error)
	GetPendingReviews(ctx context.Context, limit, offset *int) ([]domain.Review, int, error)

	ApproveReview(
		ctx context.Context,
		id uuid.UUID,
		moderatorID uuid.UUID,
		finalTitle string,
		finalText string,
		allowedMediaIDs []uuid.UUID,
		rev domain.Review,
	) error
	RejectReview(ctx context.Context, id, moderatorID uuid.UUID, reason string) error

	CreateLike(ctx context.Context, reviewID, userID uuid.UUID) error
	DeleteLike(ctx context.Context, reviewID, userID uuid.UUID) error

	GetUserReviewCount(ctx context.Context, userID uuid.UUID) (int, error)
}

type ReviewRedisRepository interface {
	InvalidateGlobalStats(ctx context.Context) error
}

func NewReviewService(
	reviewRepository ReviewRepository,
	s3 core_ports.S3Provider,
	redis ReviewRedisRepository,
) *ReviewService {
	return &ReviewService{
		reviewRepository: reviewRepository,
		s3:               s3,
		statsRedis:       redis,
	}
}
