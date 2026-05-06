package domain

import (
	"fmt"
	"strings"
	"time"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_http_types "github.com/emount4/concert_reviews/internal/core/transport/http/types"
	core_types "github.com/emount4/concert_reviews/internal/core/types"
)

// ContentStatus — тип для статуса контента (Enum)
type ContentStatus string

const (
	StatusActive   ContentStatus = "active"
	StatusHidden   ContentStatus = "hidden"
	StatusArchived ContentStatus = "archived"
)

// Artist представляет доменную модель артиста
type Artist struct {
	ArtistID    int
	Name        string
	Description *string
	PhotoURL    *string
	SocialLinks map[string]string
	Status      ContentStatus
	CreatedAt   time.Time
	DeletedAt   *time.Time
}

// Validate проверяет доменную модель на соответствие бизнес-правилам и ограничениям БД
func (a *Artist) Validate() error {
	// Проверка имени (NOT NULL и не пустая строка)
	if strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("artist name cannot be empty")
	}
	if len(a.Name) > 255 {
		return fmt.Errorf("artist name too long (max 255)")
	}

	// Проверка описания (max 2000)
	if a.Description != nil && len(*a.Description) > 2000 {
		return fmt.Errorf("description too long (max 2000)")
	}

	// Проверка фото (не пустая строка, если передана)
	if a.PhotoURL != nil {
		if strings.TrimSpace(*a.PhotoURL) == "" {
			return fmt.Errorf("photo_url cannot be empty string if provided")
		}
		if len(*a.PhotoURL) > 2048 {
			return fmt.Errorf("photo_url too long (max 2048)")
		}
	}

	return nil
}

// IsDeleted — вспомогательный метод для проверки статуса удаления
func (a *Artist) IsDeleted() bool {
	return a.DeletedAt != nil
}

type ArtistPatch struct {
	Name        core_types.Nullable[string]             `json:"name" validate:"omitempty,min=2,max=255"`
	Description core_types.Nullable[string]             `json:"description" validate:"omitempty,max=2000"`
	PhotoKey    core_types.Nullable[string]             `json:"photo_key" validate:"omitempty,max=2048"`
	SocialLinks core_http_types.NullableMapStringString `json:"social_links"`
}

func (p *ArtistPatch) Validate() error {
	if p.Name.Set && p.Name.Value == nil {
		return fmt.Errorf("name cant be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (a *Artist) ApplyPatch(patch ArtistPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate artist patch: %w", err)
	}

	tmp := *a

	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.PhotoKey.Set {
		tmp.PhotoURL = patch.PhotoKey.Value
	}

	if patch.SocialLinks.Set {
		if patch.SocialLinks.Value == nil {
			tmp.SocialLinks = nil
		} else {
			tmp.SocialLinks = *patch.SocialLinks.Value
		}
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched artist: %w", err)
	}

	*a = tmp

	return nil
}
