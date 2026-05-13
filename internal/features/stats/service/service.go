package stats_service

import (
	"context"
	"fmt"
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	"go.uber.org/zap"
)

type StatsService struct {
	statsRepository StatsRepository
	statsRedis      StatsRedis
}

type StatsRepository interface {
	GetGlobalStats(ctx context.Context) (domain.GlobalStats, error)
}

type StatsRedis interface {
	SetGlobalStats(ctx context.Context, stats domain.GlobalStats, ttl time.Duration) error
	GetGlobalStats(ctx context.Context) (domain.GlobalStats, error)
	InvalidateGlobalStats(ctx context.Context) error
}

func NewStatsService(repo StatsRepository, redis StatsRedis) *StatsService {
	return &StatsService{
		statsRepository: repo,
		statsRedis:      redis,
	}
}

func (s *StatsService) GetGlobalStats(ctx context.Context) (domain.GlobalStats, error) {
	logger := core_logger.FromContext(ctx)

	cachedStats, err := s.statsRedis.GetGlobalStats(ctx)
	if err == nil {
		return cachedStats, nil
	}

	logger.Debug("redis read", zap.Error(err))

	stats, err := s.statsRepository.GetGlobalStats(ctx)
	if err != nil {
		return domain.GlobalStats{}, fmt.Errorf("repository get global stats: %w", err)
	}

	err = s.statsRedis.SetGlobalStats(ctx, stats, 10*time.Minute)
	if err != nil {
		logger.Debug("redis write", zap.Error(err))
	}

	return stats, nil
}
