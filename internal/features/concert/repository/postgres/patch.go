package concert_postgres_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/google/uuid"
)

func (r *ConcertRepository) UpdateConcert(ctx context.Context, id uuid.UUID, patch domain.ConcertPatch) (domain.Concert, error) {
	// Use transaction to handle venue change stats atomically
	var result domain.Concert
	err := r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		// 1. If venue_id is being changed, get current concert first
		if patch.VenueID.Set {
			var oldVenueID int
			queryGetVenue := `SELECT venue_id FROM concerts WHERE concert_id = $1 AND deleted_at IS NULL`
			if err := exec.QueryRow(txCtx, queryGetVenue, id).Scan(&oldVenueID); err != nil {
				return r.handleDBError(err)
			}

			newVenueID := *patch.VenueID.Value

			// 2. Decrement old venue concerts count
			queryDecOldVenue := `UPDATE venue_stats SET concerts_count = GREATEST(0, COALESCE(concerts_count, 0) - 1) WHERE venue_id = $1`
			if _, err := exec.Exec(txCtx, queryDecOldVenue, oldVenueID); err != nil {
				return err
			}

			// 3. Increment new venue concerts count
			queryIncNewVenue := `
				INSERT INTO venue_stats (venue_id, concerts_count) VALUES ($1, 1)
				ON CONFLICT (venue_id) DO UPDATE SET concerts_count = COALESCE(venue_stats.concerts_count, 0) + 1
			`
			if _, err := exec.Exec(txCtx, queryIncNewVenue, newVenueID); err != nil {
				return err
			}
		}

		// 4. Build SET clause for other fields
		var setValues []string
		args := []any{}
		argID := 1

		if patch.Title.Set {
			setValues = append(setValues, fmt.Sprintf("title = $%d", argID))
			args = append(args, patch.Title.Value)
			argID++
		}
		if patch.VenueID.Set {
			setValues = append(setValues, fmt.Sprintf("venue_id = $%d", argID))
			args = append(args, *patch.VenueID.Value)
			argID++
		}
		if patch.Date.Set {
			setValues = append(setValues, fmt.Sprintf("date = $%d", argID))
			args = append(args, patch.Date.Value)
			argID++
		}
		if patch.PosterKey.Set {
			setValues = append(setValues, fmt.Sprintf("poster_url = $%d", argID))
			args = append(args, patch.PosterKey.Value)
			argID++
		}
		if patch.IsVerified.Set {
			setValues = append(setValues, fmt.Sprintf("is_verified = $%d", argID))
			args = append(args, patch.IsVerified.Value)
			argID++
		}

		if len(setValues) == 0 {
			// No updates, just fetch current concert
			var err error
			result, err = r.GetConcertByID(txCtx, id, uuid.Nil)
			return err
		}

		// 5. Execute UPDATE with all modified fields
		args = append(args, id)
		query := fmt.Sprintf(
			"UPDATE concerts SET %s WHERE concert_id = $%d AND deleted_at IS NULL",
			strings.Join(setValues, ", "),
			argID,
		)

		_, err := exec.Exec(txCtx, query, args...)
		if err != nil {
			return r.handleDBError(err)
		}

		var err2 error
		result, err2 = r.GetConcertByID(txCtx, id, uuid.Nil)
		return err2
	})

	if err != nil {
		return domain.Concert{}, err
	}

	return result, nil
}
