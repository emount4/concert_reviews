package review_transport_http

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type ReviewService interface {
	CreateReview(ctx context.Context, review domain.Review) (domain.Review, error)
	LikeReview(ctx context.Context, reviewID, userID uuid.UUID) error
	GetReviews(ctx context.Context, concertID *uuid.UUID, artistID *int, sort string, direction string, limit, offset *int) ([]domain.Review, error)
	GetReviewByID(ctx context.Context, id uuid.UUID) (domain.Review, error)
	GetPendingReviews(ctx context.Context, limit, offset *int) ([]domain.Review, error)
	GetLikers(ctx context.Context, reviewID uuid.UUID) ([]domain.User, error)
}

type ReviewHTTPHandler struct {
	reviewService ReviewService
}

func NewReviewHTTPHandler(s ReviewService) *ReviewHTTPHandler {
	return &ReviewHTTPHandler{
		reviewService: s,
	}
}

func (h *ReviewHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{}
}
