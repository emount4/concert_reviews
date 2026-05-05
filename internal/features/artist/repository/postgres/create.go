package artist_postgres_repository

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
)

func (r *ArtistRepository) CreateArtist(ctx context.Context, a domain.Artist) (domain.Artist, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `
		INSERT INTO artists (name, description, photo_url, social_links, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING artist_id, created_at
	`

	err := exec.QueryRow(ctx, query,
		a.Name, a.Description, a.PhotoURL, a.SocialLinks, a.Status,
	).Scan(&a.ArtistID, &a.CreatedAt)

	if err != nil {
		return domain.Artist{}, fmt.Errorf("insert artist: %w", err)
	}

	return a, nil
}
