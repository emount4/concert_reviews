package concert_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ConcertService) GetConcerts(
	ctx context.Context,
	cityID *int,
	artistID *int,
	search string,
	limit, offset *int,
) ([]domain.Concert, error) {
	if s.concertRepository == nil {
		return nil, core_errors.ErrRepositoryNotConfigured
	}

	if cityID != nil && *cityID <= 0 {
		return nil, fmt.Errorf("invalid city_id: %w", core_errors.ErrInvalidArgument)
	}
	if artistID != nil && *artistID <= 0 {
		return nil, fmt.Errorf("invalid artist_id: %w", core_errors.ErrInvalidArgument)
	}
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	return s.concertRepository.GetConcerts(ctx, cityID, artistID, search, limit, offset)
}

func (s *ConcertService) GetConcertByID(ctx context.Context, id uuid.UUID) (domain.Concert, error) {
	if id == uuid.Nil {
		return domain.Concert{}, fmt.Errorf("invalid uuid: %w", core_errors.ErrInvalidArgument)
	}

	concert, err := s.concertRepository.GetConcertByID(ctx, id)
	if err != nil {
		return domain.Concert{}, err
	}

	return concert, nil
}

func (s *ConcertService) GetConcertsAdmin(
	ctx context.Context,
	cityID *int,
	artistID *int,
	search string,
	limit, offset *int,
	includeDeleted bool,
) ([]domain.Concert, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	return s.concertRepository.GetConcertsAdmin(ctx, cityID, artistID, search, limit, offset, includeDeleted)
}

func (s *ConcertService) GetSuggestionsAdmin(ctx context.Context, limit, offset *int, status string) ([]domain.ConcertSuggestion, error) {
	return s.concertRepository.GetSuggestionsAdmin(ctx, limit, offset, status)
}

func (s *ConcertService) GetSuggestionByIDAdmin(ctx context.Context, id uuid.UUID) (domain.ConcertSuggestion, error) {
	if id == uuid.Nil {
		return domain.ConcertSuggestion{}, core_errors.ErrInvalidArgument
	}
	return s.concertRepository.GetSuggestionByIDAdmin(ctx, id)
}
