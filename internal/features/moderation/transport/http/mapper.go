package moderation_transport_http

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
)

func MapProfileModerationRequestsToResponse(items []domain.ProfileModerationRequest) ListProfileModerationRequestsResponse {
	respItems := make([]ProfileModerationRequestResponse, len(items))
	for i, item := range items {
		respItems[i] = MapProfileModerationRequestToResponse(item)
	}
	return ListProfileModerationRequestsResponse{Items: respItems}
}

func MapProfileModerationRequestToResponse(item domain.ProfileModerationRequest) ProfileModerationRequestResponse {
	return ProfileModerationRequestResponse{
		ID: item.ModerationID,
		User: ProfileModerationUserResponse{
			ID:        item.UserID,
			Username:  item.Username,
			AvatarURL: item.UserAvatarURL,
		},
		FieldName:         mapProfileModerationFieldName(item.FieldName),
		OldValue:          item.OldValue,
		NewValue:          item.NewValue,
		Status:            string(item.Status),
		ModeratedByUserID: item.ModeratedByUserID,
		CreatedAt:         item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         item.UpdatedAt.Format(time.RFC3339),
	}
}

func mapProfileModerationFieldName(fieldName string) string {
	switch fieldName {
	case "avatar_url":
		return "avatar_key"
	case "banner_url":
		return "banner_key"
	default:
		return fieldName
	}
}
