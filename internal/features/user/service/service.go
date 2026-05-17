package user_service

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository   UserRepository
	reviewRepository ReviewRepository
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
}

type ReviewRepository interface {
	GetUserReviews(ctx context.Context, userID uuid.UUID, includeStatuses []string) ([]domain.Review, error)
}

func NewUserService(userRepository UserRepository, reviewRepository ReviewRepository) *UserService {
	return &UserService{
		userRepository:   userRepository,
		reviewRepository: reviewRepository,
	}
}
