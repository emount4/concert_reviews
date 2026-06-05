package auth_service_test

import (
	"context"
	"testing"

	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
)

func TestAuthServiceLoginTGMissingRepository(t *testing.T) {
	service := auth_service.NewAuthService(nil, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})
	if _, err := service.LoginTG(context.Background(), "init"); err == nil {
		t.Fatal("expected error for missing repository")
	}
}

func TestAuthServiceLoginTGInvalidInitData(t *testing.T) {
	repo := &fakeAuthRepository{}
	service := auth_service.NewAuthService(repo, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})

	ctx := newTestContext(t)
	if _, err := service.LoginTG(ctx, "invalid"); err == nil {
		t.Fatal("expected error for invalid init data")
	}
}
