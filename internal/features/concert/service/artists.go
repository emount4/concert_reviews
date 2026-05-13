package concert_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// AddArtistToConcert adds an artist to a concert
func (s *ConcertService) AddArtistToConcert(ctx context.Context, concertID uuid.UUID, artistID int, isMain bool) error {
	if err := s.concertRepository.AddConcertArtist(ctx, concertID, artistID, isMain); err != nil {
		return fmt.Errorf("add artist to concert: %w", err)
	}
	return nil
}

// RemoveArtistFromConcert removes an artist from a concert
func (s *ConcertService) RemoveArtistFromConcert(ctx context.Context, concertID uuid.UUID, artistID int) error {
	if err := s.concertRepository.RemoveConcertArtist(ctx, concertID, artistID); err != nil {
		return fmt.Errorf("remove artist from concert: %w", err)
	}
	return nil
}

// UpdateArtistMainStatus updates the is_main flag for an artist in a concert
func (s *ConcertService) UpdateArtistMainStatus(ctx context.Context, concertID uuid.UUID, artistID int, isMain bool) error {
	if err := s.concertRepository.UpdateConcertArtistIsMain(ctx, concertID, artistID, isMain); err != nil {
		return fmt.Errorf("update artist main status: %w", err)
	}
	return nil
}
