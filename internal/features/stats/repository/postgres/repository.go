package stats_postgres_repository

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_postgres_pool "github.com/emount4/concert_reviews/internal/core/repository/postgres/pool"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
)

type StatsRepository struct {
	pool      core_postgres_pool.Pool
	txManager *core_postgres_tx.Manager
}

func NewStatsRepository(pool core_postgres_pool.Pool, txManager *core_postgres_tx.Manager) *StatsRepository {
	return &StatsRepository{
		pool:      pool,
		txManager: txManager,
	}
}

func (r *StatsRepository) GetGlobalStats(ctx context.Context) (domain.GlobalStats, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT 
			(SELECT COUNT(*) FROM users WHERE is_active = true) as users_count,
			(SELECT COUNT(*) FROM concerts WHERE deleted_at IS NULL) as concerts_count,
			(SELECT COUNT(*) FROM artists WHERE deleted_at IS NULL) as artists_count,
			(SELECT COUNT(*) FROM venues WHERE deleted_at IS NULL) as venues_count,
			(SELECT COUNT(*) FROM reviews WHERE status = 'approved' AND is_visible = true) as reviews_count
	`

	var s domain.GlobalStats
	err := r.pool.QueryRow(ctx, query).Scan(
		&s.UsersCount,
		&s.ConcertsCount,
		&s.ArtistsCount,
		&s.VenuesCount,
		&s.ReviewsCount,
	)

	if err != nil {
		return domain.GlobalStats{}, fmt.Errorf("fetch global stats: %w", err)
	}

	return s, nil
}
