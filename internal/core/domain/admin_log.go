package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	LogTargetUser    = "user"
	LogTargetReview  = "review"
	LogTargetArtist  = "artist"
	LogTargetVenue   = "venue"
	LogTargetConcert = "concert"
	LogTargetCity    = "city"
)

type AdminLog struct {
	LogID             int
	ModeratorID       uuid.UUID
	ModeratorUsername *string
	Action            string
	TargetID          string
	TargetType        string
	Details           map[string]any
	CreatedAt         time.Time
}
