package user_postgres_repository

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

func (r *UserRepository) GetProfileModerationRequests(
	ctx context.Context,
	userID uuid.UUID,
	status *domain.ModerationStatus,
) ([]domain.ProfileModerationRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT moderation_id, user_id, field_name, old_value, new_value, status,
		       moderated_by_user_id, created_at, updated_at
		FROM profile_moderation
		WHERE user_id = $1
	`
	args := []any{userID}
	if status != nil {
		query += " AND status = $2"
		args = append(args, string(*status))
	}
	query += " ORDER BY created_at DESC, moderation_id DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query profile moderation requests: %w", err)
	}
	defer rows.Close()

	requests := make([]domain.ProfileModerationRequest, 0)
	for rows.Next() {
		var req domain.ProfileModerationRequest
		var status string
		if err := rows.Scan(
			&req.ModerationID,
			&req.UserID,
			&req.FieldName,
			&req.OldValue,
			&req.NewValue,
			&status,
			&req.ModeratedByUserID,
			&req.CreatedAt,
			&req.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan profile moderation request: %w", err)
		}
		req.Status = domain.ModerationStatus(status)
		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate profile moderation requests: %w", err)
	}

	return requests, nil
}
