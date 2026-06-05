package auth_service_test

import (
	"context"
	"testing"

	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
)

func TestAuthServiceLogoutMissingRepository(t *testing.T) {
	service := auth_service.NewAuthService(nil, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})
	if err := service.Logout(context.Background(), "token"); err != auth_service.ErrAuthRepositoryNotConfigured {
		t.Fatalf("expected ErrAuthRepositoryNotConfigured, got %v", err)
	}
}
