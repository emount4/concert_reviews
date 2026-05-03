package city_transport_http

import (
	"time"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
)

type CityResponse struct {
	CityID    int       `json:"city_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
}

type NewCityRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Slug     string `json:"slug,omitempty" validate:"max=100"`
	Timezone string `json:"timezone" validate:"required"`
}

type UpdateCityRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=2"`
	Slug     *string `json:"slug" validate:"omitempty,min=2"`
	Timezone *string `json:"timezone" validate:"omitempty"`
}

type NewCityResponse CityResponse
type GetCityResponse CityResponse

type GetCitiesResponse []CityResponse

func CityDTOFromDomain(city core_models.City) CityResponse {
	return CityResponse{
		CityID:    city.CityID,
		Name:      city.Name,
		Slug:      city.Slug,
		Timezone:  city.Timezone,
		CreatedAt: city.CreatedAt,
	}

}

func CitiesDTOFromDomain(cities []core_models.City) []CityResponse {
	citiesResp := make([]CityResponse, len(cities))
	for i, domain := range cities {
		citiesResp[i] = CityDTOFromDomain(domain)
	}
	return citiesResp
}

func MapCreateDTOToDomain(dto NewCityRequest) core_models.City {
	return core_models.City{
		Name:     dto.Name,
		Slug:     dto.Slug,
		Timezone: dto.Timezone,
	}
}
