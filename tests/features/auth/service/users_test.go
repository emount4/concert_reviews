package auth_service_test

import (
	"context"
	"strings"
	"testing"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
)

func TestAuthServiceCreateUserTrimsAndLowercases(t *testing.T) {
	var captured core_domain.User
	repo := &fakeAuthRepository{
		createUserFn: func(ctx context.Context, user core_domain.User) (core_domain.User, error) {
			captured = user
			return user, nil
		},
	}
	service := auth_service.NewAuthService(repo, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})

	user := core_domain.User{Username: "  UserName ", Email: "TeSt@Example.com", RoleID: 1}
	created, err := service.CreateUser(context.Background(), user, "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.Username != "UserName" {
		t.Fatalf("expected trimmed username, got %q", captured.Username)
	}
	if captured.Email != "test@example.com" {
		t.Fatalf("expected lowercased email, got %q", captured.Email)
	}
	if !strings.HasPrefix(captured.PasswordHash, "hash:") {
		t.Fatalf("expected password hash to be set")
	}
	if created.ID == (core_domain.User{}).ID {
		t.Fatalf("expected user ID to be set")
	}
}
