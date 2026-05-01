package city_service

import (
	"context"
	"fmt"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

func (s *CityService) Create(ctx context.Context, city core_models.City) (core_models.City, error) {
	if s.cityRepository == nil {
		return core_models.City{},
			fmt.Errorf(
				"city repository: %w",
				core_errors.ErrRepositoryNotConfigured,
			)
	}

	if err := city.Validate(); err != nil {
		return core_models.City{},
			fmt.Errorf("validate city request: %w", err)
	}

	if city.Slug == "" {
		city.Slug = s.generateSlug(city.Name)
	}

	createdCity, err := s.cityRepository.CreateCity(ctx, city)

	if err != nil {
		return core_models.City{},
			fmt.Errorf("city repository: %w", err)
	}

	return createdCity, nil
}
