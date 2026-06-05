package review_postgres_repository

import (
	"encoding/json"
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

type reviewRecord struct {
	ReviewID          uuid.UUID  `db:"review_id"`
	UserID            uuid.UUID  `db:"user_id"`
	ConcertID         uuid.UUID  `db:"concert_id"`
	Title             string     `db:"title"`
	Text              string     `db:"text"`
	OriginalText      *string    `db:"original_text"`
	P1                int        `db:"p1"`
	P2                int        `db:"p2"`
	P3                int        `db:"p3"`
	P4                int        `db:"p4"`
	P5                int        `db:"p5"`
	RatingTotal       int        `db:"rating_total"`
	Status            string     `db:"status"`
	RejectionReason   *string    `db:"rejection_reason"`
	ModeratedByUserID *uuid.UUID `db:"moderated_by_user_id"`
	IsVisible         bool       `db:"is_visible"`
	CreatedAt         time.Time  `db:"created_at"`
	DeletedAt         *time.Time `db:"deleted_at"`

	// Поля из JOIN
	AuthorName       string  `db:"author_name"`
	AuthorAvatar     *string `db:"author_avatar"`
	AuthorIsDeleted  bool    `db:"author_is_deleted"`
	ConcertTitle     string  `db:"concert_title"`
	ConcertPosterURL *string `db:"concert_poster_url"`
	LikesCount       int     `db:"likes_count"`
	IsLikedByMe      bool    `db:"is_liked_by_me"`
	TotalCount       int     `db:"total_count"` // Для пагинации

	// Вложенные медиа (JSON агрегация)
	MediaJSON          []byte `db:"media_json"`
	ConcertArtistsJSON []byte `db:"concert_artists_json"`
}

func (r reviewRecord) MapToDomain() domain.Review {
	rev := domain.Review{
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
		Status:            domain.ModerationStatus(r.Status),
		RejectionReason:   r.RejectionReason,
		ModeratedByUserID: r.ModeratedByUserID,
		IsVisible:         r.IsVisible,
		CreatedAt:         r.CreatedAt,
		DeletedAt:         r.DeletedAt,
		AuthorName:        r.AuthorName,
		AuthorAvatar:      r.AuthorAvatar,
		AuthorIsDeleted:   r.AuthorIsDeleted,
		ConcertTitle:      r.ConcertTitle,
		ConcertPosterURL:  r.ConcertPosterURL,
		LikesCount:        r.LikesCount,
		IsLikedByMe:       r.IsLikedByMe,
	}

	if len(r.MediaJSON) > 0 {
		var media []domain.ReviewMedia
		_ = json.Unmarshal(r.MediaJSON, &media)
		rev.Media = media
	}

	if len(r.ConcertArtistsJSON) > 0 {
		var artists []domain.ConcertArtist
		_ = json.Unmarshal(r.ConcertArtistsJSON, &artists)
		rev.ConcertArtists = artists
	}

	return rev
}
