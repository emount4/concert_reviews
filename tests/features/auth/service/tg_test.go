package auth_service_test

import (
	"context"
	"testing"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
	"github.com/google/uuid"
)

func TestAuthServiceLinkTGNotConfigured(t *testing.T) {
	service := auth_service.NewAuthService(nil, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})
	if err := service.LinkTG(context.Background(), core_domain.User{}, "init"); err == nil {
		t.Fatal("expected error for missing repository")
	}
}

func TestAuthServiceLinkTGSuccess(t *testing.T) {
	userID := uuid.New()
	user := core_domain.User{ID: userID, Username: "user"}
	repo := &fakeAuthRepository{
		linkTGFn: func(ctx context.Context, id uuid.UUID, username string, tgID int64) error {
			return nil
		},
	}
	service := auth_service.NewAuthService(repo, nil, auth_service.Config{}, &fakeHasher{}, &fakeJWTManager{})
	ctx := newTestContext(t)
	if err := service.LinkTG(ctx, user, "invalid"); err == nil {
		t.Fatal("expected error for invalid init data")
	}
}
