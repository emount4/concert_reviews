package core_ports

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
)

type AdminLogger interface {
	Log(ctx context.Context, entry domain.AdminLog)
}
