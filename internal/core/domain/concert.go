package domain

import (
	"fmt"
	"strings"
	"time"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_types "github.com/emount4/concert_reviews/internal/core/types"
	"github.com/google/uuid"
)

// ConcertArtist описывает связь артиста с конкретным концертом
type ConcertArtist struct {
	ArtistID int
	Name     string // Подгружаем для удобства отображения
	IsMain   bool
}

type ConcertStats struct {
	ReviewsCount   int
	SumP1          int
	SumP2          int
	SumP3          int
	SumP4          int
	SumP5          int
	SumRatingTotal int64
	FavoritesCount int
	UpdatedAt      time.Time
}

// AvgRating — расчет среднего общего рейтинга
func (s *ConcertStats) AvgRating() float64 {
	if s == nil || s.ReviewsCount == 0 {
		return 0
	}
	return float64(s.SumRatingTotal) / float64(s.ReviewsCount)
}

// AvgByParam — расчет среднего по конкретному параметру (1-5)
func (s *ConcertStats) AvgByParam(param int) float64 {
	if s == nil || s.ReviewsCount == 0 {
		return 0
	}
	var sum int
	switch param {
	case 1:
		sum = s.SumP1
	case 2:
		sum = s.SumP2
	case 3:
		sum = s.SumP3
	case 4:
		sum = s.SumP4
	case 5:
		sum = s.SumP5
	default:
		return 0
	}
	return float64(sum) / float64(s.ReviewsCount)
}

type Concert struct {
	ConcertID       uuid.UUID
	VenueID         int
	Title           string
	Date            time.Time
	PosterKey       *string
	IsVerified      bool
	CreatedByUserID *uuid.UUID
	CreatedAt       time.Time
	DeletedAt       *time.Time

	Stats   *ConcertStats
	Venue   *Venue
	Artists []ConcertArtist
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
	if c.PosterKey != nil {
		if strings.TrimSpace(*c.PosterKey) == "" {
			return fmt.Errorf("poster_key cannot be empty string if provided")
		}
		if len(*c.PosterKey) > 2048 {
			return fmt.Errorf("poster_key too long (max 2048)")
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
	VenueID    core_types.Nullable[int]
	Title      core_types.Nullable[string]
	Date       core_types.Nullable[time.Time]
	PosterKey  core_types.Nullable[string]
	IsVerified core_types.Nullable[bool]
}

func (p *ConcertPatch) Validate() error {
	if p.Title.Set && (p.Title.Value == nil || strings.TrimSpace(*p.Title.Value) == "") {
		return fmt.Errorf("title cannot be empty: %w", core_errors.ErrInvalidArgument)
	}
	if p.VenueID.Set {
		if p.VenueID.Value == nil || *p.VenueID.Value <= 0 {
			return fmt.Errorf("invalid venue_id: %w", core_errors.ErrInvalidArgument)
		}
	}
	return nil
}

func (c *Concert) ApplyPatch(patch ConcertPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate concert patch: %w", err)
	}

	if patch.VenueID.Set {
		c.VenueID = *patch.VenueID.Value
	}
	if patch.Title.Set {
		c.Title = *patch.Title.Value
	}
	if patch.Date.Set {
		c.Date = *patch.Date.Value
	}
	if patch.PosterKey.Set {
		c.PosterKey = patch.PosterKey.Value
	}
	if patch.IsVerified.Set {
		c.IsVerified = *patch.IsVerified.Value
	}

	return c.Validate()
}
