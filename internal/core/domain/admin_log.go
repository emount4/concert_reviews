package domain

import (
	"github.com/google/uuid"
)

const (
	LogTargetUser    = "user"
	LogTargetReview  = "review"
	LogTargetArtist  = "artist"
	LogTargetVenue   = "venue"
	LogTargetConcert = "concert"
)

type AdminLog struct {
	ModeratorID uuid.UUID
	Action      string
	TargetID    string
	TargetType  string
	Details     map[string]any
}
