package concert_postgres_repository

import (
	core_postgres_pool "github.com/emount4/concert_reviews/internal/core/repository/postgres/pool"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
)

type ConcertRepository struct {
	pool      core_postgres_pool.Pool
	txManager *core_postgres_tx.Manager
}

func NewConcertRepository(pool core_postgres_pool.Pool, txManager *core_postgres_tx.Manager) *ConcertRepository {
	return &ConcertRepository{
		pool:      pool,
		txManager: txManager,
	}
}
