package moderation_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *ModerationRepository) GetActiveProfileRequests(
	ctx context.Context,
	limit, offset *int,
) ([]domain.ProfileModerationRequest, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const countQuery = `
		SELECT COUNT(*)
		FROM profile_moderation pm
		WHERE pm.status = 'pending'
	`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count active profile moderation requests: %w", err)
	}

	query := `
		SELECT pm.moderation_id, pm.user_id, u.username, u.avatar_url,
		       pm.field_name, pm.old_value, pm.new_value, pm.status,
		       pm.moderated_by_user_id, pm.created_at, pm.updated_at
		FROM profile_moderation pm
		JOIN users u ON u.user_id = pm.user_id
		WHERE pm.status = 'pending'
		ORDER BY pm.created_at ASC, pm.moderation_id ASC
	`
	args := []any{}
	if limit != nil && offset != nil {
		query += " LIMIT $1 OFFSET $2"
		args = append(args, *limit, *offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query active profile moderation requests: %w", err)
	}
	defer rows.Close()

	requests := make([]domain.ProfileModerationRequest, 0)
	for rows.Next() {
		var req domain.ProfileModerationRequest
		var status string
		if err := rows.Scan(
			&req.ModerationID,
			&req.UserID,
			&req.Username,
			&req.UserAvatarURL,
			&req.FieldName,
			&req.OldValue,
			&req.NewValue,
			&status,
			&req.ModeratedByUserID,
			&req.CreatedAt,
			&req.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan active profile moderation request: %w", err)
		}
		req.Status = domain.ModerationStatus(status)
		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate active profile moderation requests: %w", err)
	}

	return requests, total, nil
}

func (r *ModerationRepository) ApproveProfileRequest(
	ctx context.Context,
	id int,
	moderatorID uuid.UUID,
) (domain.ProfileModerationRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.ProfileModerationRequest{}, fmt.Errorf("begin approve profile request tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	req, err := lockProfileRequest(ctx, tx, id)
	if err != nil {
		return domain.ProfileModerationRequest{}, err
	}
	if req.Status != domain.StatusPending {
		return domain.ProfileModerationRequest{}, fmt.Errorf("profile moderation request is not pending: %w", core_errors.ErrConflict)
	}

	column, err := profileFieldColumn(req.FieldName)
	if err != nil {
		return domain.ProfileModerationRequest{}, err
	}
	if column == "username" && req.NewValue == nil {
		return domain.ProfileModerationRequest{}, fmt.Errorf("username cannot be null: %w", core_errors.ErrInvalidArgument)
	}

	queryUpdateUser := fmt.Sprintf(`UPDATE users SET %s = $1 WHERE user_id = $2`, column)
	if _, err := tx.Exec(ctx, queryUpdateUser, req.NewValue, req.UserID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ProfileModerationRequest{}, fmt.Errorf("%w: profile value conflicts with existing user", core_errors.ErrConflict)
		}
		return domain.ProfileModerationRequest{}, fmt.Errorf("apply profile moderation request: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE profile_moderation
		SET status = 'approved', moderated_by_user_id = $1, updated_at = NOW()
		WHERE moderation_id = $2
	`, moderatorID, id); err != nil {
		return domain.ProfileModerationRequest{}, fmt.Errorf("mark profile request approved: %w", err)
	}

	updated, err := getProfileRequestByID(ctx, tx, id)
	if err != nil {
		return domain.ProfileModerationRequest{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.ProfileModerationRequest{}, fmt.Errorf("commit approve profile request tx: %w", err)
	}

	return updated, nil
}

func (r *ModerationRepository) RejectProfileRequest(
	ctx context.Context,
	id int,
	moderatorID uuid.UUID,
) (domain.ProfileModerationRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.ProfileModerationRequest{}, fmt.Errorf("begin reject profile request tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	req, err := lockProfileRequest(ctx, tx, id)
	if err != nil {
		return domain.ProfileModerationRequest{}, err
	}
	if req.Status != domain.StatusPending {
		return domain.ProfileModerationRequest{}, fmt.Errorf("profile moderation request is not pending: %w", core_errors.ErrConflict)
	}

	if _, err := tx.Exec(ctx, `
		UPDATE profile_moderation
		SET status = 'rejected', moderated_by_user_id = $1, updated_at = NOW()
		WHERE moderation_id = $2
	`, moderatorID, id); err != nil {
		return domain.ProfileModerationRequest{}, fmt.Errorf("mark profile request rejected: %w", err)
	}

	updated, err := getProfileRequestByID(ctx, tx, id)
	if err != nil {
		return domain.ProfileModerationRequest{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.ProfileModerationRequest{}, fmt.Errorf("commit reject profile request tx: %w", err)
	}

	return updated, nil
}

func lockProfileRequest(ctx context.Context, tx pgx.Tx, id int) (domain.ProfileModerationRequest, error) {
	query := `
		SELECT pm.moderation_id, pm.user_id, u.username, u.avatar_url,
		       pm.field_name, pm.old_value, pm.new_value, pm.status,
		       pm.moderated_by_user_id, pm.created_at, pm.updated_at
		FROM profile_moderation pm
		JOIN users u ON u.user_id = pm.user_id
		WHERE pm.moderation_id = $1
		FOR UPDATE OF pm
	`
	req, err := scanProfileRequest(tx.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ProfileModerationRequest{}, core_errors.ErrNotFound
		}
		return domain.ProfileModerationRequest{}, fmt.Errorf("lock profile moderation request: %w", err)
	}
	return req, nil
}

func getProfileRequestByID(ctx context.Context, tx pgx.Tx, id int) (domain.ProfileModerationRequest, error) {
	query := `
		SELECT pm.moderation_id, pm.user_id, u.username, u.avatar_url,
		       pm.field_name, pm.old_value, pm.new_value, pm.status,
		       pm.moderated_by_user_id, pm.created_at, pm.updated_at
		FROM profile_moderation pm
		JOIN users u ON u.user_id = pm.user_id
		WHERE pm.moderation_id = $1
	`
	req, err := scanProfileRequest(tx.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ProfileModerationRequest{}, core_errors.ErrNotFound
		}
		return domain.ProfileModerationRequest{}, fmt.Errorf("get profile moderation request: %w", err)
	}
	return req, nil
}

func scanProfileRequest(row pgx.Row) (domain.ProfileModerationRequest, error) {
	var req domain.ProfileModerationRequest
	var status string
	if err := row.Scan(
		&req.ModerationID,
		&req.UserID,
		&req.Username,
		&req.UserAvatarURL,
		&req.FieldName,
		&req.OldValue,
		&req.NewValue,
		&status,
		&req.ModeratedByUserID,
		&req.CreatedAt,
		&req.UpdatedAt,
	); err != nil {
		return domain.ProfileModerationRequest{}, err
	}
	req.Status = domain.ModerationStatus(status)
	return req, nil
}

func profileFieldColumn(fieldName string) (string, error) {
	switch fieldName {
	case "username", "bio", "avatar_url", "banner_url":
		return fieldName, nil
	default:
		return "", fmt.Errorf("invalid profile field name: %w", core_errors.ErrInvalidArgument)
	}
}
