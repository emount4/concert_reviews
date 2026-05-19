package moderation_postgres_repository

import core_postgres_pool "github.com/emount4/concert_reviews/internal/core/repository/postgres/pool"

type ModerationRepository struct {
	pool core_postgres_pool.Pool
}

func NewModerationRepository(pool core_postgres_pool.Pool) *ModerationRepository {
	return &ModerationRepository{
		pool: pool,
	}
}
