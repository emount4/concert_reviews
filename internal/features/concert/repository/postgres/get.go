package concert_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const baseConcertSelect = `
	SELECT 
		c.concert_id, c.venue_id, c.title, c.date, c.poster_url, c.is_verified, c.created_by_user_id, c.created_at, c.deleted_at,
		v.name as venue_name,
		ct.city_id, ct.name as city_name, ct.slug as city_slug,
		cs.reviews_count, cs.sum_p1, cs.sum_p2, cs.sum_p3, cs.sum_p4, cs.sum_p5, cs.sum_rating_total, cs.updated_at as stats_updated_at,
		(
			SELECT jsonb_agg(jsonb_build_object(
				'ArtistID', a.artist_id,
				'Name', a.name,
				'IsMain', ca.is_main
			))
			FROM concert_artists ca
			JOIN artists a ON a.artist_id = ca.artist_id
			WHERE ca.concert_id = c.concert_id
		) as artists_json
	FROM concerts c
	JOIN venues v ON c.venue_id = v.venue_id
	JOIN cities ct ON v.city_id = ct.city_id
	LEFT JOIN concert_stats cs ON c.concert_id = cs.concert_id
`

func (r *ConcertRepository) GetConcertByID(ctx context.Context, id uuid.UUID) (domain.Concert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := baseConcertSelect + ` WHERE c.concert_id = $1 AND c.deleted_at IS NULL`

	var rec concertRecord
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&rec.ConcertID, &rec.VenueID, &rec.Title, &rec.Date, &rec.PosterKey, &rec.IsVerified, &rec.CreatedByUserID, &rec.CreatedAt, &rec.DeletedAt,
		&rec.VenueName, &rec.CityID, &rec.CityName, &rec.CitySlug,
		&rec.ReviewsCount, &rec.SumP1, &rec.SumP2, &rec.SumP3, &rec.SumP4, &rec.SumP5, &rec.SumRatingTotal, &rec.StatsUpdatedAt,
		&rec.ArtistsJSON,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Concert{}, core_errors.ErrNotFound
		}
		return domain.Concert{}, fmt.Errorf("scan concert: %w", err)
	}

	return rec.MapToDomain(), nil
}

func (r *ConcertRepository) GetConcerts(
	ctx context.Context,
	cityID *int,
	artistID *int,
	search string,
	sort string,
	direction string,
	limit, offset *int,
) ([]domain.Concert, int, error) {
	return r.listConcerts(ctx, cityID, artistID, search, sort, direction, limit, offset, false)
}

func (r *ConcertRepository) GetConcertsAdmin(
	ctx context.Context,
	cityID *int,
	artistID *int,
	search string,
	sort string,
	direction string,
	limit, offset *int,
	includeDeleted bool,
) ([]domain.Concert, int, error) {
	return r.listConcerts(ctx, cityID, artistID, search, sort, direction, limit, offset, includeDeleted)
}

