package core_postgres_tx

import (
	"context"
	"fmt"

	core_postgres_pool "github.com/emount4/concert_reviews/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Executor interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type contextKey struct{}

func ContextWithExecutor(ctx context.Context, exec Executor) context.Context {
	return context.WithValue(ctx, contextKey{}, exec)
}

func ExecutorFromContext(ctx context.Context) (Executor, bool) {
	exec, ok := ctx.Value(contextKey{}).(Executor)
	return exec, ok
}

type Manager struct {
	pool core_postgres_pool.Pool
}

func NewManager(pool core_postgres_pool.Pool) *Manager {
	return &Manager{pool: pool}
}

func (m *Manager) WithinTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	txCtx := ContextWithExecutor(ctx, tx)
	if err := fn(txCtx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != pgx.ErrTxClosed {
			return fmt.Errorf("rollback tx: %v: %w", rollbackErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
