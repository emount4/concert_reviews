package city_service

import (
	"context"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
)

type CityService struct {
	cityRepository CityRepository
}

type CityRepository interface {
	CreateCity(ctx context.Context, city core_models.City) (core_models.City, error)
	GetCities(ctx context.Context, limit, offset *int) ([]core_models.City, error)
	GetCityByID(ctx context.Context, id int) (core_models.City, error)

	DeleteCity(ctx context.Context, id int) error
	HasVenues(ctx context.Context, id int) (bool, error)

	UpdateCity(ctx context.Context, city core_models.City) (core_models.City, error)
}

func NewCityService(
	cityRepository CityRepository,
) *CityService {
	return &CityService{
		cityRepository: cityRepository,
	}
}
