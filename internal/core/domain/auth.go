package domain

import "time"

type AuthResponse struct {
	User         User
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
