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

func MapAdminUsersToResponse(users []domain.User) ListAdminUsersResponse {
	items := make([]AdminUserResponse, len(users))
	for i, user := range users {
		items[i] = MapAdminUserToResponse(user)
	}
	return ListAdminUsersResponse{Items: items}
}

func MapAdminUserToResponse(user domain.User) AdminUserResponse {
	resp := AdminUserResponse{
		ID:               user.ID,
		Email:            user.Email,
		Username:         user.Username,
		Bio:              user.Bio,
		AvatarURL:        user.AvatarURL,
		BannerURL:        user.BannerURL,
		TelegramID:       user.TelegramID,
		TelegramUsername: user.TelegramUsername,
		RoleID:           user.RoleID,
		IsEmailVerified:  user.IsEmailVerified,
		IsActive:         user.IsActive,
		IsDeleted:        !user.IsActive,
		IsBanned:         user.IsBanned,
		BannedByUserID:   user.BannedByUserID,
		CreatedAt:        user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        user.UpdatedAt.Format(time.RFC3339),
	}

	if user.Stats != nil {
		resp.Stats = &AdminUserStatsResponse{
			ReviewsCount:       user.Stats.ReviewsCount,
			LikesGivenCount:    user.Stats.LikesGivenCount,
			LikesReceivedCount: user.Stats.LikesReceivedCount,
		}
	}

	return resp
}

func MapAdminLogsToResponse(logs []domain.AdminLog) ListAdminLogsResponse {
	items := make([]AdminLogResponse, len(logs))
	for i, item := range logs {
		items[i] = MapAdminLogToResponse(item)
	}
	return ListAdminLogsResponse{Items: items}
}

func MapAdminLogToResponse(item domain.AdminLog) AdminLogResponse {
	return AdminLogResponse{
		ID: item.LogID,
		Moderator: AdminLogModeratorResponse{
			ID:       item.ModeratorID,
			Username: item.ModeratorUsername,
		},
		Action:     item.Action,
		TargetID:   item.TargetID,
		TargetType: item.TargetType,
		Details:    item.Details,
		CreatedAt:  item.CreatedAt.Format(time.RFC3339),
	}
}
