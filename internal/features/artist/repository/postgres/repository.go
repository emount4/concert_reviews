package artist_postgres_repository

import core_postgres_pool "github.com/emount4/concert_reviews/internal/core/repository/postgres/pool"

type ArtistRepository struct {
	pool core_postgres_pool.Pool
}

func NewArtistRepository(pool core_postgres_pool.Pool) *ArtistRepository {
	return &ArtistRepository{pool: pool}
}
