package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *AuthRepository) LinkTG(
	ctx context.Context,
	userID uuid.UUID,
	username string,
	tgID int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	sql := `SELECT user_id FROM users WHERE tg_id = $1 AND user_id != $2`

	var existingUserID uuid.UUID
	err := exec.QueryRow(ctx, sql, tgID, userID).Scan(&existingUserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("check tg_id uniqueness: %w", err)
	}
	if err == nil {
		return fmt.Errorf("%w: telegram account already linked to another user", core_errors.ErrConflict)
	}

	sql = `
		UPDATE users
		SET tg_id = $1, tg_username = $2, updated_at = NOW()
		WHERE user_id = $3
	`

	ct, err := exec.Exec(ctx, sql, tgID, username, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("telegram account already linked to another user, %w", core_errors.ErrConflict)
		}
		return fmt.Errorf("link telegram: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}

	return nil
}
