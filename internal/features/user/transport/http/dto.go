package user_transport_http

import (
	core_http_types "github.com/emount4/concert_reviews/internal/core/transport/http/types"
	review_transport_http "github.com/emount4/concert_reviews/internal/features/reviews/transport/http"
	"github.com/google/uuid"
)

// --- Requests ---

type UpdateUserRequest struct {
	Username  core_http_types.Nullable[string] `json:"username"`
	Bio       core_http_types.Nullable[string] `json:"bio"`
	AvatarKey core_http_types.Nullable[string] `json:"avatar_key"`
	BannerKey core_http_types.Nullable[string] `json:"banner_key"`
}

// --- Responses ---
type UserStatsDTO struct {
	ReviewsCount       int `json:"reviews_count"`
	LikesGivenCount    int `json:"likes_given_count"`
	LikesReceivedCount int `json:"likes_received_count"`
}

type UserMeResponse struct {
	ID               string                                 `json:"id"`
	Email            string                                 `json:"email"`
	Username         string                                 `json:"username"`
	Bio              *string                                `json:"bio,omitempty"`
	AvatarURL        *string                                `json:"avatar_url,omitempty"`
	BannerURL        *string                                `json:"banner_url,omitempty"`
	TelegramID       *int64                                 `json:"telegram_id,omitempty"`
	TelegramUsername *string                                `json:"telegram_username,omitempty"`
	RoleID           int                                    `json:"role_id"`
	IsEmailVerified  bool                                   `json:"is_email_verified"`
	IsBanned         bool                                   `json:"is_banned"`
	CreatedAt        string                                 `json:"created_at"`
	Stats            *UserStatsDTO                          `json:"stats,omitempty"`
	Reviews          []review_transport_http.ReviewResponse `json:"reviews,omitempty"`
}

type PublicProfileResponse struct {
	ID        string                                 `json:"id"`
	Username  string                                 `json:"username"`
	Bio       *string                                `json:"bio,omitempty"`
	AvatarURL *string                                `json:"avatar_url,omitempty"`
	BannerURL *string                                `json:"banner_url,omitempty"`
	CreatedAt string                                 `json:"created_at"`
	Stats     *UserStatsDTO                          `json:"stats,omitempty"`
	Reviews   []review_transport_http.ReviewResponse `json:"reviews,omitempty"`
}

type ProfileModerationRequestResponse struct {
	ID                int        `json:"id"`
	FieldName         string     `json:"field_name"`
	OldValue          *string    `json:"old_value"`
	NewValue          *string    `json:"new_value"`
	Status            string     `json:"status"`
	ModeratedByUserID *uuid.UUID `json:"moderated_by_user_id,omitempty"`
	CreatedAt         string     `json:"created_at"`
	UpdatedAt         string     `json:"updated_at"`
}

type ListProfileModerationRequestsResponse struct {
	Items []ProfileModerationRequestResponse `json:"items"`
}