// Приватный метод для объединения логики фильтрации
func (r *ConcertRepository) listConcerts(
	ctx context.Context,
	cityID *int,
	artistID *int,
	search string,
	sort string,
	direction string,
	limit, offset *int,
	includeDeleted bool,
) ([]domain.Concert, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var conditions []string
	args := []any{}

	// Обработка мягкого удаления
	if !includeDeleted {
		conditions = append(conditions, "c.deleted_at IS NULL")
	}

	// Фильтр по городу
	if cityID != nil {
		args = append(args, *cityID)
		conditions = append(conditions, fmt.Sprintf("v.city_id = $%d", len(args)))
	}

	// Поиск по названию
	if search != "" {
		args = append(args, "%"+search+"%")
		conditions = append(conditions, fmt.Sprintf("c.title ILIKE $%d", len(args)))
	}

	// Фильтр по артисту (через подзапрос, так как связь многие-ко-многим)
	if artistID != nil {
		args = append(args, *artistID)
		conditions = append(conditions, fmt.Sprintf(`c.concert_id IN (SELECT concert_id FROM concert_artists WHERE artist_id = $%d)`, len(args)))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	sortExpr := "c.date"
	switch sort {
	case "title":
		sortExpr = "c.title"
	case "rating":
		sortExpr = "COALESCE(CASE WHEN COALESCE(cs.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(cs.sum_rating_total, 0)::numeric / cs.reviews_count END, 0)"
	case "reviews":
		sortExpr = "COALESCE(cs.reviews_count, 0)"
	case "p1":
		sortExpr = "COALESCE(CASE WHEN COALESCE(cs.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(cs.sum_p1, 0)::numeric / cs.reviews_count END, 0)"
	case "p2":
		sortExpr = "COALESCE(CASE WHEN COALESCE(cs.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(cs.sum_p2, 0)::numeric / cs.reviews_count END, 0)"
	case "p3":
		sortExpr = "COALESCE(CASE WHEN COALESCE(cs.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(cs.sum_p3, 0)::numeric / cs.reviews_count END, 0)"
	case "p4":
		sortExpr = "COALESCE(CASE WHEN COALESCE(cs.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(cs.sum_p4, 0)::numeric / cs.reviews_count END, 0)"
	case "p5":
		sortExpr = "COALESCE(CASE WHEN COALESCE(cs.reviews_count, 0) = 0 THEN 0 ELSE COALESCE(cs.sum_p5, 0)::numeric / cs.reviews_count END, 0)"
	}

	dir := strings.ToUpper(direction)
	if dir != "ASC" {
		dir = "DESC"
	}

	countQuery := `SELECT COUNT(*) FROM concerts c JOIN venues v ON c.venue_id = v.venue_id JOIN cities ct ON v.city_id = ct.city_id` + whereClause
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count list: %w", err)
	}

	query := baseConcertSelect + whereClause + fmt.Sprintf(" ORDER BY %s %s, c.date DESC", sortExpr, dir)

	if limit != nil && offset != nil {
		args = append(args, *limit, *offset)
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query list: %w", err)
	}
	defer rows.Close()

	var result []domain.Concert
	for rows.Next() {
		var rec concertRecord
		err := rows.Scan(
			&rec.ConcertID, &rec.VenueID, &rec.Title, &rec.Date, &rec.PosterKey, &rec.IsVerified, &rec.CreatedByUserID, &rec.CreatedAt, &rec.DeletedAt,
			&rec.VenueName, &rec.CityID, &rec.CityName, &rec.CitySlug,
			&rec.ReviewsCount, &rec.SumP1, &rec.SumP2, &rec.SumP3, &rec.SumP4, &rec.SumP5, &rec.SumRatingTotal, &rec.StatsUpdatedAt,
			&rec.ArtistsJSON,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan list row: %w", err)
		}
		result = append(result, rec.MapToDomain())
	}

	if result == nil {
		result = []domain.Concert{}
	}
	return result, total, nil
}

// GetSuggestionsAdmin — Список предложений для модерации
func (r *ConcertRepository) GetSuggestionsAdmin(ctx context.Context, limit, offset *int, status string) ([]domain.ConcertSuggestion, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT suggestion_id, user_id, raw_artist_name, raw_venue_name, concert_date, info, created_at
		FROM concert_suggestions
		ORDER BY created_at DESC
	`
	args := []any{}

	if limit != nil && offset != nil {
		argNum := len(args) + 1
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argNum, argNum+1)
		args = append(args, *limit, *offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query suggestions: %w", err)
	}
	defer rows.Close()

	var result []domain.ConcertSuggestion
	for rows.Next() {
		var s domain.ConcertSuggestion
		err := rows.Scan(&s.SuggestionID, &s.UserID, &s.RawArtistName, &s.RawVenueName, &s.ConcertDate, &s.Info, &s.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan suggestion: %w", err)
		}
		result = append(result, s)
	}

	if result == nil {
		result = []domain.ConcertSuggestion{}
	}
	return result, nil
}

// GetSuggestionByIDAdmin — Детали предложения
func (r *ConcertRepository) GetSuggestionByIDAdmin(ctx context.Context, id uuid.UUID) (domain.ConcertSuggestion, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT suggestion_id, user_id, raw_artist_name, raw_venue_name, concert_date, info, created_at
		FROM concert_suggestions
		WHERE suggestion_id = $1
	`
	var s domain.ConcertSuggestion
	err := r.pool.QueryRow(ctx, query, id).Scan(&s.SuggestionID, &s.UserID, &s.RawArtistName, &s.RawVenueName, &s.ConcertDate, &s.Info, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ConcertSuggestion{}, core_errors.ErrNotFound
		}
		return domain.ConcertSuggestion{}, err
	}
	return s, nil
}
