package stats_redis_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
)

func (r *StatsRedisRepository) SetGlobalStats(ctx context.Context, stats domain.GlobalStats, ttl time.Duration) error {
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, "stats:global", data, ttl).Err()
}

func (r *StatsRedisRepository) GetGlobalStats(ctx context.Context) (domain.GlobalStats, error) {
	data, err := r.client.Get(ctx, "stats:global").Bytes()
	if err != nil {
		return domain.GlobalStats{}, fmt.Errorf("cannot get data from redis: %w", err)
	}

	var stats domain.GlobalStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return domain.GlobalStats{}, err
	}

	return stats, nil
}

func (r *StatsRedisRepository) InvalidateGlobalStats(ctx context.Context) error {
	return r.client.Del(ctx, "stats:global").Err()
}
