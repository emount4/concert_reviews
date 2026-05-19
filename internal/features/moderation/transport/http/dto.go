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
