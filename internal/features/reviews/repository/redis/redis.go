package review_redis_repository

import (
	"context"
	"fmt"

	core_redis "github.com/emount4/concert_reviews/internal/core/repository/redis"
)

type StatsReviewRepository struct {
	client *core_redis.Client
}

func NewReviewRedisRepository(client *core_redis.Client) *StatsReviewRepository {
	return &StatsReviewRepository{
		client: client,
	}
}

func (r *StatsReviewRepository) InvalidateGlobalStats(ctx context.Context) error {
	err := r.client.Del(ctx, "stats:global").Err()
	if err != nil {
		return fmt.Errorf("redis delete global stats key: %w", err)
	}
	return nil
}
