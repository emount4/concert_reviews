package core_postgres_repository

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_postgres_pool "github.com/emount4/concert_reviews/internal/core/repository/postgres/pool"
	"go.uber.org/zap"
)

type AuditRepository struct {
	pool core_postgres_pool.Pool
	log  *core_logger.Logger
}

func NewAuditRepository(pool core_postgres_pool.Pool, log *core_logger.Logger) *AuditRepository {
	return &AuditRepository{
		pool: pool,
		log:  log,
	}
}

func (r *AuditRepository) Log(ctx context.Context, l domain.AdminLog) {
	query := `
		INSERT INTO moderation_logs (moderator_user_id, action, target_id, target_type, details)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, l.ModeratorID, l.Action, l.TargetID, l.TargetType, l.Details)
	if err != nil {
		r.log.Error("failed to write admin audit log", zap.Error(err))
	}
}
