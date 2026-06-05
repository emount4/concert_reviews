package moderation_postgres_repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const adminUserSelect = `
	SELECT
		u.user_id, u.email, u.password_hash, u.tg_id, u.tg_username, u.role_id,
		u.username, u.bio, u.avatar_url, u.banner_url, u.is_email_verified,
		u.is_active, u.is_banned, u.banned_by_user_id, u.created_at, u.updated_at,
		COALESCE(s.reviews_count, 0),
		COALESCE(s.likes_given_count, 0),
		COALESCE(s.likes_received_count, 0),
		s.updated_at
	FROM users u
	LEFT JOIN user_stats s ON s.user_id = u.user_id
`

func (r *ModerationRepository) GetUsers(
	ctx context.Context,
	search string,
	limit, offset *int,
) ([]domain.User, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	args := []any{}
	where := ""
	if search = strings.TrimSpace(search); search != "" {
		args = append(args, "%"+search+"%")
		where = " WHERE u.username ILIKE $1 OR u.email ILIKE $1"
	}

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users u"+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count admin users: %w", err)
	}

	query := adminUserSelect + where + " ORDER BY u.created_at DESC, u.user_id ASC"
	if limit != nil && offset != nil {
		args = append(args, *limit, *offset)
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query admin users: %w", err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		user, err := scanAdminUser(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan admin user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate admin users rows: %w", err)
	}

	return users, total, nil
}

func (r *ModerationRepository) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	user, err := scanAdminUser(r.pool.QueryRow(ctx, adminUserSelect+" WHERE u.user_id = $1", id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, core_errors.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("query admin user by id: %w", err)
	}

	return user, nil
}

