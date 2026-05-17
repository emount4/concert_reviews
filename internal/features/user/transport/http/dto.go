package user_transport_http

import (
	core_http_types "github.com/emount4/concert_reviews/internal/core/transport/http/types"
)

// --- Requests ---

type UpdateUserRequest struct {
	Username  core_http_types.Nullable[string] `json:"username"`
	Bio       core_http_types.Nullable[string] `json:"bio"`
	AvatarURL core_http_types.Nullable[string] `json:"avatar_url"`
	BannerURL core_http_types.Nullable[string] `json:"banner_url"`
}

// --- Responses ---
type UserStatsDTO struct {
	ReviewsCount       int `json:"reviews_count"`
	LikesGivenCount    int `json:"likes_given_count"`
	LikesReceivedCount int `json:"likes_received_count"`
}

type UserMeResponse struct {
	ID               string        `json:"id"`
	Email            string        `json:"email"`
	Username         string        `json:"username"`
	Bio              *string       `json:"bio,omitempty"`
	AvatarURL        *string       `json:"avatar_url,omitempty"`
	BannerURL        *string       `json:"banner_url,omitempty"`
	TelegramID       *int64        `json:"telegram_id,omitempty"`
	TelegramUsername *string       `json:"telegram_username,omitempty"`
	RoleID           int           `json:"role_id"`
	IsEmailVerified  bool          `json:"is_email_verified"`
	IsBanned         bool          `json:"is_banned"`
	CreatedAt        string        `json:"created_at"`
	Stats            *UserStatsDTO `json:"stats,omitempty"`
	Reviews          []interface{} `json:"reviews,omitempty"` // можно переиспользовать ReviewResponse из reviews/transport/http
}

type PublicProfileResponse struct {
	Username  string        `json:"username"`
	Bio       *string       `json:"bio,omitempty"`
	AvatarURL *string       `json:"avatar_url,omitempty"`
	BannerURL *string       `json:"banner_url,omitempty"`
	CreatedAt string        `json:"created_at"`
	Stats     *UserStatsDTO `json:"stats,omitempty"`
	Reviews   []interface{} `json:"reviews,omitempty"`
}
