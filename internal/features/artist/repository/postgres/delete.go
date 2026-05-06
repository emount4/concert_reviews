package artist_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	artist_service "github.com/emount4/concert_reviews/internal/features/artist/service"
	"github.com/jackc/pgx/v5"
)

func (r *ArtistRepository) DeleteArtistHard(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `DELETE FROM artists WHERE artist_id = $1 AND deleted_at IS NULL`

	tag, err := exec.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec hard delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}

	return nil
}

func (r *ArtistRepository) DeleteArtistSoft(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `
		UPDATE artists 
		SET deleted_at = NOW() 
		WHERE artist_id = $1 AND deleted_at IS NULL
	`

	tag, err := exec.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec soft delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}

	return nil
}

func (r *ArtistRepository) GetArtistDependencies(
	ctx context.Context,
	id int,
) (artist_service.ArtistDependencies, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	// Учитываем, что отзывы привязаны к концертам, а не напрямую к артистам
	query := `
		SELECT 
			COUNT(DISTINCT ca.concert_id) AS concerts_count,
			COUNT(DISTINCT r.review_id) AS reviews_count
		FROM artists a
		LEFT JOIN concert_artists ca ON a.artist_id = ca.artist_id
		LEFT JOIN concerts c ON ca.concert_id = c.concert_id 
			AND c.deleted_at IS NULL
		LEFT JOIN reviews r ON c.concert_id = r.concert_id 
			AND r.deleted_at IS NULL AND r.status = 'approved'
		WHERE a.artist_id = $1 AND a.deleted_at IS NULL
		GROUP BY a.artist_id
	`

	var deps artist_service.ArtistDependencies
	err := exec.QueryRow(ctx, query, id).Scan(&deps.ConcertsCount, &deps.ReviewsCount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return artist_service.ArtistDependencies{}, nil
		}
		return artist_service.ArtistDependencies{}, fmt.Errorf("count dependencies: %w", err)
	}

	return deps, nil
}

func (r *ArtistRepository) RestoreArtist(ctx context.Context, id int) (domain.Artist, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `
		UPDATE artists 
		SET deleted_at = NULL 
		WHERE artist_id = $1 AND deleted_at IS NOT NULL
		RETURNING artist_id, name, description, photo_url, social_links, status, created_at, deleted_at
	`

	var record ArtistRecord
	err := exec.QueryRow(ctx, query, id).Scan(
		&record.ArtistID, &record.Name, &record.Description, &record.PhotoURL,
		&record.SocialLinks, &record.Status, &record.CreatedAt, &record.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Artist{}, core_errors.ErrNotFound
		}
		return domain.Artist{}, fmt.Errorf("restore artist query: %w", err)
	}

	return record.MapToDomain(), nil
}
