package concert_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ConcertService) CreateConcert(
	ctx context.Context,
	concert domain.Concert,
	artists []domain.ConcertArtist,
) (domain.Concert, error) {

	if err := concert.Validate(); err != nil {
		return domain.Concert{}, fmt.Errorf("validate concert: %w", err)
	}

	if concert.PosterKey != nil {
		if _, err := s.s3.FileExists(ctx, *concert.PosterKey); err != nil {
			return domain.Concert{}, fmt.Errorf("concert poster not found in S3: %w", err)
		}
	}

	created, err := s.concertRepository.CreateConcert(ctx, concert, artists)
	if err != nil {
		return domain.Concert{}, fmt.Errorf("create concert in repository: %w", err)
	}

	if s.statsCache != nil {
		_ = s.statsCache.InvalidateGlobalStats(ctx)
	}

	return created, nil
}

func (s *ConcertService) SuggestConcert(ctx context.Context, suggestion domain.ConcertSuggestion) (domain.ConcertSuggestion, error) {
	if err := suggestion.Validate(); err != nil {
		return domain.ConcertSuggestion{}, fmt.Errorf("validate suggestion: %w", err)
	}

	pendingCount, err := s.concertRepository.CountPendingSuggestions(ctx, suggestion.UserID)
	if err != nil {
		return domain.ConcertSuggestion{}, fmt.Errorf("check user pending suggestions limit: %w", err)
	}

	if pendingCount >= 3 {
		return domain.ConcertSuggestion{}, fmt.Errorf(
			"user %s reached limit of 3 pending suggestions: %w",
			suggestion.UserID,
			core_errors.ErrConflict,
		)
	}

	return s.concertRepository.CreateSuggestion(ctx, suggestion)
}

func (s *ConcertService) RestoreConcert(ctx context.Context, id uuid.UUID) (domain.Concert, error) {
	if s.concertRepository == nil {
		return domain.Concert{}, fmt.Errorf("concert repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	if id == uuid.Nil {
		return domain.Concert{}, fmt.Errorf("invalid concert uuid: %w", core_errors.ErrInvalidArgument)
	}

	restored, err := s.concertRepository.RestoreConcert(ctx, id)
	if err != nil {
		return domain.Concert{}, fmt.Errorf("restore concert in repository: %w", err)
	}

	if s.statsCache != nil {
		_ = s.statsCache.InvalidateGlobalStats(ctx)
	}

	return restored, nil
}
