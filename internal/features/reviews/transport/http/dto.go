package review_transport_http

import (
	"time"

	"github.com/google/uuid"
)

type CreateReviewRequest struct {
	ConcertID uuid.UUID `json:"concert_id" validate:"required"`
	Title     string    `json:"title" validate:"required,max=255"`
	Text      string    `json:"text" validate:"required,min=100,max=8000"`
	P1        int       `json:"p1" validate:"required,min=1,max=10"`
	P2        int       `json:"p2" validate:"required,min=1,max=10"`
	P3        int       `json:"p3" validate:"required,min=1,max=10"`
	P4        int       `json:"p4" validate:"required,min=1,max=10"`
	P5        int       `json:"p5" validate:"required,min=1,max=10"`
	MediaKeys []string  `json:"media_keys" validate:"omitempty,max=10"`
}

type ApproveReviewRequest struct {
	FinalText       string      `json:"final_text" validate:"required"`
	AllowedMediaIDs []uuid.UUID `json:"allowed_media_ids"`
}
type ReviewMediaResponse struct {
	MediaID   uuid.UUID `json:"media_id"`
	ReviewID  uuid.UUID `json:"review_id"`
	MediaURL  string    `json:"media_url"` // Твой KEY
	MediaType string    `json:"media_type"`
	FileSize  *int64    `json:"file_size"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ReviewResponse struct {
	ReviewID          uuid.UUID  `json:"review_id"`
	UserID            uuid.UUID  `json:"user_id"`
	ConcertID         uuid.UUID  `json:"concert_id"`
	Title             string     `json:"title"`
	Text              string     `json:"text"`
	OriginalText      *string    `json:"original_text"`
	P1                int        `json:"p1"`
	P2                int        `json:"p2"`
	P3                int        `json:"p3"`
	P4                int        `json:"p4"`
	P5                int        `json:"p5"`
	RatingTotal       int        `json:"rating_total"`
	Status            string     `json:"status"`
	RejectionReason   *string    `json:"rejection_reason"`
	ModeratedByUserID *uuid.UUID `json:"moderated_by_user_id"`
	IsVisible         bool       `json:"is_visible"`
	CreatedAt         time.Time  `json:"created_at"`
	DeletedAt         *time.Time `json:"deleted_at"`

	AuthorName   string                `json:"author_name,omitempty"`
	AuthorAvatar *string               `json:"author_avatar_url,omitempty"`
	ConcertTitle string                `json:"concert_title,omitempty"`
	Media        []ReviewMediaResponse `json:"media,omitempty"`
	LikesCount   int                   `json:"likes_count"`
	IsLikedByMe  bool                  `json:"is_liked_by_me"`
}

type ListReviewsResponse struct {
	Items []ReviewResponse `json:"items"`
}
