package venue_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

// DeleteVenueHard — проверка зависимостей → физическое удаление
func (s *VenueService) DeleteVenueHard(ctx context.Context, id int) error {
	if s.venueRepository == nil {
		return fmt.Errorf("venue repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	deps, err := s.venueRepository.GetVenueDependencies(ctx, id)
	if err != nil {
		return fmt.Errorf("check venue dependencies: %w", err)
	}
	if deps.HasAny() {
		return fmt.Errorf(
			"cannot delete venue with linked entities (concerts: %d, reviews: %d): %w",
			deps.ConcertsCount, deps.ReviewsCount, core_errors.ErrConflict,
		)
	}

	if err := s.venueRepository.DeleteVenueHard(ctx, id); err != nil {
		return fmt.Errorf("hard delete venue: %w", err)
	}
	if s.statsCache != nil {
		_ = s.statsCache.InvalidateGlobalStats(ctx)
	}
	return nil
}

// DeleteVenueSoft — установка deleted_at
func (s *VenueService) DeleteVenueSoft(ctx context.Context, id int) error {
	if s.venueRepository == nil {
		return fmt.Errorf("venue repository: %w", core_errors.ErrRepositoryNotConfigured)
	}
	if err := s.venueRepository.DeleteVenueSoft(ctx, id); err != nil {
		return fmt.Errorf("soft delete venue: %w", err)
	}
	if s.statsCache != nil {
		_ = s.statsCache.InvalidateGlobalStats(ctx)
	}
	return nil
}

// RestoreVenue — снятие метки удаления
func (s *VenueService) RestoreVenue(ctx context.Context, id int) (domain.Venue, error) {
	if s.venueRepository == nil {
		return domain.Venue{}, fmt.Errorf("venue repository: %w", core_errors.ErrRepositoryNotConfigured)
	}
	venue, err := s.venueRepository.RestoreVenue(ctx, id)
	if err != nil {
		return domain.Venue{}, fmt.Errorf("restore venue: %w", err)
	}
	if s.statsCache != nil {
		_ = s.statsCache.InvalidateGlobalStats(ctx)
	}
	return venue, nil
}
