package city_transport_http

import "time"

type NewCityRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Slug     string `json:"slug,omitempty" validate:"max=100"`
	Timezone string `json:"timezone" validate:"required"`
}

type NewCityResponse struct {
	CityID    int       `json:"city_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
}
