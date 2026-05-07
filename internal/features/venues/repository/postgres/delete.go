package venue_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
)

// VenueDependencies — структура для проверки связей
type VenueDependencies struct {
	ConcertsCount int `db:"concerts_count"`
	ReviewsCount  int `db:"reviews_count"`
}

func (d VenueDependencies) HasAny() bool {
	return d.ConcertsCount > 0 || d.ReviewsCount > 0
}

// DeleteVenueHard — физическое удаление
func (r *VenueRepository) DeleteVenueHard(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `DELETE FROM venues WHERE venue_id = $1 AND deleted_at IS NULL`
	tag, err := exec.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec hard delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}
	return nil
}

// DeleteVenueSoft — логическое удаление
func (r *VenueRepository) DeleteVenueSoft(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `UPDATE venues SET deleted_at = NOW() WHERE venue_id = $1 AND deleted_at IS NULL`
	tag, err := exec.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec soft delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}
	return nil
}

// RestoreVenue — восстановление записи
func (r *VenueRepository) RestoreVenue(ctx context.Context, id int) (domain.Venue, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `
		UPDATE venues SET deleted_at = NULL 
		WHERE venue_id = $1 AND deleted_at IS NOT NULL
		RETURNING venue_id, city_id, name, address, capacity, photo_url, description, social_links, status, created_at, deleted_at
	`

	var rec VenueRecord
	err := exec.QueryRow(ctx, query, id).Scan(
		&rec.VenueID, &rec.CityID, &rec.Name, &rec.Address, &rec.Capacity,
		&rec.PhotoURL, &rec.Description, &rec.SocialLinks, &rec.Status,
		&rec.CreatedAt, &rec.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Venue{}, core_errors.ErrNotFound
		}
		return domain.Venue{}, fmt.Errorf("restore venue query: %w", err)
	}

	return rec.MapToDomain(), nil
}

// GetVenueDependencies — подсчёт связанных концертов и отзывов
func (r *VenueRepository) GetVenueDependencies(ctx context.Context, id int) (domain.VenueDependencies, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	// Отзывы привязаны к концертам, концерты к площадкам
	query := `
		SELECT 
			COUNT(DISTINCT c.concert_id) AS concerts_count,
			COUNT(DISTINCT rv.review_id) AS reviews_count
		FROM venues v
		LEFT JOIN concerts c ON v.venue_id = c.venue_id AND c.deleted_at IS NULL
		LEFT JOIN reviews rv ON c.concert_id = rv.concert_id 
			AND rv.deleted_at IS NULL AND rv.status = 'approved'
		WHERE v.venue_id = $1 AND v.deleted_at IS NULL
		GROUP BY v.venue_id
	`

	var deps domain.VenueDependencies
	err := exec.QueryRow(ctx, query, id).Scan(&deps.ConcertsCount, &deps.ReviewsCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.VenueDependencies{}, nil // Нет связей или площадка не найдена
		}
		return domain.VenueDependencies{}, fmt.Errorf("count venue dependencies: %w", err)
	}

	return deps, nil
}
