package venue_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

func (s *VenueService) UpdateVenue(
	ctx context.Context,
	id int,
	patch domain.VenuePatch,
) (domain.Venue, error) {
	if s.venueRepository == nil {
		return domain.Venue{}, fmt.Errorf("venue repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	if patch.CityID.Set && patch.CityID.Value != nil {
		exists, err := s.venueRepository.CityExists(ctx, *patch.CityID.Value)
		if err != nil {
			return domain.Venue{}, fmt.Errorf("check city existence: %w", err)
		}
		if !exists {
			return domain.Venue{}, fmt.Errorf("city not found: %w", core_errors.ErrInvalidArgument)
		}
	}

	if patch.PhotoKey.Set && patch.PhotoKey.Value != nil {
		if _, err := s.s3.FileExists(ctx, *patch.PhotoKey.Value); err != nil {
			return domain.Venue{}, fmt.Errorf("venue is not exist: %w", err)
		}
	}

	updatedVenue, err := s.venueRepository.PatchVenue(ctx, id, patch)
	if err != nil {
		return domain.Venue{}, fmt.Errorf("patch venue repository: %w", err)
	}

	return updatedVenue, nil
}
