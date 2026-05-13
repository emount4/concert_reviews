package artist_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

// GetArtistByID возвращает доменную модель по int ID.
func (s *ArtistService) GetArtistByID(ctx context.Context, id int) (domain.Artist, error) {

	artist, err := s.artistRepository.GetArtistByID(ctx, id)
	if err != nil {
		return domain.Artist{}, fmt.Errorf("get artist by id: %w", err)
	}

	return artist, nil
}

// GetArtists возвращает срез доменных моделей с учетом пагинации.
func (s *ArtistService) GetArtists(ctx context.Context, search string, sort string, direction string, hasReviews *bool, limit, offset *int) ([]domain.Artist, int, error) {
	if s.artistRepository == nil {
		return nil, 0, core_errors.ErrRepositoryNotConfigured
	}
	if limit != nil && *limit < 0 {
		return nil, 0, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, 0, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	artists, total, err := s.artistRepository.GetArtists(ctx, search, sort, direction, hasReviews, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get artists: %w", err)
	}

	return artists, total, nil
}

// internal/core/service/artist/artist_service.go

// GetArtistsAdmin — получение списка артистов для админки с фильтрами
func (s *ArtistService) GetArtistsAdmin(
	ctx context.Context,
	search string,
	sort string,
	direction string,
	hasReviews *bool,
	limit, offset *int,
	includeDeleted bool,
	status string,
) ([]domain.Artist, int, error) {
	if s.artistRepository == nil {
		return nil, 0, fmt.Errorf("artist repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	// Валидация пагинации
	if limit != nil && *limit < 0 {
		return nil, 0, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, 0, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	// Валидация статуса (если передан)
	if status != "" && !isValidContentStatus(status) {
		return nil, 0, fmt.Errorf("invalid status filter: %w", core_errors.ErrInvalidArgument)
	}

	artists, total, err := s.artistRepository.GetArtistsAdmin(
		ctx,
		search,
		sort,
		direction,
		hasReviews,
		limit,
		offset,
		includeDeleted,
		status,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("get artists admin from repository: %w", err)
	}

	return artists, total, nil
}

// Вспомогательная функция для валидации статуса
func isValidContentStatus(s string) bool {
	switch domain.ContentStatus(s) {
	case domain.StatusActive, domain.StatusHidden, domain.StatusArchived:
		return true
	}
	return false
}
