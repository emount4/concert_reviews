package concert_postgres_repository

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/google/uuid"
)

func (r *ConcertRepository) DeleteConcertSoft(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `UPDATE concerts SET deleted_at = NOW() WHERE concert_id = $1 AND deleted_at IS NULL`
	res, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}
	return nil
}

// RestoreConcert — Снимает флаг удаления (статистика не меняется, т.к. soft delete её не трогает)
func (r *ConcertRepository) RestoreConcert(ctx context.Context, id uuid.UUID) (domain.Concert, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `UPDATE concerts SET deleted_at = NULL WHERE concert_id = $1 AND deleted_at IS NOT NULL`
	res, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return domain.Concert{}, err
	}
	if res.RowsAffected() == 0 {
		return domain.Concert{}, core_errors.ErrNotFound
	}
	return r.GetConcertByID(ctx, id, uuid.Nil)
}

func (r *ConcertRepository) DeleteConcertHard(ctx context.Context, id uuid.UUID) error {
	// Используем транзакцию, так как нужно обновить статы в других таблицах
	return r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		// 1. Получаем ID площадки перед удалением
		var venueID int
		queryGetVenue := `SELECT venue_id FROM concerts WHERE concert_id = $1`
		if err := exec.QueryRow(txCtx, queryGetVenue, id).Scan(&venueID); err != nil {
			return r.handleDBError(err)
		}

		// 2. Получаем ID артистов с is_main=true перед удалением
		var mainArtistIDs []int
		queryGetMainArtists := `
			SELECT ARRAY(SELECT artist_id FROM concert_artists WHERE concert_id = $1 AND is_main = true)
		`
		if err := exec.QueryRow(txCtx, queryGetMainArtists, id).Scan(&mainArtistIDs); err != nil {
			return r.handleDBError(err)
		}

		// 3. Декремент счетчика у площадки
		queryDecVenue := `UPDATE venue_stats SET concerts_count = GREATEST(0, COALESCE(concerts_count, 0) - 1) WHERE venue_id = $1`
		if _, err := exec.Exec(txCtx, queryDecVenue, venueID); err != nil {
			return err
		}

		// 4. Декремент счетчика только у основных артистов (is_main=true)
		if len(mainArtistIDs) > 0 {
			queryDecArtist := `UPDATE artist_stats SET concerts_count = GREATEST(0, COALESCE(concerts_count, 0) - 1) WHERE artist_id = ANY($1)`
			if _, err := exec.Exec(txCtx, queryDecArtist, mainArtistIDs); err != nil {
				return err
			}
		}

		// 5. Физическое удаление (Связи concert_artists и concert_stats удалятся сами по ON DELETE CASCADE)
		queryDelete := `DELETE FROM concerts WHERE concert_id = $1`
		res, err := exec.Exec(txCtx, queryDelete, id)
		if err != nil {
			return r.handleDBError(err)
		}
		if res.RowsAffected() == 0 {
			return core_errors.ErrNotFound
		}

		return nil
	})
}

func (r *ConcertRepository) DeleteSuggestionAdmin(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM concert_suggestions WHERE suggestion_id = $1`
	res, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete suggestion: %w", err)
	}
	if res.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}
	return nil
}
