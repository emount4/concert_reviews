package auth_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
	"github.com/google/uuid"
)

func TestAuthServiceRefreshInvalidSession(t *testing.T) {
	repo := &fakeAuthRepository{
		getSessionFn: func(ctx context.Context, oldToken string) (core_domain.RefreshToken, error) {
			return core_domain.RefreshToken{}, errors.New("not found")
		},
	}
	tx := &fakeTxManager{}
	service := auth_service.NewAuthService(repo, tx, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})

	_, err := service.Refresh(context.Background(), "old")
	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("expected unauthorized, got %v", err)
	}
}

func TestAuthServiceRefreshInactiveUserDeletesSessions(t *testing.T) {
	userID := uuid.New()
	repo := &fakeAuthRepository{
		getSessionFn: func(ctx context.Context, oldToken string) (core_domain.RefreshToken, error) {
			return core_domain.RefreshToken{UserID: userID}, nil
		},
		getUserByIDFn: func(ctx context.Context, id uuid.UUID) (core_domain.User, error) {
			return core_domain.User{ID: id, IsActive: false}, nil
		},
	}
	tx := &fakeTxManager{}
	service := auth_service.NewAuthService(repo, tx, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})

	_, err := service.Refresh(context.Background(), "old")
	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("expected unauthorized, got %v", err)
	}
	if repo.deleteAllCalls != 1 {
		t.Fatalf("expected DeleteAllUserSessions, got %d", repo.deleteAllCalls)
	}
}

func TestAuthServiceRefreshSuccess(t *testing.T) {
	userID := uuid.New()
	repo := &fakeAuthRepository{
		getSessionFn: func(ctx context.Context, oldToken string) (core_domain.RefreshToken, error) {
			return core_domain.RefreshToken{UserID: userID}, nil
		},
		getUserByIDFn: func(ctx context.Context, id uuid.UUID) (core_domain.User, error) {
			return core_domain.User{ID: id, RoleID: 1, IsActive: true, IsBanned: false}, nil
		},
	}
	tx := &fakeTxManager{}
	jwt := &fakeJWTManager{
		generateFn: func(id uuid.UUID, roleID int, accessTTL, refreshTTL time.Duration) (core_domain.AuthResponse, error) {
			return core_domain.AuthResponse{
				AccessToken:  "access",
				RefreshToken: "refresh",
				ExpiresAt:    time.Now().Add(time.Hour),
			}, nil
		},
	}
	service := auth_service.NewAuthService(repo, tx, auth_service.Config{}, &fakeHasher{}, jwt)

	resp, err := service.Refresh(context.Background(), "old")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "access" || resp.RefreshToken != "refresh" {
		t.Fatalf("unexpected tokens: %+v", resp)
	}
	if repo.deleteSessionCalls != 1 {
		t.Fatalf("expected DeleteSession to be called once, got %d", repo.deleteSessionCalls)
	}
	if repo.createSessionCalls != 1 {
		t.Fatalf("expected CreateSession to be called once, got %d", repo.createSessionCalls)
	}
}

func TestAuthServiceLogout(t *testing.T) {
	repo := &fakeAuthRepository{}
	service := auth_service.NewAuthService(repo, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})

	if err := service.Logout(context.Background(), "token"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.deleteSessionCalls != 1 {
		t.Fatalf("expected DeleteSession to be called once, got %d", repo.deleteSessionCalls)
	}
}
