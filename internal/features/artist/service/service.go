package artist_service

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_ports "github.com/emount4/concert_reviews/internal/core/domain/ports"
)

type ArtistService struct {
	artistRepository ArtistRepository
	s3               core_ports.S3Provider
	statsCache       GlobalStatsCache
}

type ArtistRepository interface {
	CreateArtist(ctx context.Context, artist domain.Artist) (domain.Artist, error)
	GetArtists(ctx context.Context, search string, sort string, direction string, hasReviews *bool, limit, offset *int) ([]domain.Artist, int, error)
	GetArtistByID(ctx context.Context, id int) (domain.Artist, error)
	PatchArtist(ctx context.Context, id int, artist domain.Artist) (domain.Artist, error)
	DeleteArtistHard(ctx context.Context, id int) error
	DeleteArtistSoft(ctx context.Context, id int) error

	GetArtistDependencies(ctx context.Context, id int) (ArtistDependencies, error)
	RestoreArtist(ctx context.Context, id int) (domain.Artist, error)

	GetArtistsAdmin(
		ctx context.Context,
		search string,
		sort string,
		direction string,
		hasReviews *bool,
		limit, offset *int,
		includeDeleted bool,
		status string, // "", "active", "hidden", "archived"
	) ([]domain.Artist, int, error)
}

type GlobalStatsCache interface {
	InvalidateGlobalStats(ctx context.Context) error
}

func NewArtistService(
	artistRepository ArtistRepository,
	s3 core_ports.S3Provider,
	statsCache ...GlobalStatsCache,
) *ArtistService {
	var cache GlobalStatsCache
	if len(statsCache) > 0 {
		cache = statsCache[0]
	}

	return &ArtistService{
		artistRepository: artistRepository,
		s3:               s3,
		statsCache:       cache,
	}
}
