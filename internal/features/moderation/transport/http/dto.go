package moderation_transport_http

import "github.com/google/uuid"

type ProfileModerationUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
}

type ProfileModerationRequestResponse struct {
	ID                int                           `json:"id"`
	User              ProfileModerationUserResponse `json:"user"`
	FieldName         string                        `json:"field_name"`
	OldValue          *string                       `json:"old_value"`
	NewValue          *string                       `json:"new_value"`
	Status            string                        `json:"status"`
	ModeratedByUserID *uuid.UUID                    `json:"moderated_by_user_id,omitempty"`
	CreatedAt         string                        `json:"created_at"`
	UpdatedAt         string                        `json:"updated_at"`
}

type ListProfileModerationRequestsResponse struct {
	Items     []ProfileModerationRequestResponse `json:"items"`
	PageCount int                                `json:"page_count"`
}

type AdminUserStatsResponse struct {
	ReviewsCount       int `json:"reviews_count"`
	LikesGivenCount    int `json:"likes_given_count"`
	LikesReceivedCount int `json:"likes_received_count"`
}

type AdminUserResponse struct {
	ID               uuid.UUID               `json:"id"`
	Email            string                  `json:"email"`
	Username         string                  `json:"username"`
	Bio              *string                 `json:"bio,omitempty"`
	AvatarURL        *string                 `json:"avatar_url,omitempty"`
	BannerURL        *string                 `json:"banner_url,omitempty"`
	TelegramID       *int64                  `json:"telegram_id,omitempty"`
	TelegramUsername *string                 `json:"telegram_username,omitempty"`
	RoleID           int                     `json:"role_id"`
	IsEmailVerified  bool                    `json:"is_email_verified"`
	IsActive         bool                    `json:"is_active"`
	IsDeleted        bool                    `json:"is_deleted"`
	IsBanned         bool                    `json:"is_banned"`
	BannedByUserID   *uuid.UUID              `json:"banned_by_user_id,omitempty"`
	CreatedAt        string                  `json:"created_at"`
	UpdatedAt        string                  `json:"updated_at"`
	Stats            *AdminUserStatsResponse `json:"stats,omitempty"`
}

type ListAdminUsersResponse struct {
	Items     []AdminUserResponse `json:"items"`
	PageCount int                 `json:"page_count"`
}

type SetUserBanRequest struct {
	IsBanned bool `json:"is_banned"`
}

type SetUserRoleRequest struct {
	RoleID int `json:"role_id" validate:"required,min=1,max=3"`
}

type AdminLogModeratorResponse struct {
	ID       uuid.UUID `json:"id"`
	Username *string   `json:"username,omitempty"`
}

type AdminLogResponse struct {
	ID         int                       `json:"id"`
	Moderator  AdminLogModeratorResponse `json:"moderator"`
	Action     string                    `json:"action"`
	TargetID   string                    `json:"target_id"`
	TargetType string                    `json:"target_type"`
	Details    map[string]any            `json:"details"`
	CreatedAt  string                    `json:"created_at"`
}

type ListAdminLogsResponse struct {
	Items     []AdminLogResponse `json:"items"`
	PageCount int                `json:"page_count"`
}
