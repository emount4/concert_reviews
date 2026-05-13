package review_postgres_repository

import (
	core_postgres_pool "github.com/emount4/concert_reviews/internal/core/repository/postgres/pool"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
)

type ReviewRepository struct {
	txManager *core_postgres_tx.Manager
	pool      core_postgres_pool.Pool
}

func NewReviewRepository(
	pool core_postgres_pool.Pool,
	txManager *core_postgres_tx.Manager,
) *ReviewRepository {
	return &ReviewRepository{
		pool:      pool,
		txManager: txManager,
	}
}
