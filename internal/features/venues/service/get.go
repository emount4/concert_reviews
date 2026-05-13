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
	sort string,
	direction string,
	capacityFrom, capacityTo *int,
	limit, offset *int,
) ([]domain.Venue, int, error) {
	if s.venueRepository == nil {
		return nil, 0, fmt.Errorf("venue repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	if limit != nil && *limit < 0 {
		return nil, 0, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, 0, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if capacityFrom != nil && *capacityFrom < 0 {
		return nil, 0, fmt.Errorf("capacity_from must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if capacityTo != nil && *capacityTo < 0 {
		return nil, 0, fmt.Errorf("capacity_to must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if capacityFrom != nil && capacityTo != nil && *capacityFrom > *capacityTo {
		return nil, 0, fmt.Errorf("capacity_from must be <= capacity_to: %w", core_errors.ErrInvalidArgument)
	}

	venues, total, err := s.venueRepository.GetVenues(ctx, cityID, search, sort, direction, capacityFrom, capacityTo, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get venues: %w", err)
	}

	return venues, total, nil
}

func (s *VenueService) GetVenuesAdmin(
	ctx context.Context,
	cityID *int,
	search string,
	sort string,
	direction string,
	capacityFrom, capacityTo *int,
	limit, offset *int,
	includeDeleted bool,
	status string,
) ([]domain.Venue, int, error) {
	if s.venueRepository == nil {
		return nil, 0, fmt.Errorf("venue repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	if status != "" && !isValidContentStatus(status) {
		return nil, 0, fmt.Errorf("invalid status filter: %w", core_errors.ErrInvalidArgument)
	}
	if capacityFrom != nil && *capacityFrom < 0 {
		return nil, 0, fmt.Errorf("capacity_from must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if capacityTo != nil && *capacityTo < 0 {
		return nil, 0, fmt.Errorf("capacity_to must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if capacityFrom != nil && capacityTo != nil && *capacityFrom > *capacityTo {
		return nil, 0, fmt.Errorf("capacity_from must be <= capacity_to: %w", core_errors.ErrInvalidArgument)
	}

	venues, total, err := s.venueRepository.GetVenuesAdmin(ctx, cityID, search, sort, direction, capacityFrom, capacityTo, limit, offset, includeDeleted, status)
	if err != nil {
		return nil, 0, fmt.Errorf("get venues admin: %w", err)
	}

	return venues, total, nil
}

func isValidContentStatus(s string) bool {
	switch domain.ContentStatus(s) {
	case domain.StatusActive, domain.StatusHidden, domain.StatusArchived:
		return true
	}
	return false
}
