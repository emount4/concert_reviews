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

	// Validate patch first
	if err := patch.Validate(); err != nil {
		return domain.Artist{}, fmt.Errorf("validate patch: %w", err)
	}

	// Check S3 file before applying patch
	if patch.PhotoKey.Set && patch.PhotoKey.Value != nil {
		if _, err := s.s3.FileExists(ctx, *patch.PhotoKey.Value); err != nil {
			return domain.Artist{}, fmt.Errorf("photo file not found in S3: %w", err)
		}
	}

	// Apply patch to domain model
	if err := artist.ApplyPatch(patch); err != nil {
		return domain.Artist{}, fmt.Errorf("apply artist patch: %w", err)
	}

	domainArtist, err := s.artistRepository.PatchArtist(ctx, id, artist)

	if err != nil {
		return domain.Artist{}, fmt.Errorf("patch artist: %w", err)
	}
	return domainArtist, nil
}
