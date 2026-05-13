package stats_redis_repository

import core_redis "github.com/emount4/concert_reviews/internal/core/repository/redis"

type StatsRedisRepository struct {
	client *core_redis.Client
}

func NewStatsRedisRepository(client *core_redis.Client) *StatsRedisRepository {
	return &StatsRedisRepository{
		client: client,
	}
}
