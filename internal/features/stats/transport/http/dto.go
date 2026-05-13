package stats_transport_http

import "github.com/emount4/concert_reviews/internal/core/domain"

type GlobalStatsResponse struct {
	UsersCount    int `json:"users_count"`
	ConcertsCount int `json:"concerts_count"`
	ArtistsCount  int `json:"artists_count"`
	VenuesCount   int `json:"venues_count"`
	ReviewsCount  int `json:"reviews_count"`
}

func MapDomainToGlobalStatsResponse(s domain.GlobalStats) GlobalStatsResponse {
	return GlobalStatsResponse(s)
}
