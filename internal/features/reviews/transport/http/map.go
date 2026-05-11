package review_transport_http

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

// MapCreateRequestToDomain — маппинг из DTO в Domain
func MapCreateRequestToDomain(req CreateReviewRequest, userID uuid.UUID) domain.Review {
	return domain.Review{
		ReviewID:  uuid.New(), // Генерируем новый ID
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

// MapDomainToReviewResponse — маппинг из Domain в Response DTO
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
		CreatedAt:         r.CreatedAt.Format(time.RFC3339),

		// Автор (используем поля из твоей структуры User)
		Author: AuthorBriefResponse{
			ID:        r.UserID,
			Username:  r.AuthorName, // Поле AuthorName заполняется в репозитории через JOIN
			AvatarURL: r.AuthorAvatar,
		},

		ConcertTitle: r.ConcertTitle,
		LikesCount:   r.LikesCount,
		IsLikedByMe:  r.IsLikedByMe,
	}

	if r.DeletedAt != nil {
		dAt := r.DeletedAt.Format(time.RFC3339)
		resp.DeletedAt = &dAt
	}

	// Маппинг вложенных медиа (ReviewMediaDto на фронте)
	if len(r.Media) > 0 {
		resp.Media = make([]ReviewMediaResponse, len(r.Media))
		for i, m := range r.Media {
			resp.Media[i] = ReviewMediaResponse{
				MediaID:   m.MediaID,
				ReviewID:  m.ReviewID,
				MediaURL:  m.MediaURL,
				MediaType: m.MediaType,
				FileSize:  m.FileSize,
				Status:    string(m.Status),
				CreatedAt: m.CreatedAt.Format(time.RFC3339),
			}
		}
	}

	return resp
}

func MapDomainListToReviewResponse(reviews []domain.Review) []ReviewResponse {
	res := make([]ReviewResponse, len(reviews))
	for i, r := range reviews {
		res[i] = MapDomainToReviewResponse(r)
	}
	return res
}

// Маппинг для списка лайкнувших
func MapUsersToLikersResponse(users []domain.User) []LikerResponse {
	resp := make([]LikerResponse, len(users))
	for i, u := range users {
		resp[i] = LikerResponse{
			ID:        u.ID,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
		}
	}
	return resp
}
