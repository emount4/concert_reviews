package user_transport_http

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
)

func MapUpdateUserRequestToDomain(req UpdateUserRequest) domain.UserPatch {
	return domain.UserPatch{
		Username:  req.Username.ToDomain(),
		Bio:       req.Bio.ToDomain(),
		AvatarKey: req.AvatarKey.ToDomain(),
		BannerKey: req.BannerKey.ToDomain(),
	}
}

func MapDomainToMeResponse(u domain.User) UserMeResponse {
	resp := UserMeResponse{
		ID:               u.ID.String(),
		Email:            u.Email,
		Username:         u.Username,
		Bio:              u.Bio,
		AvatarURL:        u.AvatarURL,
		BannerURL:        u.BannerURL,
		TelegramID:       u.TelegramID,
		TelegramUsername: u.TelegramUsername,
		RoleID:           u.RoleID,
		IsEmailVerified:  u.IsEmailVerified,
		IsBanned:         u.IsBanned,
		CreatedAt:        u.CreatedAt.Format(time.RFC3339),
	}

	if u.Stats != nil {
		resp.Stats = &UserStatsDTO{
			ReviewsCount:       u.Stats.ReviewsCount,
			LikesGivenCount:    u.Stats.LikesGivenCount,
			LikesReceivedCount: u.Stats.LikesReceivedCount,
		}
	}
	return resp
}

func MapDomainToPublicResponse(u domain.User) PublicProfileResponse {
	resp := PublicProfileResponse{
		ID:        u.ID.String(),
		Username:  u.Username,
		Bio:       u.Bio,
		AvatarURL: u.AvatarURL,
		BannerURL: u.BannerURL,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}

	if u.Stats != nil {
		resp.Stats = &UserStatsDTO{
			ReviewsCount:       u.Stats.ReviewsCount,
			LikesGivenCount:    u.Stats.LikesGivenCount,
			LikesReceivedCount: u.Stats.LikesReceivedCount,
		}
	}
	return resp
}

func MapProfileModerationRequestsToResponse(items []domain.ProfileModerationRequest) ListProfileModerationRequestsResponse {
	respItems := make([]ProfileModerationRequestResponse, len(items))
	for i, item := range items {
		respItems[i] = ProfileModerationRequestResponse{
			ID:                item.ModerationID,
			FieldName:         mapProfileModerationFieldName(item.FieldName),
			OldValue:          item.OldValue,
			NewValue:          item.NewValue,
			Status:            string(item.Status),
			ModeratedByUserID: item.ModeratedByUserID,
			CreatedAt:         item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:         item.UpdatedAt.Format(time.RFC3339),
		}
	}
	return ListProfileModerationRequestsResponse{Items: respItems}
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
