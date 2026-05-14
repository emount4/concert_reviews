package review_transport_http

import (
	"context"
	"net/http"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type ReviewService interface {
	CreateReview(ctx context.Context, review domain.Review) (domain.Review, error)
	GetReviews(ctx context.Context, userID *uuid.UUID, concertID *uuid.UUID, artistID, venueID *int, sort string, direction string, limit, offset *int) ([]domain.Review, int, error)
	GetReviewByID(ctx context.Context, id uuid.UUID) (domain.Review, error)
	GetPendingReviews(ctx context.Context, limit, offset *int) ([]domain.Review, int, error)

	LikeReview(ctx context.Context, reviewID, userID uuid.UUID) error
	GetLikers(ctx context.Context, reviewID uuid.UUID) ([]domain.User, error)

	ApproveReview(
		ctx context.Context,
		id uuid.UUID,
		moderatorID uuid.UUID,
		finalTitle string,
		finalText string,
		allowedMediaIDs []uuid.UUID,
	) (domain.Review, error)

	RejectReview(ctx context.Context, id, moderatorID uuid.UUID, reason string) error
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
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/review",
			Handler: h.CreateReview,
			Access:  core_http_server.AccessAuthOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/review",
			Handler: h.GetReviews,
		},
		{
			Method:  http.MethodGet,
			Path:    "/review/{id}",
			Handler: h.GetReview,
		},
		{
			Method:  http.MethodGet,
			Path:    "/admin/review",
			Handler: h.GetPendingReviews,
			Access:  core_http_server.AccessAdminOnly,
		},
		{
			Method:  http.MethodPost,
			Path:    "/review/{id}/like",
			Handler: h.LikeReview,
			Access:  core_http_server.AccessAuthOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/review/{id}/like",
			Handler: h.DeleteLike,
			Access:  core_http_server.AccessAuthOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/review/{id}/like",
			Handler: h.GetLikers,
		},

		{
			Method:  http.MethodPost,
			Path:    "/review/{id}/approve",
			Handler: h.ApproveReview,
			Access:  core_http_server.AccessAdminOnly,
		},
	}
}
