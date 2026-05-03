package city_service

import (
	"context"
	"fmt"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

func (s *CityService) Create(
	ctx context.Context,
	city core_models.City,
) (core_models.City, error) {
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

func (s *CityService) GetCities(
	ctx context.Context,
	limit, offset *int,
) ([]core_models.City, error) {
	if s.cityRepository == nil {
		return []core_models.City{},
			fmt.Errorf(
				"city repository: %w",
				core_errors.ErrRepositoryNotConfigured,
			)
	}

	if limit != nil && *limit < 0 {
		return nil,
			fmt.Errorf(
				"limit must be non-negative: %w",
				core_errors.ErrInvalidArgument,
			)
	}

	if offset != nil && *offset < 0 {
		return nil,
			fmt.Errorf(
				"offset must be non-negative: %w",
				core_errors.ErrInvalidArgument,
			)
	}

	cities, err := s.cityRepository.GetCities(ctx, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("get cities from repository: %w", err)
	}

	return cities, nil
}

func (s *CityService) GetCityByID(
	ctx context.Context,
	id int,
) (core_models.City, error) {
	if s.cityRepository == nil {
		return core_models.City{},
			fmt.Errorf(
				"city repository: %w",
				core_errors.ErrRepositoryNotConfigured,
			)
	}

	city, err := s.cityRepository.GetCityByID(ctx, id)
	if err != nil {
		return core_models.City{}, err
	}

	return city, nil
}

func (s *CityService) Delete(ctx context.Context, id int) error {
	if s.cityRepository == nil {
		return fmt.Errorf("city repository: %w", core_errors.ErrRepositoryNotConfigured)
	}

	// 1. Проверяем, есть ли площадки в этом городе
	hasVenues, err := s.cityRepository.HasVenues(ctx, id)
	if err != nil {
		return fmt.Errorf("check city dependencies: %w", err)
	}
	if hasVenues {
		return fmt.Errorf("cannot delete city with linked venues: %w", core_errors.ErrConflict)
	}

	// 2. Выполняем удаление
	if err := s.cityRepository.DeleteCity(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *CityService) Update(ctx context.Context, id int, name, slug, tz *string) (core_models.City, error) {
	city, err := s.cityRepository.GetCityByID(ctx, id)
	if err != nil {
		return core_models.City{}, err
	}

	if name != nil {
		city.Name = *name
	}
	if slug != nil {
		city.Slug = *slug
	}
	if tz != nil {
		city.Timezone = *tz
	}

	if err := city.Validate(); err != nil {
		return core_models.City{}, fmt.Errorf("validate updated city: %w", err)
	}

	return s.cityRepository.UpdateCity(ctx, city)
}
