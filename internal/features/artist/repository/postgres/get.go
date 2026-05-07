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

// GetByID — получение по UUID с учетом Soft Delete
func (r *ArtistRepository) GetArtistByID(ctx context.Context, id int) (domain.Artist, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT a.artist_id, a.name, a.description, a.photo_url, a.social_links,
		       COALESCE(s.reviews_count, 0) AS reviews_count,
		       COALESCE(s.sum_rating_total, 0) AS sum_rating_total,
		       COALESCE(s.concerts_count, 0) AS concerts_count,
		       COALESCE(s.favorites_count, 0) AS favorites_count,
		       COALESCE(s.updated_at, a.created_at) AS stats_updated_at,
		       a.status, a.created_at, a.deleted_at
		FROM artists a
		LEFT JOIN artist_stats s ON a.artist_id = s.artist_id
		WHERE a.artist_id = $1 AND a.deleted_at IS NULL
	`

	var rec ArtistRecord
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&rec.ArtistID, &rec.Name, &rec.Description, &rec.PhotoURL,
		&rec.SocialLinks, &rec.StatsReviewsCount, &rec.StatsSumRatingTotal,
		&rec.StatsConcertsCount, &rec.StatsFavoritesCount, &rec.StatsUpdatedAt,
		&rec.Status, &rec.CreatedAt, &rec.DeletedAt,
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
		SELECT a.artist_id, a.name, a.description, a.photo_url, a.social_links,
		       COALESCE(s.reviews_count, 0) AS reviews_count,
		       COALESCE(s.sum_rating_total, 0) AS sum_rating_total,
		       COALESCE(s.concerts_count, 0) AS concerts_count,
		       COALESCE(s.favorites_count, 0) AS favorites_count,
		       COALESCE(s.updated_at, a.created_at) AS stats_updated_at,
		       a.status, a.created_at, a.deleted_at
		FROM artists a
		LEFT JOIN artist_stats s ON a.artist_id = s.artist_id
		WHERE a.deleted_at IS NULL
	`
	args := []any{}
	argIdx := 1

	// Если есть поиск по имени
	if search != "" {
		query += fmt.Sprintf(" AND a.name ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += ` ORDER BY a.name ASC`

	// Пагинация
	if limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, *limit)
		argIdx++
	}
	if offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, *offset)
		argIdx++
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
			&rec.SocialLinks, &rec.StatsReviewsCount, &rec.StatsSumRatingTotal,
			&rec.StatsConcertsCount, &rec.StatsFavoritesCount, &rec.StatsUpdatedAt,
			&rec.Status, &rec.CreatedAt, &rec.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan artist: %w", err)
		}
		artists = append(artists, rec.MapToDomain())
	}

	return artists, nil
}

// GetArtistsAdmin — список с фильтрами для админки
func (r *ArtistRepository) GetArtistsAdmin(
	ctx context.Context,
	search string,
	limit, offset *int,
	includeDeleted bool,
	status string,
) ([]domain.Artist, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	// Базовый запрос
	query := `
		SELECT a.artist_id, a.name, a.description, a.photo_url, a.social_links,
		       COALESCE(s.reviews_count, 0) AS reviews_count,
		       COALESCE(s.sum_rating_total, 0) AS sum_rating_total,
		       COALESCE(s.updated_at, a.created_at) AS stats_updated_at,
		       a.status, a.created_at, a.deleted_at
		FROM artists a
		LEFT JOIN artist_stats s ON a.artist_id = s.artist_id
		WHERE 1=1
	`

	args := make([]interface{}, 0)
	argIdx := 1

	// Фильтр по поиску (по имени, регистронезависимый)
	if search != "" {
		query += fmt.Sprintf(" AND a.name ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	// Фильтр по статусу
	if status != "" {
		query += fmt.Sprintf(" AND a.status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	// Фильтр по deleted_at
	if !includeDeleted {
		query += " AND a.deleted_at IS NULL"
	}

	// Сортировка по умолчанию: новые сверху
	query += " ORDER BY a.created_at DESC"

	// Пагинация
	if limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, *limit)
		argIdx++
	}
	if offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, *offset)
		argIdx++
	}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query artists admin: %w", err)
	}
	defer rows.Close()

	var artists []domain.Artist
	for rows.Next() {
		var record ArtistRecord
		err := rows.Scan(
			&record.ArtistID, &record.Name, &record.Description, &record.PhotoURL,
			&record.SocialLinks, &record.StatsReviewsCount, &record.StatsSumRatingTotal, &record.StatsUpdatedAt,
			&record.StatsConcertsCount, &record.StatsFavoritesCount, &record.StatsUpdatedAt,
			&record.Status, &record.CreatedAt, &record.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan artist row: %w", err)
		}
		artists = append(artists, record.MapToDomain())
	}

	return artists, nil
}
