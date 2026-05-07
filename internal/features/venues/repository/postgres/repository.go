package venue_postgres_repository

import core_postgres_pool "github.com/emount4/concert_reviews/internal/core/repository/postgres/pool"

type VenueRepository struct {
	pool core_postgres_pool.Pool
}

func NewVenueRepository(pool core_postgres_pool.Pool) *VenueRepository {
	return &VenueRepository{pool: pool}
}
