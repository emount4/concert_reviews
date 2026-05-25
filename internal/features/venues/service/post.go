package venue_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

func (s *VenueService) CreateVenue(ctx context.Context, venue domain.Venue) (domain.Venue, error) {
	if s.venueRepository == nil {
		return domain.Venue{}, fmt.Errorf("venue repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	if err := venue.Validate(); err != nil {
		return domain.Venue{}, fmt.Errorf("validate venue: %w", err)
	}

	if venue.PhotoURL != nil {
		_, err := s.s3.FileExists(ctx, *venue.PhotoURL)
		if err != nil {
			return domain.Venue{}, fmt.Errorf("failed to find photo: %w", err)
		}
	}

	if ok, err := s.venueRepository.CityExists(ctx, venue.CityID); !ok || err != nil {
		return domain.Venue{}, fmt.Errorf("city not found: %w", core_errors.ErrInvalidArgument)
	}

	createdVenue, err := s.venueRepository.CreateVenue(ctx, venue)
	if err != nil {
		return domain.Venue{}, fmt.Errorf("create venue repository: %w", err)
	}

	if s.statsCache != nil {
		_ = s.statsCache.InvalidateGlobalStats(ctx)
	}

	return createdVenue, nil
}
