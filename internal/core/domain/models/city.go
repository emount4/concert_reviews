package core_models

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

type City struct {
	CityID    int        `db:"city_id"`
	Name      string     `db:"name"`
	Slug      string     `db:"slug"`
	Timezone  string     `db:"timezone"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func NewCity(name, slug, timezone string) City {
	return City{
		Name:     name,
		Slug:     slug,
		Timezone: timezone,
	}
}

func (c *City) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("%w: name is required", core_errors.ErrInvalidArgument)
	}

	if strings.TrimSpace(c.Timezone) == "" {
		return fmt.Errorf("%w: timezone is required", core_errors.ErrInvalidArgument)
	}

	// Пробуем загрузить как IANA timezone
	_, err := time.LoadLocation(c.Timezone)
	if err != nil {
		// Если не получилось — проверяем, может это UTC offset
		if !isValidUTCOffset(c.Timezone) {
			return fmt.Errorf("%w: invalid timezone '%s': must be IANA name (e.g. 'Europe/Moscow') or UTC offset (e.g. 'UTC+3', 'UTC-5:30')",
				core_errors.ErrInvalidArgument, c.Timezone)
		}
	}

	return nil
}

func isValidUTCOffset(s string) bool {
	// Форматы: UTC+3, UTC-5, UTC+5:30, UTC-03:00
	matched, _ := regexp.MatchString(`^UTC[+-]\d{1,2}(:\d{2})?$`, s)
	return matched
}
