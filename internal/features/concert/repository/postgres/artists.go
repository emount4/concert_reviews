package concert_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/google/uuid"
)

// AddConcertArtist adds an artist to a concert and updates artist stats if is_main=true
func (r *ConcertRepository) AddConcertArtist(ctx context.Context, concertID uuid.UUID, artistID int, isMain bool) error {
	return r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		// 1. Check if concert exists
		var concertExists bool
		queryCheckConcert := `SELECT EXISTS(SELECT 1 FROM concerts WHERE concert_id = $1 AND deleted_at IS NULL)`
		if err := exec.QueryRow(txCtx, queryCheckConcert, concertID).Scan(&concertExists); err != nil {
			return err
		}
		if !concertExists {
			return core_errors.ErrNotFound
		}

		// 2. Check if artist exists
		var artistExists bool
		queryCheckArtist := `SELECT EXISTS(SELECT 1 FROM artists WHERE artist_id = $1 AND deleted_at IS NULL)`
		if err := exec.QueryRow(txCtx, queryCheckArtist, artistID).Scan(&artistExists); err != nil {
			return err
		}
		if !artistExists {
			return core_errors.ErrNotFound
		}

		// 3. Check if artist is already linked to this concert
		var alreadyLinked bool
		queryCheckLink := `SELECT EXISTS(SELECT 1 FROM concert_artists WHERE concert_id = $1 AND artist_id = $2)`
		if err := exec.QueryRow(txCtx, queryCheckLink, concertID, artistID).Scan(&alreadyLinked); err != nil {
			return err
		}
		if alreadyLinked {
			return fmt.Errorf("artist already linked to this concert")
		}

		// 4. Insert concert_artist link
		queryInsertLink := `INSERT INTO concert_artists (concert_id, artist_id, is_main) VALUES ($1, $2, $3)`
		if _, err := exec.Exec(txCtx, queryInsertLink, concertID, artistID, isMain); err != nil {
			return err
		}

		// 5. If is_main=true, increment artist concerts count
		if isMain {
			queryIncArtistStats := `
				INSERT INTO artist_stats (artist_id, concerts_count) VALUES ($1, 1)
				ON CONFLICT (artist_id) DO UPDATE SET concerts_count = COALESCE(artist_stats.concerts_count, 0) + 1
			`
			if _, err := exec.Exec(txCtx, queryIncArtistStats, artistID); err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveConcertArtist removes an artist from a concert and updates artist stats if was_main=true
func (r *ConcertRepository) RemoveConcertArtist(ctx context.Context, concertID uuid.UUID, artistID int) error {
	return r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		// 1. Check if link exists and get is_main flag
		var isMain bool
		queryGetLink := `SELECT is_main FROM concert_artists WHERE concert_id = $1 AND artist_id = $2`
		err := exec.QueryRow(txCtx, queryGetLink, concertID, artistID).Scan(&isMain)
		if err != nil {
			return core_errors.ErrNotFound
		}

		// 2. Delete concert_artist link
		queryDeleteLink := `DELETE FROM concert_artists WHERE concert_id = $1 AND artist_id = $2`
		if _, err := exec.Exec(txCtx, queryDeleteLink, concertID, artistID); err != nil {
			return err
		}

		// 3. If was_main=true, decrement artist concerts count
		if isMain {
			queryDecArtistStats := `UPDATE artist_stats SET concerts_count = GREATEST(0, COALESCE(concerts_count, 0) - 1) WHERE artist_id = $1`
			if _, err := exec.Exec(txCtx, queryDecArtistStats, artistID); err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateConcertArtistIsMain updates the is_main flag for an artist in a concert
func (r *ConcertRepository) UpdateConcertArtistIsMain(ctx context.Context, concertID uuid.UUID, artistID int, isMain bool) error {
	return r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		// 1. Get current is_main flag
		var currentIsMain bool
		queryGetLink := `SELECT is_main FROM concert_artists WHERE concert_id = $1 AND artist_id = $2`
		err := exec.QueryRow(txCtx, queryGetLink, concertID, artistID).Scan(&currentIsMain)
		if err != nil {
			return core_errors.ErrNotFound
		}

		// If no change needed, return
		if currentIsMain == isMain {
			return nil
		}

		// 2. Update is_main flag
		queryUpdateLink := `UPDATE concert_artists SET is_main = $1 WHERE concert_id = $2 AND artist_id = $3`
		if _, err := exec.Exec(txCtx, queryUpdateLink, isMain, concertID, artistID); err != nil {
			return err
		}

		// 3. Update artist stats
		if isMain {
			// Transitioning from non-main to main: increment
			queryIncArtistStats := `
				INSERT INTO artist_stats (artist_id, concerts_count) VALUES ($1, 1)
				ON CONFLICT (artist_id) DO UPDATE SET concerts_count = COALESCE(artist_stats.concerts_count, 0) + 1
			`
			if _, err := exec.Exec(txCtx, queryIncArtistStats, artistID); err != nil {
				return err
			}
		} else {
			// Transitioning from main to non-main: decrement
			queryDecArtistStats := `UPDATE artist_stats SET concerts_count = GREATEST(0, COALESCE(concerts_count, 0) - 1) WHERE artist_id = $1`
			if _, err := exec.Exec(txCtx, queryDecArtistStats, artistID); err != nil {
				return err
			}
		}

		return nil
	})
}
