package domain

import (
	"fmt"
	"strings"
	"time"
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
