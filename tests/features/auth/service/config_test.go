package auth_service_test

import (
	"testing"
	"time"

	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
)

func TestAuthServiceConfigFields(t *testing.T) {
	cfg := auth_service.Config{
		PasswordSalt:    "salt",
		JWTSigningKey:   "secret",
		AccessTokenTTL:  time.Minute,
		RefreshTokenTTL: time.Hour,
		BotToken:        "bot",
	}
	if cfg.PasswordSalt == "" || cfg.JWTSigningKey == "" || cfg.BotToken == "" {
		t.Fatal("expected config fields to be set")
	}
}
