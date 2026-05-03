package postgres_city_repository

import (
	"context"
	"errors"
	"fmt"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/jackc/pgx/v5"
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

func (r *CityRepository) GetCities(ctx context.Context, limit, offset *int) ([]core_models.City, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT city_id, name, slug, timezone, created_at 
	FROM cities
	ORDER BY city_id ASC
	LIMIT  $1
	OFFSET $2 
	`

	rows, err := r.pool.Query(
		ctx,
		query,
		limit,
		offset,
	)

	if err != nil {
		return nil, fmt.Errorf("select cities: %w", err)
	}

	defer rows.Close()

	var cityModels []core_models.City
	for rows.Next() {
		var cityModel core_models.City

		err := rows.Scan(
			&cityModel.CityID,
			&cityModel.Name,
			&cityModel.Slug,
			&cityModel.Timezone,
			&cityModel.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan city: %w", err)
		}

		cityModels = append(cityModels, cityModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return cityModels, nil
}

func (r *CityRepository) GetCityByID(ctx context.Context, id int) (core_models.City, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT city_id, name, slug, timezone, created_at 
		FROM cities 
		WHERE city_id = $1 AND deleted_at IS NULL
	`

	var city core_models.City
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&city.CityID,
		&city.Name,
		&city.Slug,
		&city.Timezone,
		&city.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_models.City{}, core_errors.ErrNotFound
		}
		return core_models.City{}, fmt.Errorf("select city by id: %w", err)
	}

	return city, nil
}

func (r *CityRepository) DeleteCity(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM cities WHERE city_id = $1`

	res, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("execute delete query: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("city with id ='%d': %w", id, core_errors.ErrNotFound)
	}

	return nil
}

// Дополнительный метод для проверки зависимостей
func (r *CityRepository) HasVenues(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM venues WHERE city_id = $1 AND deleted_at IS NULL)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	return exists, err
}

func (r *CityRepository) UpdateCity(ctx context.Context, city core_models.City) (core_models.City, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE cities SET name = $1, slug = $2, timezone = $3 
		WHERE city_id = $4 RETURNING created_at`
	err := r.pool.QueryRow(ctx, query, city.Name, city.Slug, city.Timezone, city.CityID).Scan(&city.CreatedAt)
	return city, err
}