func (r *ModerationRepository) SetUserBan(
	ctx context.Context,
	moderatorID uuid.UUID,
	targetID uuid.UUID,
	isBanned bool,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.User{}, fmt.Errorf("begin set user ban tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	target, err := lockAdminUser(ctx, tx, targetID)
	if err != nil {
		return domain.User{}, err
	}

	var bannedBy *uuid.UUID
	if isBanned {
		bannedBy = &moderatorID
	}
	if _, err := tx.Exec(ctx, `
		UPDATE users
		SET is_banned = $1, banned_by_user_id = $2
		WHERE user_id = $3
	`, isBanned, bannedBy, targetID); err != nil {
		return domain.User{}, fmt.Errorf("update user ban: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM auth_tokens WHERE user_id = $1`, targetID); err != nil {
		return domain.User{}, fmt.Errorf("delete user sessions after ban change: %w", err)
	}

	action := "user_unbanned"
	if isBanned {
		action = "user_banned"
	}
	if err := insertModerationLog(ctx, tx, moderatorID, action, targetID.String(), domain.LogTargetUser, map[string]any{
		"old_is_banned": target.IsBanned,
		"new_is_banned": isBanned,
		"old_role_id":   target.RoleID,
	}); err != nil {
		return domain.User{}, err
	}

	updated, err := getAdminUserByIDTx(ctx, tx, targetID)
	if err != nil {
		return domain.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, fmt.Errorf("commit set user ban tx: %w", err)
	}

	return updated, nil
}

func (r *ModerationRepository) SetUserRole(
	ctx context.Context,
	moderatorID uuid.UUID,
	targetID uuid.UUID,
	roleID int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.User{}, fmt.Errorf("begin set user role tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	target, err := lockAdminUser(ctx, tx, targetID)
	if err != nil {
		return domain.User{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE users SET role_id = $1 WHERE user_id = $2`, roleID, targetID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return domain.User{}, fmt.Errorf("role not found: %w", core_errors.ErrInvalidArgument)
		}
		return domain.User{}, fmt.Errorf("update user role: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM auth_tokens WHERE user_id = $1`, targetID); err != nil {
		return domain.User{}, fmt.Errorf("delete user sessions after role change: %w", err)
	}

	if err := insertModerationLog(ctx, tx, moderatorID, "user_role_changed", targetID.String(), domain.LogTargetUser, map[string]any{
		"old_role_id": target.RoleID,
		"new_role_id": roleID,
	}); err != nil {
		return domain.User{}, err
	}

	updated, err := getAdminUserByIDTx(ctx, tx, targetID)
	if err != nil {
		return domain.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, fmt.Errorf("commit set user role tx: %w", err)
	}

	return updated, nil
}

func (r *ModerationRepository) AnonymizeUser(
	ctx context.Context,
	moderatorID uuid.UUID,
	targetID uuid.UUID,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.User{}, fmt.Errorf("begin anonymize user tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	target, err := lockAdminUser(ctx, tx, targetID)
	if err != nil {
		return domain.User{}, err
	}
	if target.RoleID == domain.RoleSuperAdminID {
		return domain.User{}, fmt.Errorf("super admin cannot anonymize another super admin: %w", core_errors.ErrForbidden)
	}
	if !target.IsActive {
		return domain.User{}, fmt.Errorf("user is already anonymized: %w", core_errors.ErrConflict)
	}

	deletedEmail := fmt.Sprintf("deleted_%s@concert.bot", targetID.String())
	deletedUsername := fmt.Sprintf("deleted_user_%s", strings.ReplaceAll(targetID.String(), "-", "")[:12])

	if _, err := tx.Exec(ctx, `
		UPDATE users
		SET email = $1,
		    username = $2,
		    password_hash = '',
		    bio = NULL,
		    avatar_url = NULL,
		    banner_url = NULL,
		    tg_id = NULL,
		    tg_username = NULL,
		    is_active = false,
		    is_banned = false,
		    banned_by_user_id = NULL
		WHERE user_id = $3
	`, deletedEmail, deletedUsername, targetID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, fmt.Errorf("anonymized user identity conflict: %w", core_errors.ErrConflict)
		}
		return domain.User{}, fmt.Errorf("anonymize user: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM auth_tokens WHERE user_id = $1`, targetID); err != nil {
		return domain.User{}, fmt.Errorf("delete user sessions after anonymization: %w", err)
	}

	if err := insertModerationLog(ctx, tx, moderatorID, "user_anonymized", targetID.String(), domain.LogTargetUser, map[string]any{
		"old_email":    target.Email,
		"old_username": target.Username,
		"old_role_id":  target.RoleID,
		"new_username": deletedUsername,
	}); err != nil {
		return domain.User{}, err
	}

	updated, err := getAdminUserByIDTx(ctx, tx, targetID)
	if err != nil {
		return domain.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, fmt.Errorf("commit anonymize user tx: %w", err)
	}

	return updated, nil
}

func lockAdminUser(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.User, error) {
	user, err := scanAdminUser(tx.QueryRow(ctx, adminUserSelect+" WHERE u.user_id = $1 FOR UPDATE OF u", id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, core_errors.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("lock admin user: %w", err)
	}
	return user, nil
}

func getAdminUserByIDTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.User, error) {
	user, err := scanAdminUser(tx.QueryRow(ctx, adminUserSelect+" WHERE u.user_id = $1", id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, core_errors.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("get admin user by id: %w", err)
	}
	return user, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanAdminUser(row rowScanner) (domain.User, error) {
	var user domain.User
	var stats domain.UserStats
	var statsUpdatedAt *time.Time
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
		&stats.ReviewsCount,
		&stats.LikesGivenCount,
		&stats.LikesReceivedCount,
		&statsUpdatedAt,
	); err != nil {
		return domain.User{}, err
	}
	if statsUpdatedAt != nil {
		stats.UpdatedAt = *statsUpdatedAt
	}
	user.Stats = &stats
	return user, nil
}

func insertModerationLog(
	ctx context.Context,
	tx pgx.Tx,
	moderatorID uuid.UUID,
	action string,
	targetID string,
	targetType string,
	details map[string]any,
) error {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("marshal moderation log details: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO moderation_logs (moderator_user_id, action, target_id, target_type, details)
		VALUES ($1, $2, $3, $4::target_type_enum, $5::jsonb)
	`, moderatorID, action, targetID, targetType, detailsJSON); err != nil {
		return fmt.Errorf("insert moderation log: %w", err)
	}

	return nil
}
