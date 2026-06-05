package auth_service_test

import (
	"context"
	"strings"
	"testing"
	"time"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
	"github.com/google/uuid"
)

func TestAuthServiceRegisterMissingRepository(t *testing.T) {
	service := auth_service.NewAuthService(nil, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})
	_, err := service.Register(context.Background(), core_domain.User{}, "pass")
	if err != auth_service.ErrAuthRepositoryNotConfigured {
		t.Fatalf("expected ErrAuthRepositoryNotConfigured, got %v", err)
	}
}

func TestAuthServiceRegisterMissingTxManager(t *testing.T) {
	service := auth_service.NewAuthService(&fakeAuthRepository{}, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})
	_, err := service.Register(context.Background(), core_domain.User{}, "pass")
	if err != auth_service.ErrTxManagerNotConfigured {
		t.Fatalf("expected ErrTxManagerNotConfigured, got %v", err)
	}
}

func TestAuthServiceRegisterValidationError(t *testing.T) {
	repo := &fakeAuthRepository{}
	tx := &fakeTxManager{}
	service := auth_service.NewAuthService(repo, tx, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})

	user := core_domain.User{Username: "ab", Email: "bad"}
	_, err := service.Register(context.Background(), user, "123")
	if err == nil || !strings.Contains(err.Error(), "validate register request") {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestAuthServiceRegisterSuccess(t *testing.T) {
	repo := &fakeAuthRepository{}
	tx := &fakeTxManager{}
	cache := &fakeStatsCache{}

	repo.createUserFn = func(ctx context.Context, user core_domain.User) (core_domain.User, error) {
		return user, nil
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
	service := auth_service.NewAuthService(repo, tx, auth_service.Config{}, &fakeHasher{}, jwt, cache)

	user := core_domain.User{Username: "  UserName ", Email: "TeSt@Example.com", RoleID: 1}
	resp, err := service.Register(context.Background(), user, "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "access" || resp.RefreshToken != "refresh" {
		t.Fatalf("unexpected tokens: %+v", resp)
	}
	if repo.createSessionCalls != 1 {
		t.Fatalf("expected CreateSession to be called once, got %d", repo.createSessionCalls)
	}
	if !tx.called {
		t.Fatal("expected transaction to be used")
	}
	if !cache.called {
		t.Fatal("expected stats cache invalidation")
	}
}
