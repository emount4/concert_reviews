package venue_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/jackc/pgx/v5"
)

// GetVenueByID — получение площадки по ID (исключаем удалённые)
func (r *VenueRepository) GetVenueByID(ctx context.Context, id int) (domain.Venue, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT v.venue_id, v.city_id, c.name, c.slug, v.name, v.address, v.capacity, v.social_links,
		       COALESCE(vs.reviews_count, 0) AS reviews_count,
		       COALESCE(vs.sum_rating_total, 0) AS sum_rating_total,
		       COALESCE(vs.concerts_count, 0) AS concerts_count,
		       COALESCE(vs.favorites_count, 0) AS favorites_count,
		       COALESCE(vs.updated_at, v.created_at) AS stats_updated_at,
		       v.photo_url, v.description, v.status, v.created_at, v.deleted_at
		FROM venues v
		LEFT JOIN cities c ON v.city_id = c.city_id
		LEFT JOIN venue_stats vs ON v.venue_id = vs.venue_id
		WHERE v.venue_id = $1 AND v.deleted_at IS NULL
	`

	var rec VenueRecord
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&rec.VenueID, &rec.CityID, &rec.CityName, &rec.CitySlug, &rec.Name, &rec.Address, &rec.Capacity, &rec.SocialLinks,
		&rec.StatsReviewsCount, &rec.StatsSumRatingTotal, &rec.StatsConcertsCount, &rec.StatsFavoritesCount, &rec.StatsUpdatedAt, &rec.PhotoURL,
		&rec.Description, &rec.Status, &rec.CreatedAt, &rec.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Venue{}, core_errors.ErrNotFound
		}
		return domain.Venue{}, fmt.Errorf("select venue by id: %w", err)
	}

	return rec.MapToDomain(), nil
}

// GetVenues — список с пагинацией, поиском и фильтром по городу
func (r *VenueRepository) GetVenues(
	ctx context.Context,
	cityID *int,
	search string,
	limit, offset *int,
) ([]domain.Venue, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	// Базовый запрос (только активные)
	query := `
		SELECT v.venue_id, v.city_id, c.name, c.slug, v.name, v.address, v.capacity, v.social_links,
		       COALESCE(vs.reviews_count, 0) AS reviews_count,
		       COALESCE(vs.sum_rating_total, 0) AS sum_rating_total,
		       COALESCE(vs.concerts_count, 0) AS concerts_count,
		       COALESCE(vs.favorites_count, 0) AS favorites_count,
		       COALESCE(vs.updated_at, v.created_at) AS stats_updated_at,
		       v.photo_url, v.description, v.status, v.created_at, v.deleted_at
		FROM venues v
		LEFT JOIN cities c ON v.city_id = c.city_id
		LEFT JOIN venue_stats vs ON v.venue_id = vs.venue_id
		WHERE v.deleted_at IS NULL
	`
	args := make([]any, 0)
	argIdx := 1

	// Фильтр по городу
	if cityID != nil && *cityID > 0 {
		query += fmt.Sprintf(" AND v.city_id = $%d", argIdx)
		args = append(args, *cityID)
		argIdx++
	}

	// Поиск по имени (регистронезависимый)
	if search != "" {
		query += fmt.Sprintf(" AND v.name ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	// Сортировка по умолчанию: по имени
	query += ` ORDER BY v.name ASC`

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
		return nil, fmt.Errorf("query venues list: %w", err)
	}
	defer rows.Close()

	var venues []domain.Venue
	for rows.Next() {
		var rec VenueRecord
		err := rows.Scan(
			&rec.VenueID, &rec.CityID, &rec.CityName, &rec.CitySlug, &rec.Name, &rec.Address, &rec.Capacity, &rec.SocialLinks,
			&rec.StatsReviewsCount, &rec.StatsSumRatingTotal, &rec.StatsConcertsCount, &rec.StatsFavoritesCount, &rec.StatsUpdatedAt, &rec.PhotoURL,
			&rec.Description, &rec.Status, &rec.CreatedAt, &rec.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan venue: %w", err)
		}
		venues = append(venues, rec.MapToDomain())
	}

	return venues, nil
}

// GetVenuesAdmin — расширенный список для админки (с флагом includeDeleted и фильтром по статусу)
func (r *VenueRepository) GetVenuesAdmin(
	ctx context.Context,
	cityID *int,
	search string,
	limit, offset *int,
	includeDeleted bool,
	status string,
) ([]domain.Venue, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `
		SELECT v.venue_id, v.city_id, c.name, c.slug, v.name, v.address, v.capacity, v.social_links,
		       COALESCE(vs.reviews_count, 0) AS reviews_count,
		       COALESCE(vs.sum_rating_total, 0) AS sum_rating_total,
		       COALESCE(vs.concerts_count, 0) AS concerts_count,
		       COALESCE(vs.favorites_count, 0) AS favorites_count,
		       COALESCE(vs.updated_at, v.created_at) AS stats_updated_at,
		       v.photo_url, v.description, v.status, v.created_at, v.deleted_at
		FROM venues v
		LEFT JOIN cities c ON v.city_id = c.city_id
		LEFT JOIN venue_stats vs ON v.venue_id = vs.venue_id
		WHERE 1=1
	`
	args := make([]any, 0)
	argIdx := 1

	// Фильтр по городу
	if cityID != nil && *cityID > 0 {
		query += fmt.Sprintf(" AND v.city_id = $%d", argIdx)
		args = append(args, *cityID)
		argIdx++
	}

	// Поиск по имени
	if search != "" {
		query += fmt.Sprintf(" AND v.name ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	// Фильтр по статусу
	if status != "" {
		query += fmt.Sprintf(" AND v.status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	// Фильтр по deleted_at (если не включены удалённые)
	if !includeDeleted {
		query += " AND v.deleted_at IS NULL"
	}

	// Сортировка: новые сверху
	query += " ORDER BY v.created_at DESC"

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
		return nil, fmt.Errorf("query venues admin: %w", err)
	}
	defer rows.Close()

	var venues []domain.Venue
	for rows.Next() {
		var rec VenueRecord
		err := rows.Scan(
			&rec.VenueID, &rec.CityID, &rec.CityName, &rec.CitySlug, &rec.Name, &rec.Address, &rec.Capacity, &rec.SocialLinks,
			&rec.StatsReviewsCount, &rec.StatsSumRatingTotal, &rec.StatsConcertsCount, &rec.StatsFavoritesCount, &rec.StatsUpdatedAt, &rec.PhotoURL,
			&rec.Description, &rec.Status, &rec.CreatedAt, &rec.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan venue row: %w", err)
		}
		venues = append(venues, rec.MapToDomain())
	}

	return venues, nil
}
