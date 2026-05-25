package moderation_postgres_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

func (r *ModerationRepository) GetAdminLogs(
	ctx context.Context,
	moderatorID *uuid.UUID,
	targetType *string,
	action *string,
	limit, offset *int,
) ([]domain.AdminLog, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	conditions := []string{"1=1"}
	args := []any{}

	if moderatorID != nil {
		args = append(args, *moderatorID)
		conditions = append(conditions, fmt.Sprintf("ml.moderator_user_id = $%d", len(args)))
	}
	if targetType != nil && strings.TrimSpace(*targetType) != "" {
		args = append(args, strings.TrimSpace(*targetType))
		conditions = append(conditions, fmt.Sprintf("ml.target_type = $%d::target_type_enum", len(args)))
	}
	if action != nil && strings.TrimSpace(*action) != "" {
		args = append(args, strings.TrimSpace(*action))
		conditions = append(conditions, fmt.Sprintf("ml.action = $%d", len(args)))
	}

	whereClause := " WHERE " + strings.Join(conditions, " AND ")
	countQuery := "SELECT COUNT(*) FROM moderation_logs ml" + whereClause

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count admin logs: %w", err)
	}

	query := `
		SELECT
			ml.log_id,
			ml.moderator_user_id,
			u.username,
			ml.action,
			COALESCE(ml.target_id, ''),
			COALESCE(ml.target_type::text, ''),
			COALESCE(ml.details, '{}'::jsonb),
			ml.created_at
		FROM moderation_logs ml
		LEFT JOIN users u ON u.user_id = ml.moderator_user_id
	` + whereClause + `
		ORDER BY ml.created_at DESC, ml.log_id DESC
	`
	if limit != nil && offset != nil {
		args = append(args, *limit, *offset)
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query admin logs: %w", err)
	}
	defer rows.Close()

	logs := make([]domain.AdminLog, 0)
	for rows.Next() {
		var item domain.AdminLog
		var moderatorID *uuid.UUID
		var details []byte
		if err := rows.Scan(
			&item.LogID,
			&moderatorID,
			&item.ModeratorUsername,
			&item.Action,
			&item.TargetID,
			&item.TargetType,
			&details,
			&item.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan admin log: %w", err)
		}
		if moderatorID != nil {
			item.ModeratorID = *moderatorID
		}
		if len(details) > 0 {
			if err := json.Unmarshal(details, &item.Details); err != nil {
				return nil, 0, fmt.Errorf("unmarshal admin log details: %w", err)
			}
		}
		if item.Details == nil {
			item.Details = map[string]any{}
		}
		logs = append(logs, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate admin logs rows: %w", err)
	}

	return logs, total, nil
}
