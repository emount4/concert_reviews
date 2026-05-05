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
func (s *ArtistService) GetArtists(ctx context.Context, search string, limit, offset *int) ([]domain.Artist, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	artists, err := s.artistRepository.GetArtists(ctx, search, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get artists: %w", err)
	}

	return artists, nil
}
