package auth_service_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
	"github.com/google/uuid"
)

func TestAuthServiceLoginValidationError(t *testing.T) {
	repo := &fakeAuthRepository{}
	service := auth_service.NewAuthService(repo, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})

	_, err := service.Login(context.Background(), "bad-email", "123")
	if err == nil || !strings.Contains(err.Error(), "validate login request") {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestAuthServiceLoginUserNotFound(t *testing.T) {
	repo := &fakeAuthRepository{
		getUserByEmailFn: func(ctx context.Context, email string) (core_domain.User, error) {
			return core_domain.User{}, core_errors.ErrNotFound
		},
	}
	service := auth_service.NewAuthService(repo, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})

	_, err := service.Login(context.Background(), "user@example.com", "pass")
	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("expected unauthorized, got %v", err)
	}
}

func TestAuthServiceLoginInvalidPassword(t *testing.T) {
	repo := &fakeAuthRepository{
		getUserByEmailFn: func(ctx context.Context, email string) (core_domain.User, error) {
			return core_domain.User{
				ID:           uuid.New(),
				Email:        email,
				PasswordHash: "hash:other",
				Username:     "user",
				RoleID:       1,
				IsActive:     true,
				IsBanned:     false,
			}, nil
		},
	}
	service := auth_service.NewAuthService(repo, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})

	_, err := service.Login(context.Background(), "user@example.com", "pass")
	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("expected unauthorized, got %v", err)
	}
}

func TestAuthServiceLoginInactiveUser(t *testing.T) {
	repo := &fakeAuthRepository{
		getUserByEmailFn: func(ctx context.Context, email string) (core_domain.User, error) {
			return core_domain.User{
				ID:           uuid.New(),
				Email:        email,
				PasswordHash: "hash:pass",
				Username:     "user",
				RoleID:       1,
				IsActive:     false,
				IsBanned:     false,
			}, nil
		},
	}
	service := auth_service.NewAuthService(repo, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})

	_, err := service.Login(context.Background(), "user@example.com", "pass")
	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("expected unauthorized, got %v", err)
	}
}

func TestAuthServiceLoginSuccessCreatesSession(t *testing.T) {
	userID := uuid.New()
	repo := &fakeAuthRepository{
		getUserByEmailFn: func(ctx context.Context, email string) (core_domain.User, error) {
			return core_domain.User{
				ID:           userID,
				Email:        email,
				PasswordHash: "hash:pass",
				Username:     "user",
				RoleID:       2,
				IsActive:     true,
				IsBanned:     false,
			}, nil
		},
	}
	jwt := &fakeJWTManager{
		generateFn: func(id uuid.UUID, roleID int, accessTTL, refreshTTL time.Duration) (core_domain.AuthResponse, error) {
			return core_domain.AuthResponse{
				AccessToken:  "access",
				RefreshToken: "refresh",
				ExpiresAt:    time.Now().Add(time.Hour),
			}, nil
		},
	}
	service := auth_service.NewAuthService(repo, nil, auth_service.Config{}, &fakeHasher{}, jwt)

	resp, err := service.Login(context.Background(), "user@example.com", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "access" || resp.RefreshToken != "refresh" {
		t.Fatalf("unexpected tokens: %+v", resp)
	}
	if resp.User.ID != userID {
		t.Fatalf("expected user to be set in response")
	}
	if repo.createSessionCalls != 1 {
		t.Fatalf("expected CreateSession to be called once, got %d", repo.createSessionCalls)
	}
	if repo.lastCreateSession.AccessToken != "access" {
		t.Fatalf("expected CreateSession to receive tokens")
	}
}
