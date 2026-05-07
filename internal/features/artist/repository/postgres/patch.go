package artist_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/jackc/pgx/v5"
)

func (r *ArtistRepository) PatchArtist(ctx context.Context, id int, artist domain.Artist) (domain.Artist, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `
		UPDATE artists
		SET 
			name = $1,
			description = $2,
			photo_url = $3,
			social_links = $4,
			status = $5
		WHERE artist_id = $6 AND deleted_at IS NULL
		RETURNING artist_id, name, description, photo_url, social_links, status, created_at, deleted_at
	`

	err := exec.QueryRow(ctx, query,
		artist.Name,
		artist.Description,
		artist.PhotoURL,
		artist.SocialLinks, // pgx автоматически закодирует map[string]string в JSONB
		artist.Status,
		id,
	).Scan(
		&artist.ArtistID,
		&artist.Name,
		&artist.Description,
		&artist.PhotoURL,
		&artist.SocialLinks, // pgx автоматически декодирует JSONB в map[string]string
		&artist.Status,
		&artist.CreatedAt,
		&artist.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Artist{}, core_errors.ErrNotFound
		}
		return domain.Artist{}, fmt.Errorf("update artist: %w", err)
	}

	return artist, nil
}
