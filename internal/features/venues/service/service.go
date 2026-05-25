package venue_service

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_ports "github.com/emount4/concert_reviews/internal/core/domain/ports"
)

type VenueService struct {
	venueRepository VenueRepository
	s3              core_ports.S3Provider
	statsCache      GlobalStatsCache
}

type VenueRepository interface {
	CreateVenue(ctx context.Context, venue domain.Venue) (domain.Venue, error)
	GetVenues(
		ctx context.Context,
		cityID *int,
		search string,
		sort string,
		direction string,
		capacityFrom, capacityTo *int,
		limit, offset *int,
	) ([]domain.Venue, int, error)
	GetVenueByID(ctx context.Context, id int) (domain.Venue, error)
	PatchVenue(ctx context.Context, id int, patch domain.VenuePatch) (domain.Venue, error)
	DeleteVenueHard(ctx context.Context, id int) error
	DeleteVenueSoft(ctx context.Context, id int) error

	GetVenueDependencies(ctx context.Context, id int) (domain.VenueDependencies, error)
	RestoreVenue(ctx context.Context, id int) (domain.Venue, error)

	GetVenuesAdmin(
		ctx context.Context,
		cityID *int,
		search string,
		sort string,
		direction string,
		capacityFrom, capacityTo *int,
		limit, offset *int,
		includeDeleted bool,
		status string,
	) ([]domain.Venue, int, error)

	CityExists(ctx context.Context, id int) (bool, error)
}

type GlobalStatsCache interface {
	InvalidateGlobalStats(ctx context.Context) error
}

func NewVenueService(
	venueRepository VenueRepository,
	s3 core_ports.S3Provider,
	statsCache ...GlobalStatsCache,
) *VenueService {
	var cache GlobalStatsCache
	if len(statsCache) > 0 {
		cache = statsCache[0]
	}

	return &VenueService{
		venueRepository: venueRepository,
		s3:              s3,
		statsCache:      cache,
	}
}
