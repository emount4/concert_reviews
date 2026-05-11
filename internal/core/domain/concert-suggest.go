package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ConcertSuggestion struct {
	SuggestionID  uuid.UUID
	UserID        uuid.UUID
	RawArtistName string    // "сырое" имя артиста
	RawVenueName  string    // "сырое" название площадки
	ConcertDate   time.Time // обязательно по требованию
	Info          *string   // доп. комментарий
	CreatedAt     time.Time
}

// Validate проверяет бизнес-правила при создании предложения
func (s *ConcertSuggestion) Validate() error {
	// Дата обязательна
	if s.ConcertDate.IsZero() {
		return fmt.Errorf("concert_date is required")
	}

	// Хотя бы одно из полей должно быть заполнено
	if strings.TrimSpace(s.RawArtistName) == "" && strings.TrimSpace(s.RawVenueName) == "" {
		return fmt.Errorf("raw_artist_name or raw_venue_name must be provided")
	}

	// Ограничения длины как в БД
	if len(s.RawArtistName) > 255 {
		return fmt.Errorf("raw_artist_name too long (max 255)")
	}
	if len(s.RawVenueName) > 255 {
		return fmt.Errorf("raw_venue_name too long (max 255)")
	}

	// Info опционален, но с лимитом
	if s.Info != nil && len(*s.Info) > 2000 {
		return fmt.Errorf("info too long (max 2000)")
	}

	return nil
}
