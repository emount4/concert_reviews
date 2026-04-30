package auth_postgres_repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	hash := r.hashToken(response.RefreshToken)

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

func (r *AuthRepository) GetSession(
	ctx context.Context,
	token string,
) (core_domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	hashedToken := r.hashToken(token)

	query := `
	SELECT user_id, token, expires_at 
	FROM auth_tokens 
	WHERE token = $1
	`

	var session core_domain.RefreshToken
	err := exec.QueryRow(ctx, query, hashedToken).Scan(
		&session.UserID,
		&session.Token,
		&session.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.RefreshToken{}, core_errors.ErrNotFound
		}
		return core_domain.RefreshToken{}, fmt.Errorf("select session: %w", err)
	}

	// Проверка на протухание прямо в репозитории
	if time.Now().After(session.ExpiresAt) {
		_ = r.DeleteSession(ctx, token) // удаляем протухший
		return core_domain.RefreshToken{}, core_errors.ErrUnauthorized
	}

	return session, nil
}

func (r *AuthRepository) DeleteSession(
	ctx context.Context,
	token string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	hashedToken := r.hashToken(token)
	query := `DELETE FROM auth_tokens WHERE token = $1`
	if _, err := exec.Exec(ctx, query, hashedToken); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (r *AuthRepository) DeleteAllUserSessions(
	ctx context.Context,
	userID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `DELETE FROM auth_tokens WHERE user_id = $1`
	if _, err := exec.Exec(ctx, query, userID); err != nil {
		return fmt.Errorf("delete user sessions: %w", err)
	}
	return nil
}

func (r *AuthRepository) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
