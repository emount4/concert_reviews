package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *AuthRepository) CreateUser(
	ctx context.Context,
	user core_domain.User,
) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `
	INSERT INTO users (email, password_hash, username)
	VALUES ($1, $2, $3)
	RETURNING user_id, role_id, is_email_verified, is_active, is_banned, created_at, updated_at
	`

	row := exec.QueryRow(ctx, query, user.Email, user.PasswordHash, user.Username)
	if err := row.Scan(
		&user.ID,
		&user.RoleID,
		&user.IsEmailVerified,
		&user.IsActive,
		&user.IsBanned,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return core_domain.User{}, fmt.Errorf("%w: user already exists", core_errors.ErrConflict)
		}
		return core_domain.User{}, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}

func (r *AuthRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `
	SELECT
		user_id,
		email,
		password_hash,
		tg_id,
		tg_username,
		role_id,
		username,
		bio,
		avatar_url,
		banner_url,
		is_email_verified,
		is_active,
		is_banned,
		banned_by_user_id,
		created_at,
		updated_at
	FROM users
	WHERE email = $1
	LIMIT 1
	`

	var user core_domain.User
	row := exec.QueryRow(ctx, query, email)
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.TelegramID,
		&user.TelegramUsername,
		&user.RoleID,
		&user.Username,
		&user.Bio,
		&user.AvatarURL,
		&user.BannerURL,
		&user.IsEmailVerified,
		&user.IsActive,
		&user.IsBanned,
		&user.BannedByUserID,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.User{}, core_errors.ErrNotFound
		}
		return core_domain.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

func (r *AuthRepository) GetUserByTelegramID(
	ctx context.Context,
	tgID int64,
) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	query := `
	SELECT
		user_id,
		email,
		password_hash,
		tg_id,
		tg_username,
		role_id,
		username,
		bio,
		avatar_url,
		banner_url,
		is_email_verified,
		is_active,
		is_banned,
		banned_by_user_id,
		created_at,
		updated_at
	FROM users
	WHERE tg_id = $1
	LIMIT 1
	`

	var user core_domain.User
	row := exec.QueryRow(ctx, query, tgID)
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.TelegramID,
		&user.TelegramUsername,
		&user.RoleID,
		&user.Username,
		&user.Bio,
		&user.AvatarURL,
		&user.BannerURL,
		&user.IsEmailVerified,
		&user.IsActive,
		&user.IsBanned,
		&user.BannedByUserID,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_domain.User{}, core_errors.ErrNotFound
		}
		return core_domain.User{}, fmt.Errorf("get user by telegram id: %w", err)
	}

	return user, nil
}
