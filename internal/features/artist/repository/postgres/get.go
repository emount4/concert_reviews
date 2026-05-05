package artist_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

// GetByID — получение по UUID с учетом Soft Delete
func (r *ArtistRepository) GetArtistByID(ctx context.Context, id int) (domain.Artist, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT artist_id, name, description, photo_url, social_links, status, created_at, deleted_at
		FROM artists
		WHERE artist_id = $1 AND deleted_at IS NULL
	`

	var rec ArtistRecord
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&rec.ArtistID, &rec.Name, &rec.Description, &rec.PhotoURL,
		&rec.SocialLinks, &rec.Status, &rec.CreatedAt, &rec.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Artist{}, core_errors.ErrNotFound
		}
		return domain.Artist{}, fmt.Errorf("select artist by id: %w", err)
	}

	return rec.MapToDomain(), nil
}

// List — поиск и пагинация
func (r *ArtistRepository) GetArtists(ctx context.Context, search string, limit, offset *int) ([]domain.Artist, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	// Базовый запрос
	query := `
		SELECT artist_id, name, description, photo_url, social_links, status, created_at, deleted_at
		FROM artists
		WHERE deleted_at IS NULL
	`
	args := []any{}

	// Если есть поиск по имени
	if search != "" {
		query += ` AND name ILIKE $1`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY name ASC`

	// Пагинация
	if limit != nil && offset != nil {
		argNum := len(args) + 1
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argNum, argNum+1)
		args = append(args, *limit, *offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query artists list: %w", err)
	}
	defer rows.Close()

	var artists []domain.Artist
	for rows.Next() {
		var rec ArtistRecord
		err := rows.Scan(
			&rec.ArtistID, &rec.Name, &rec.Description, &rec.PhotoURL,
			&rec.SocialLinks, &rec.Status, &rec.CreatedAt, &rec.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan artist: %w", err)
		}
		artists = append(artists, rec.MapToDomain())
	}

	return artists, nil
}
