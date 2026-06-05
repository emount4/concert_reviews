package auth_service_test

import (
	"testing"
	"time"

	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
	"github.com/google/uuid"
)

func TestJWTManagerGenerateAndParse(t *testing.T) {
	manager := auth_service.NewManager("secret")
	userID := uuid.New()
	roleID := 2

	now := time.Now()
	resp, err := manager.Generate(userID, roleID, time.Minute, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatalf("expected tokens to be set")
	}
	if !resp.ExpiresAt.After(now) {
		t.Fatalf("expected refresh expiry in the future")
	}

	claims, err := manager.Parse(resp.AccessToken)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if claims.UserID != userID || claims.RoleID != roleID {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestJWTManagerRefreshTokenLength(t *testing.T) {
	manager := auth_service.NewManager("secret")
	token, err := manager.NewRefreshToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(token) != 64 {
		t.Fatalf("expected refresh token length 64, got %d", len(token))
	}
}
