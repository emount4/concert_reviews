package auth_postgres_repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *AuthRepository) CreateSession(
	ctx context.Context,
	response core_domain.AuthResponse,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	hashBytes := sha256.Sum256([]byte(response.RefreshToken))
	hash := hex.EncodeToString(hashBytes[:])

	query := `
	INSERT INTO auth_tokens (user_id, token, expires_at)
	VALUES ($1, $2, $3)
	`

	if _, err := exec.Exec(ctx, query, response.User.ID, hash, response.ExpiresAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: token already exists", core_errors.ErrConflict)
		}
		return fmt.Errorf("insert auth token: %w", err)
	}

	return nil
}
