package user_transport_http

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
)

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
