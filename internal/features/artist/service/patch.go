package artist_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
)

func (s *ArtistService) PatchArtist(
	ctx context.Context,
	id int,
	patch domain.ArtistPatch,
) (domain.Artist, error) {
	artist, err := s.artistRepository.GetArtistByID(ctx, id)
	if err != nil {
		return domain.Artist{}, fmt.Errorf("get artist: %w", err)
	}

	err = artist.ApplyPatch(patch)
	if patch.PhotoKey.Set && patch.PhotoKey.Value != nil {
		if _, err := s.s3.FileExists(ctx, *patch.PhotoKey.Value); err != nil {
			return domain.Artist{}, fmt.Errorf("file is not exist: %w", err)
		}
	}
	if err != nil {
		return domain.Artist{}, fmt.Errorf("apply artist patch: %w", err)
	}

	domainArtist, err := s.artistRepository.PatchArtist(ctx, id, artist)

	if err != nil {
		return domain.Artist{}, fmt.Errorf("patch artist: %w", err)
	}
	return domainArtist, nil
}
