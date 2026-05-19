package user_service

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_ports "github.com/emount4/concert_reviews/internal/core/domain/ports"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository   UserRepository
	reviewRepository ReviewRepository
	s3               core_ports.S3Provider
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	SubmitProfilePatch(ctx context.Context, current domain.User, patch domain.UserPatch) error
	GetProfileModerationRequests(ctx context.Context, userID uuid.UUID, status *domain.ModerationStatus) ([]domain.ProfileModerationRequest, error)
}

type ReviewRepository interface {
	GetUserReviews(ctx context.Context, userID uuid.UUID, viewerID *uuid.UUID, includeStatuses []string) ([]domain.Review, error)
	GetLikedReviews(ctx context.Context, userID uuid.UUID, viewerID *uuid.UUID, limit, offset *int) ([]domain.Review, int, error)
}

func NewUserService(
	userRepository UserRepository,
	reviewRepository ReviewRepository,
	s3 core_ports.S3Provider,
) *UserService {
	return &UserService{
		userRepository:   userRepository,
		reviewRepository: reviewRepository,
		s3:               s3,
	}
}
