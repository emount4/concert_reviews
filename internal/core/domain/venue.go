package domain

import (
	"fmt"
	"strings"
	"time"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_http_types "github.com/emount4/concert_reviews/internal/core/transport/http/types"
	core_types "github.com/emount4/concert_reviews/internal/core/types"
)

type City struct {
	CityID int
	Name   string
	Slug   string
}

type Venue struct {
	VenueID     int
	CityID      int
	City        *City
	Name        string
	Address     *string
	Capacity    *int
	SocialLinks map[string]string
	PhotoURL    *string
	Description *string
	Stats       *ContentStats
	Status      ContentStatus
	CreatedAt   time.Time
	DeletedAt   *time.Time
}

func (v *Venue) Validate() error {
	if v.CityID <= 0 {
		return fmt.Errorf("city_id must be positive")
	}

	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("venue name cannot be empty")
	}
	if len(v.Name) > 255 {
		return fmt.Errorf("venue name too long (max 255)")
	}

	if v.Address != nil && strings.TrimSpace(*v.Address) == "" {
		return fmt.Errorf("address cannot be empty string if provided")
	}

	if v.Capacity != nil {
		if *v.Capacity <= 0 || *v.Capacity > 500000 {
			return fmt.Errorf("capacity must be between 1 and 500000")
		}
	}

	if v.PhotoURL != nil {
		if strings.TrimSpace(*v.PhotoURL) == "" {
			return fmt.Errorf("photo_url cannot be empty string if provided")
		}
		if len(*v.PhotoURL) > 2048 {
			return fmt.Errorf("photo_url too long (max 2048)")
		}
	}

	if v.Description != nil && len(*v.Description) > 2000 {
		return fmt.Errorf("description too long (max 2000)")
	}

	return nil
}

func (v *Venue) IsDeleted() bool {
	return v.DeletedAt != nil
}

type VenuePatch struct {
	CityID      core_types.Nullable[int]                `json:"city_id" validate:"omitempty,min=1"`
	Name        core_types.Nullable[string]             `json:"name" validate:"omitempty,min=2,max=255"`
	Address     core_types.Nullable[string]             `json:"address" validate:"omitempty"`
	Capacity    core_types.Nullable[int]                `json:"capacity" validate:"omitempty,min=1,max=500000"`
	SocialLinks core_http_types.NullableMapStringString `json:"social_links"`
	PhotoKey    core_types.Nullable[string]             `json:"photo_key" validate:"omitempty,max=2048"`
	Description core_types.Nullable[string]             `json:"description" validate:"omitempty,max=2000"`
	Status      core_types.Nullable[string]             `json:"status" validate:"omitempty,oneof=active hidden archived"`
}

func (p *VenuePatch) Validate() error {
	if p.Name.Set && p.Name.Value == nil {
		return fmt.Errorf("name cannot be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if p.CityID.Set {
		if p.CityID.Value == nil {
			return fmt.Errorf("city_id cannot be patched to NULL: %w", core_errors.ErrInvalidArgument)
		}
		if *p.CityID.Value <= 0 {
			return fmt.Errorf("city_id must be positive: %w", core_errors.ErrInvalidArgument)
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

func (v *Venue) ApplyPatch(patch VenuePatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate venue patch: %w", err)
	}

	tmp := *v

	if patch.CityID.Set && patch.CityID.Value != nil {
		tmp.CityID = *patch.CityID.Value
	}

	if patch.Name.Set && patch.Name.Value != nil {
		tmp.Name = *patch.Name.Value
	}

	if patch.Address.Set {
		tmp.Address = patch.Address.Value
	}

	if patch.Capacity.Set {
		tmp.Capacity = patch.Capacity.Value
	}

	if patch.SocialLinks.Set {
		if patch.SocialLinks.Value == nil {
			tmp.SocialLinks = nil
		} else {
			tmp.SocialLinks = *patch.SocialLinks.Value
		}
	}

	if patch.PhotoKey.Set {
		tmp.PhotoURL = patch.PhotoKey.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.Status.Set && patch.Status.Value != nil {
		tmp.Status = ContentStatus(*patch.Status.Value)
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched venue: %w", err)
	}

	*v = tmp

	return nil
}

type VenueDependencies struct {
	ConcertsCount int
	ReviewsCount  int
}

func (d VenueDependencies) HasAny() bool {
	return d.ConcertsCount > 0 || d.ReviewsCount > 0
}
