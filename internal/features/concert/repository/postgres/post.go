package concert_postgres_repository

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/google/uuid"
)

func (r *ConcertRepository) CreateConcert(
	ctx context.Context,
	concert domain.Concert,
	artists []domain.ConcertArtist,
) (domain.Concert, error) {

	err := r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		queryConcert := `
			INSERT INTO concerts (venue_id, title, date, poster_url, is_verified)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING concert_id, created_at
		`
		err := exec.QueryRow(txCtx, queryConcert,
			concert.VenueID, concert.Title, concert.Date, concert.PosterKey, concert.IsVerified,
		).Scan(&concert.ConcertID, &concert.CreatedAt)

		if err != nil {
			return r.handleDBError(err)
		}

		queryArtistLink := `INSERT INTO concert_artists (concert_id, artist_id, is_main) VALUES ($1, $2, $3)`
		for _, a := range artists {
			if _, err := exec.Exec(txCtx, queryArtistLink, concert.ConcertID, a.ArtistID, a.IsMain); err != nil {
				return r.handleDBError(err)
			}
		}

		queryVenueStats := `
			INSERT INTO venue_stats (venue_id, concerts_count) VALUES ($1, 1)
			ON CONFLICT (venue_id) DO UPDATE SET concerts_count = COALESCE(venue_stats.concerts_count, 0) + 1
		`
		if _, err := exec.Exec(txCtx, queryVenueStats, concert.VenueID); err != nil {
			return fmt.Errorf("update venue concerts count: %w", err)
		}

		queryArtistStats := `
			INSERT INTO artist_stats (artist_id, concerts_count) VALUES ($1, 1)
			ON CONFLICT (artist_id) DO UPDATE SET concerts_count = COALESCE(artist_stats.concerts_count, 0) + 1
		`
		for _, a := range artists {
			if a.IsMain {
				if _, err := exec.Exec(txCtx, queryArtistStats, a.ArtistID); err != nil {
					return fmt.Errorf("update artist %d concerts count: %w", a.ArtistID, err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return domain.Concert{}, err
	}

	return r.GetConcertByID(ctx, concert.ConcertID, uuid.Nil)
}

func (r *ConcertRepository) CreateSuggestion(ctx context.Context, s domain.ConcertSuggestion) (domain.ConcertSuggestion, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO concert_suggestions (user_id, raw_artist_name, raw_venue_name, concert_date, info)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING suggestion_id,  created_at
	`

	err := r.pool.QueryRow(ctx, query,
		s.UserID, s.RawArtistName, s.RawVenueName, s.ConcertDate, s.Info,
	).Scan(&s.SuggestionID, &s.CreatedAt)

	if err != nil {
		return domain.ConcertSuggestion{}, fmt.Errorf("insert suggestion: %w", err)
	}

	return s, nil
}

func (r *ConcertRepository) CountPendingSuggestions(ctx context.Context, userID uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT COUNT(*) 
		FROM concert_suggestions 
		WHERE user_id = $1 AND status = 'pending'
	`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count pending suggestions: %w", err)
	}

	return count, nil
}
