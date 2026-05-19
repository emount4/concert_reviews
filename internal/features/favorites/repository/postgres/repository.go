package favorites_postgres_repository

import (
	core_postgres_pool "github.com/emount4/concert_reviews/internal/core/repository/postgres/pool"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
)

type FavoritesRepository struct {
	pool      core_postgres_pool.Pool
	txManager *core_postgres_tx.Manager
}

func NewFavoritesRepository(
	pool core_postgres_pool.Pool,
	txManager *core_postgres_tx.Manager,
) *FavoritesRepository {
	return &FavoritesRepository{
		pool:      pool,
		txManager: txManager,
	}
}
