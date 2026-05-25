package artist_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

type ArtistDependencies struct {
	ConcertsCount int
	ReviewsCount  int
}

func (d ArtistDependencies) HasAny() bool {
	return d.ConcertsCount > 0 || d.ReviewsCount > 0
}

// DeleteArtistHard — полное физическое удаление артиста из БД
func (s *ArtistService) DeleteArtistHard(ctx context.Context, id int) error {
	if s.artistRepository == nil {
		return fmt.Errorf("artist repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	// Проверяем зависимости (концерты, отзывы и т.д.)
	dependencies, err := s.artistRepository.GetArtistDependencies(ctx, id)
	if err != nil {
		return fmt.Errorf("check artist dependencies: %w", err)
	}
	if dependencies.HasAny() {
		return fmt.Errorf(
			"cannot delete artist with linked entities (concerts: %d, reviews: %d): %w",
			dependencies.ConcertsCount,
			dependencies.ReviewsCount,
			core_errors.ErrConflict,
		)
	}

	// Выполняем физическое удаление
	if err := s.artistRepository.DeleteArtistHard(ctx, id); err != nil {
		return fmt.Errorf("hard delete artist: %w", err)
	}

	if s.statsCache != nil {
		_ = s.statsCache.InvalidateGlobalStats(ctx)
	}

	return nil
}

func (s *ArtistService) DeleteArtistSoft(ctx context.Context, id int) error {
	if s.artistRepository == nil {
		return fmt.Errorf("artist repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	if err := s.artistRepository.DeleteArtistSoft(ctx, id); err != nil {
		return fmt.Errorf("soft delete artist: %w", err)
	}

	if s.statsCache != nil {
		_ = s.statsCache.InvalidateGlobalStats(ctx)
	}

	return nil
}

func (s *ArtistService) RestoreArtist(ctx context.Context, id int) (domain.Artist, error) {
	if s.artistRepository == nil {
		return domain.Artist{}, fmt.Errorf("artist repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	restored, err := s.artistRepository.RestoreArtist(ctx, id)
	if err != nil {
		return domain.Artist{}, fmt.Errorf("restore artist: %w", err)
	}

	if s.statsCache != nil {
		_ = s.statsCache.InvalidateGlobalStats(ctx)
	}

	return restored, nil
}
