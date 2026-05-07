package venue_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

func (s *VenueService) GetVenueByID(ctx context.Context, id int) (domain.Venue, error) {
	if s.venueRepository == nil {
		return domain.Venue{}, fmt.Errorf("venue repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	venue, err := s.venueRepository.GetVenueByID(ctx, id)
	if err != nil {
		return domain.Venue{}, fmt.Errorf("get venue by id: %w", err)
	}

	return venue, nil
}

func (s *VenueService) GetVenues(
	ctx context.Context,
	cityID *int,
	search string,
	limit, offset *int,
) ([]domain.Venue, error) {
	if s.venueRepository == nil {
		return nil, fmt.Errorf("venue repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	venues, err := s.venueRepository.GetVenues(ctx, cityID, search, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get venues: %w", err)
	}

	return venues, nil
}

func (s *VenueService) GetVenuesAdmin(
	ctx context.Context,
	cityID *int,
	search string,
	limit, offset *int,
	includeDeleted bool,
	status string,
) ([]domain.Venue, error) {
	if s.venueRepository == nil {
		return nil, fmt.Errorf("venue repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	// Валидация статуса
	if status != "" && !isValidContentStatus(status) {
		return nil, fmt.Errorf("invalid status filter: %w", core_errors.ErrInvalidArgument)
	}

	venues, err := s.venueRepository.GetVenuesAdmin(ctx, cityID, search, limit, offset, includeDeleted, status)
	if err != nil {
		return nil, fmt.Errorf("get venues admin: %w", err)
	}

	return venues, nil
}

func isValidContentStatus(s string) bool {
	switch domain.ContentStatus(s) {
	case domain.StatusActive, domain.StatusHidden, domain.StatusArchived:
		return true
	}
	return false
}
