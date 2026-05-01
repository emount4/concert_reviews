package postgres_city_repository

import (
	"context"
	"errors"
	"fmt"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *CityRepository) CreateCity(ctx context.Context, city core_models.City) (core_models.City, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `
		INSERT INTO cities (name, slug, timezone)
		VALUES ($1, $2, $3)
		RETURNING city_id, created_at, deleted_at
	`

	row := exec.QueryRow(ctx, query, city.Name, city.Slug, city.Timezone)

	if err := row.Scan(
		&city.CityID,
		&city.CreatedAt,
		&city.DeletedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return core_models.City{}, fmt.Errorf("%w: city already exists", core_errors.ErrConflict)
		}
		return core_models.City{}, fmt.Errorf("insert city: %w", err)
	}

	return city, nil
}
