package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReviewLike struct {
	LikeID    int
	UserID    uuid.UUID
	ReviewID  uuid.UUID
	CreatedAt time.Time
}
