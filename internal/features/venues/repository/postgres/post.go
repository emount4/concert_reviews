package venue_postgres_repository

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
)

func (r *VenueRepository) CreateVenue(ctx context.Context, venue domain.Venue) (domain.Venue, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `
		INSERT INTO venues (
			city_id, name, address, capacity, photo_url, description, social_links, status
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING venue_id, created_at, deleted_at
	`

	// Выполняем вставку и получаем ID и время создания
	err := exec.QueryRow(ctx, query,
		venue.CityID,
		venue.Name,
		venue.Address,
		venue.Capacity,
		venue.PhotoURL,
		venue.Description,
		venue.SocialLinks,
		venue.Status,
	).Scan(&venue.VenueID, &venue.CreatedAt, &venue.DeletedAt)

	if err != nil {
		return domain.Venue{}, fmt.Errorf("insert venue: %w", err)
	}

	return venue, nil
}

func (r *VenueRepository) CityExists(ctx context.Context, id int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `SELECT EXISTS(SELECT 1 FROM cities WHERE city_id = $1 AND deleted_at IS NULL)`

	var exists bool
	err := exec.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check city existence: %w", err)
	}

	return exists, nil
}
