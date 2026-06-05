package artist_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
func (r *ArtistRepository) GetArtists(ctx context.Context, search string, sort string, direction string, hasReviews *bool, limit, offset *int) ([]domain.Artist, int, error) {
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

	if hasReviews != nil {
		if *hasReviews {
			query += " AND COALESCE(s.reviews_count, 0) > 0"
		} else {
			query += " AND COALESCE(s.reviews_count, 0) = 0"
		}
	}

	sortExpr := "a.name"
	switch sort {
	case "rating":
		sortExpr = "COALESCE(CASE WHEN COALESCE(s.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(s.sum_rating_total, 0)::numeric / s.reviews_count END, 0)"
	case "reviews":
		sortExpr = "COALESCE(s.reviews_count, 0)"
	case "p1":
		sortExpr = "COALESCE(CASE WHEN COALESCE(s.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(s.sum_p1, 0)::numeric / s.reviews_count END, 0)"
	case "p2":
		sortExpr = "COALESCE(CASE WHEN COALESCE(s.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(s.sum_p2, 0)::numeric / s.reviews_count END, 0)"
	case "p3":
		sortExpr = "COALESCE(CASE WHEN COALESCE(s.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(s.sum_p3, 0)::numeric / s.reviews_count END, 0)"
	case "p4":
		sortExpr = "COALESCE(CASE WHEN COALESCE(s.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(s.sum_p4, 0)::numeric / s.reviews_count END, 0)"
	case "p5":
		sortExpr = "COALESCE(CASE WHEN COALESCE(s.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(s.sum_p5, 0)::numeric / s.reviews_count END, 0)"
	}

	dir := strings.ToUpper(direction)
	if dir != "ASC" {
		dir = "DESC"
	}

	countQuery := `
		SELECT COUNT(*)
		FROM artists a
		LEFT JOIN artist_stats s ON a.artist_id = s.artist_id
		WHERE a.deleted_at IS NULL`
	countArgs := []any{}
	argIdxCount := 1
	if search != "" {
		countQuery += fmt.Sprintf(" AND a.name ILIKE $%d", argIdxCount)
		countArgs = append(countArgs, "%"+search+"%")
		argIdxCount++
	}
	if hasReviews != nil {
		if *hasReviews {
			countQuery += " AND COALESCE(s.reviews_count, 0) > 0"
		} else {
			countQuery += " AND COALESCE(s.reviews_count, 0) = 0"
		}
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count artists: %w", err)
	}

	query += fmt.Sprintf(" ORDER BY %s %s, a.artist_id ASC", sortExpr, dir)

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
		return nil, 0, fmt.Errorf("query artists list: %w", err)
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
			return nil, 0, fmt.Errorf("scan artist: %w", err)
		}
		artists = append(artists, rec.MapToDomain())
	}

	if artists == nil {
		artists = []domain.Artist{}
	}
	return artists, total, nil
}

// GetArtistsAdmin — список с фильтрами для админки
func (r *ArtistRepository) GetArtistsAdmin(
	ctx context.Context,
	search string,
	sort string,
	direction string,
	hasReviews *bool,
	limit, offset *int,
	includeDeleted bool,
	status string,
) ([]domain.Artist, int, error) {
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
		       COALESCE(s.concerts_count, 0) AS concerts_count,
		       COALESCE(s.favorites_count, 0) AS favorites_count,
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
	if status == "deleted" {
		query += " AND a.deleted_at IS NOT NULL"
	} else if status != "" {
		query += fmt.Sprintf(" AND a.status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	// Фильтр по deleted_at
	if !includeDeleted && status != "deleted" {
		query += " AND a.deleted_at IS NULL"
	}

	if hasReviews != nil {
		if *hasReviews {
			query += " AND COALESCE(s.reviews_count, 0) > 0"
		} else {
			query += " AND COALESCE(s.reviews_count, 0) = 0"
		}
	}

	sortExpr := "a.name"
	switch sort {
	case "rating":
		sortExpr = "COALESCE(CASE WHEN COALESCE(s.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(s.sum_rating_total, 0)::numeric / s.reviews_count END, 0)"
	case "reviews":
		sortExpr = "COALESCE(s.reviews_count, 0)"
	case "p1":
		sortExpr = "COALESCE(CASE WHEN COALESCE(s.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(s.sum_p1, 0)::numeric / s.reviews_count END, 0)"
	case "p2":
		sortExpr = "COALESCE(CASE WHEN COALESCE(s.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(s.sum_p2, 0)::numeric / s.reviews_count END, 0)"
	case "p3":
		sortExpr = "COALESCE(CASE WHEN COALESCE(s.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(s.sum_p3, 0)::numeric / s.reviews_count END, 0)"
	case "p4":
		sortExpr = "COALESCE(CASE WHEN COALESCE(s.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(s.sum_p4, 0)::numeric / s.reviews_count END, 0)"
	case "p5":
		sortExpr = "COALESCE(CASE WHEN COALESCE(s.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(s.sum_p5, 0)::numeric / s.reviews_count END, 0)"
	}

	dir := strings.ToUpper(direction)
	if dir != "ASC" {
		dir = "DESC"
	}

	countQuery := `
		SELECT COUNT(*)
		FROM artists a
		LEFT JOIN artist_stats s ON a.artist_id = s.artist_id
		WHERE 1=1`
	countArgs := make([]interface{}, 0)
	argIdxCount := 1
	if search != "" {
		countQuery += fmt.Sprintf(" AND a.name ILIKE $%d", argIdxCount)
		countArgs = append(countArgs, "%"+search+"%")
		argIdxCount++
	}
	if status == "deleted" {
		countQuery += " AND a.deleted_at IS NOT NULL"
	} else if status != "" {
		countQuery += fmt.Sprintf(" AND a.status = $%d", argIdxCount)
		countArgs = append(countArgs, status)
		argIdxCount++
	}
	if !includeDeleted && status != "deleted" {
		countQuery += " AND a.deleted_at IS NULL"
	}
	if hasReviews != nil {
		if *hasReviews {
			countQuery += " AND COALESCE(s.reviews_count, 0) > 0"
		} else {
			countQuery += " AND COALESCE(s.reviews_count, 0) = 0"
		}
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count artists admin: %w", err)
	}

	// Сортировка по умолчанию: новые сверху
	query += fmt.Sprintf(" ORDER BY %s %s, a.created_at DESC", sortExpr, dir)

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
		return nil, 0, fmt.Errorf("query artists admin: %w", err)
	}
	defer rows.Close()

	var artists []domain.Artist
	for rows.Next() {
		var record ArtistRecord
		err := rows.Scan(
			&record.ArtistID, &record.Name, &record.Description, &record.PhotoURL,
			&record.SocialLinks, &record.StatsReviewsCount, &record.StatsSumRatingTotal,
			&record.StatsConcertsCount, &record.StatsFavoritesCount, &record.StatsUpdatedAt,
			&record.Status, &record.CreatedAt, &record.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan artist row: %w", err)
		}
		artists = append(artists, record.MapToDomain())
	}

	if artists == nil {
		artists = []domain.Artist{}
	}
	return artists, total, nil
}
