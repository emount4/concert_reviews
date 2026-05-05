package postgres_city_repository

import core_postgres_pool "github.com/emount4/concert_reviews/internal/core/repository/postgres/pool"

type CityRepository struct {
	pool core_postgres_pool.Pool
}

func NewCityRepository(
	pool core_postgres_pool.Pool,
) *CityRepository {
	return &CityRepository{
		pool: pool,
	}
}
