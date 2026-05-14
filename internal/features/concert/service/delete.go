package concert_service

import (
	"context"
	"fmt"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ConcertService) DeleteConcertSoft(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("concert id is nil: %w", core_errors.ErrInvalidArgument)
	}

	if err := s.concertRepository.DeleteConcertSoft(ctx, id); err != nil {
		return fmt.Errorf("soft delete concert %s in repository: %w", id, err)
	}

	return nil
}

func (s *ConcertService) DeleteConcertHard(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("concert id is nil: %w", core_errors.ErrInvalidArgument)
	}

	concert, err := s.concertRepository.GetConcertByID(ctx, id, uuid.Nil)
	if err != nil {
		return fmt.Errorf("fetch concert %s before hard delete: %w", id, err)
	}

	if err := s.concertRepository.DeleteConcertHard(ctx, id); err != nil {
		return fmt.Errorf("hard delete concert %s from database: %w", id, err)
	}

	if concert.PosterKey != nil && *concert.PosterKey != "" {
		if err := s.s3.DeleteObject(ctx, *concert.PosterKey); err != nil {
			return fmt.Errorf("cleanup concert poster %s from S3: %w", *concert.PosterKey, err)
		}
	}

	return nil
}

func (s *ConcertService) DeleteSuggestionAdmin(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("suggestion id is nil: %w", core_errors.ErrInvalidArgument)
	}

	_, err := s.concertRepository.GetSuggestionByIDAdmin(ctx, id)
	if err != nil {
		return fmt.Errorf("verify suggestion %s existence: %w", id, err)
	}

	if err := s.concertRepository.DeleteSuggestionAdmin(ctx, id); err != nil {
		return fmt.Errorf("remove suggestion %s from repository: %w", id, err)
	}

	return nil
}
