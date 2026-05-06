package artist_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
)

func (s *ArtistService) CreateArtist(ctx context.Context, artist domain.Artist) (domain.Artist, error) {
	err := artist.Validate()

	if err != nil {
		return domain.Artist{}, fmt.Errorf("validate artist: %w", err)
	}

	if artist.PhotoURL != nil {
		if _, err := s.s3.FileExists(ctx, *artist.PhotoURL); err != nil {
			return domain.Artist{}, fmt.Errorf("artist image not found: %w", err)
		}
	}

	artist.Status = domain.StatusActive

	createdArtist, err := s.artistRepository.CreateArtist(ctx, artist)
	if err != nil {
		return domain.Artist{}, fmt.Errorf("create artist in repo: %w", err)
	}

	return createdArtist, nil
}
