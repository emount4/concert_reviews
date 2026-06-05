package concert_transport_http

import (
	"time"

	core_http_types "github.com/emount4/concert_reviews/internal/core/transport/http/types"
)

// --- Requests ---

type ArtistAssignmentDTO struct {
	ArtistID int  `json:"artist_id" validate:"required"`
	IsMain   bool `json:"is_main"`
}

type CreateConcertRequest struct {
	VenueID   int                   `json:"venue_id" validate:"required"`
	Title     string                `json:"title" validate:"required,max=255"`
	Date      time.Time             `json:"date" validate:"required"`
	PosterKey string                `json:"poster_key" validate:"omitempty,max=2048"`
	Artists   []ArtistAssignmentDTO `json:"artists" validate:"required,min=1"`
}

type UpdateConcertRequest struct {
	VenueID    core_http_types.Nullable[int]       `json:"venue_id"`
	Title      core_http_types.Nullable[string]    `json:"title"`
	Date       core_http_types.Nullable[time.Time] `json:"date"`
	PosterKey  core_http_types.Nullable[string]    `json:"poster_key"`
	IsVerified core_http_types.Nullable[bool]      `json:"is_verified"`
}

type AddConcertArtistRequest struct {
	ArtistID int  `json:"artist_id" validate:"required"`
	IsMain   bool `json:"is_main"`
}

type UpdateConcertArtistRequest struct {
	IsMain bool `json:"is_main"`
}

// --- Responses ---

type VenueBriefResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

type ArtistBriefResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	IsMain bool   `json:"is_main"`
}

type ConcertStatsResponse struct {
	ReviewsCount   int     `json:"reviews_count"`
	AvgRatingTotal float64 `json:"avg_rating_total"`
	AvgP1          float64 `json:"avg_p1"`
	AvgP2          float64 `json:"avg_p2"`
	AvgP3          float64 `json:"avg_p3"`
	AvgP4          float64 `json:"avg_p4"`
	AvgP5          float64 `json:"avg_p5"`
	FavoritesCount int     `json:"favorites_count"`
	UpdatedAt      string  `json:"updated_at"`
}

type ConcertResponse struct {
	ID               string                `json:"id"`
	Title            string                `json:"title"`
	Date             string                `json:"date"`
	PosterURL        string                `json:"poster_url,omitempty"` // Тут будет лежать KEY
	IsVerified       bool                  `json:"is_verified"`
	Venue            VenueBriefResponse    `json:"venue"`
	Artists          []ArtistBriefResponse `json:"artists"`
	Stats            *ConcertStatsResponse `json:"stats,omitempty"` // Вложенная структура
	CreatedAt        string                `json:"created_at"`
	UserReviewStatus *string               `json:"user_review_status,omitempty"`
}

type ListConcertsResponse struct {
	Items     []ConcertResponse `json:"items"`
	PageCount int               `json:"page_count"`
}

// --- Suggestion Requests ---

type CreateSuggestionRequest struct {
	ArtistName string    `json:"artist_name" validate:"required_without=VenueName,max=255"`
	VenueName  string    `json:"venue_name" validate:"required_without=ArtistName,max=255"`
	Date       time.Time `json:"date" validate:"required"`
	Info       string    `json:"info" validate:"omitempty,max=2000"`
}

// --- Suggestion Responses ---

type SuggestionResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ArtistName string    `json:"artist_name"`
	VenueName  string    `json:"venue_name"`
	Date       time.Time `json:"date"`
	Info       *string   `json:"info,omitempty"`
	CreatedAt  string    `json:"created_at"`
}

type ListSuggestionsResponse struct {
	Items []SuggestionResponse `json:"items"`
}

type ConcertResponseAdmin struct {
	ID              string                `json:"id"`
	VenueID         int                   `json:"venue_id"`
	Title           string                `json:"title"`
	Date            string                `json:"date"`
	PosterURL       string                `json:"poster_url,omitempty"`
	IsVerified      bool                  `json:"is_verified"`
	CreatedByUserID *string               `json:"created_by_user_id,omitempty"`
	CreatedAt       string                `json:"created_at"`
	DeletedAt       *string               `json:"deleted_at,omitempty"`
	Venue           VenueBriefResponse    `json:"venue"`
	Artists         []ArtistBriefResponse `json:"artists"`
	Stats           *ConcertStatsResponse `json:"stats,omitempty"`
}

type ListConcertsAdminResponse struct {
	Items     []ConcertResponseAdmin `json:"items"`
	PageCount int                    `json:"page_count"`
}
