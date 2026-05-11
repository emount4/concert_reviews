package review_transport_http

import (
	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

func MapCreateRequestToDomain(req CreateReviewRequest, userID uuid.UUID) domain.Review {
	return domain.Review{
		UserID:    userID,
		ConcertID: req.ConcertID,
		Title:     req.Title,
		Text:      req.Text,
		P1:        req.P1,
		P2:        req.P2,
		P3:        req.P3,
		P4:        req.P4,
		P5:        req.P5,
		Status:    domain.StatusPending,
		IsVisible: true,
	}
}

func MapDomainToReviewResponse(r domain.Review) ReviewResponse {
	resp := ReviewResponse{
		ReviewID:          r.ReviewID,
		UserID:            r.UserID,
		ConcertID:         r.ConcertID,
		Title:             r.Title,
		Text:              r.Text,
		OriginalText:      r.OriginalText,
		P1:                r.P1,
		P2:                r.P2,
		P3:                r.P3,
		P4:                r.P4,
		P5:                r.P5,
		RatingTotal:       r.RatingTotal,
		Status:            string(r.Status),
		RejectionReason:   r.RejectionReason,
		ModeratedByUserID: r.ModeratedByUserID,
		IsVisible:         r.IsVisible,
		CreatedAt:         r.CreatedAt,
		DeletedAt:         r.DeletedAt,

		// Поля обогащения
		AuthorName:   r.AuthorName,
		AuthorAvatar: r.AuthorAvatar,
		ConcertTitle: r.ConcertTitle,
		LikesCount:   r.LikesCount,
		IsLikedByMe:  r.IsLikedByMe,
	}

	// Маппим вложенные медиа
	if len(r.Media) > 0 {
		resp.Media = make([]ReviewMediaResponse, len(r.Media))
		for i, m := range r.Media {
			resp.Media[i] = ReviewMediaResponse{
				MediaID:   m.MediaID,
				ReviewID:  m.ReviewID,
				MediaURL:  m.MediaURL, // Отдаем KEY
				MediaType: m.MediaType,
				FileSize:  m.FileSize,
				Status:    string(m.Status),
				CreatedAt: m.CreatedAt,
			}
		}
	}

	return resp
}

func MapDomainListToReviewResponse(reviews []domain.Review) ListReviewsResponse {
	items := make([]ReviewResponse, len(reviews))
	for i, r := range reviews {
		items[i] = MapDomainToReviewResponse(r)
	}
	return ListReviewsResponse{Items: items}
}
