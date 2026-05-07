package domain

import (
	"fmt"
	"strings"
	"time"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_types "github.com/emount4/concert_reviews/internal/core/types"
	"github.com/google/uuid"
)

type Concert struct {
	ConcertID       uuid.UUID
	VenueID         int
	Title           string
	Date            time.Time
	PosterURL       *string
	IsVerified      bool
	CreatedByUserID *uuid.UUID
	Status          ContentStatus
	CreatedAt       time.Time
	DeletedAt       *time.Time
}

func (c *Concert) Validate() error {
	if c.VenueID <= 0 {
		return fmt.Errorf("venue_id must be positive")
	}

	if strings.TrimSpace(c.Title) == "" {
		return fmt.Errorf("concert title cannot be empty")
	}
	if len(c.Title) > 255 {
		return fmt.Errorf("title too long (max 255)")
	}

	if c.Date.IsZero() {
		return fmt.Errorf("date is required")
	}

	if c.PosterURL != nil {
		if strings.TrimSpace(*c.PosterURL) == "" {
			return fmt.Errorf("poster_url cannot be empty string if provided")
		}
		if len(*c.PosterURL) > 2048 {
			return fmt.Errorf("poster_url too long (max 2048)")
		}
	}

	return nil
}

func (c *Concert) IsDeleted() bool {
	return c.DeletedAt != nil
}

func (c *Concert) IsUpcoming() bool {
	return c.Date.After(time.Now())
}

type ConcertPatch struct {
	VenueID    core_types.Nullable[int]       `json:"venue_id"`
	Title      core_types.Nullable[string]    `json:"title"`
	Date       core_types.Nullable[time.Time] `json:"date"`
	PosterKey  core_types.Nullable[string]    `json:"poster_key"`
	IsVerified core_types.Nullable[bool]      `json:"is_verified"`
	Status     core_types.Nullable[string]    `json:"status"`
}

func (p *ConcertPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("title cannot be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if p.VenueID.Set {
		if p.VenueID.Value == nil {
			return fmt.Errorf("venue_id cannot be patched to NULL: %w", core_errors.ErrInvalidArgument)
		}
		if *p.VenueID.Value <= 0 {
			return fmt.Errorf("venue_id must be positive: %w", core_errors.ErrInvalidArgument)
		}
	}

	if p.Status.Set && p.Status.Value != nil {
		switch ContentStatus(*p.Status.Value) {
		case StatusActive, StatusHidden, StatusArchived:
			// OK
		default:
			return fmt.Errorf("invalid status value: %w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

// ApplyPatch применяет изменения из патча к доменной модели
func (c *Concert) ApplyPatch(patch ConcertPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate concert patch: %w", err)
	}

	tmp := *c

	if patch.VenueID.Set && patch.VenueID.Value != nil {
		tmp.VenueID = *patch.VenueID.Value
	}
	if patch.Title.Set && patch.Title.Value != nil {
		tmp.Title = *patch.Title.Value
	}
	if patch.Date.Set && !patch.Date.Value.IsZero() {
		tmp.Date = *patch.Date.Value
	}
	if patch.PosterKey.Set {
		tmp.PosterURL = patch.PosterKey.Value
	}
	if patch.IsVerified.Set {
		tmp.IsVerified = *patch.IsVerified.Value
	}
	if patch.Status.Set && patch.Status.Value != nil {
		tmp.Status = ContentStatus(*patch.Status.Value)
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched concert: %w", err)
	}

	*c = tmp
	return nil
}
