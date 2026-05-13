package concert_postgres_repository

import (
	"errors"
	"fmt"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *ConcertRepository) handleDBError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503": // Foreign Key Violation
			return fmt.Errorf("related entity not found: %w", core_errors.ErrNotFound)
		case "23505": // Unique Violation
			return fmt.Errorf("conflict occurred: %w", core_errors.ErrConflict)
		}
	}
	return err
}
