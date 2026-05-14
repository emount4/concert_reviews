package concert_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ConcertService) UpdateConcert(
	ctx context.Context,
	id uuid.UUID,
	patch domain.ConcertPatch,
) (domain.Concert, error) {
	if s.concertRepository == nil {
		return domain.Concert{}, core_errors.ErrRepositoryNotConfigured
	}

	concert, err := s.concertRepository.GetConcertByID(ctx, id, uuid.Nil)
	if err != nil {
		return domain.Concert{}, fmt.Errorf("fetch concert for update: %w", err)
	}

	// Validate patch and check S3 file before applying
	if err := patch.Validate(); err != nil {
		return domain.Concert{}, fmt.Errorf("validate patch: %w", err)
	}

	if patch.PosterKey.Set && patch.PosterKey.Value != nil {
		posterKey := *patch.PosterKey.Value
		if posterKey != "" {
			if _, err := s.s3.FileExists(ctx, posterKey); err != nil {
				return domain.Concert{}, fmt.Errorf("patched poster not found in S3: %w", err)
			}
		}
	}

	if err := concert.ApplyPatch(patch); err != nil {
		return domain.Concert{}, fmt.Errorf("apply patch to domain: %w", err)
	}

	updated, err := s.concertRepository.UpdateConcert(ctx, id, patch)
	if err != nil {
		return domain.Concert{}, fmt.Errorf("update concert in database: %w", err)
	}

	return updated, nil
}
