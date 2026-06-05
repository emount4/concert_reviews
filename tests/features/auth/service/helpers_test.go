package auth_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
	"github.com/google/uuid"
)

type fakeAuthRepository struct {
	createUserFn            func(ctx context.Context, user core_domain.User) (core_domain.User, error)
	getUserByEmailFn        func(ctx context.Context, email string) (core_domain.User, error)
	getUserByTelegramIDFn   func(ctx context.Context, tgID int64) (core_domain.User, error)
	getUserByIDFn           func(ctx context.Context, id uuid.UUID) (core_domain.User, error)
	createSessionFn         func(ctx context.Context, response core_domain.AuthResponse) error
	getSessionFn            func(ctx context.Context, oldToken string) (core_domain.RefreshToken, error)
	deleteSessionFn         func(ctx context.Context, token string) error
	deleteAllUserSessionsFn func(ctx context.Context, userID uuid.UUID) error
	linkTGFn                func(ctx context.Context, userID uuid.UUID, username string, tgID int64) error

	lastCreateSession  core_domain.AuthResponse
	createSessionCalls int
	deleteSessionCalls int
	deleteAllCalls     int
}

func (f *fakeAuthRepository) CreateUser(ctx context.Context, user core_domain.User) (core_domain.User, error) {
	if f.createUserFn == nil {
		return core_domain.User{}, errors.New("CreateUser not stubbed")
	}
	return f.createUserFn(ctx, user)
}

func (f *fakeAuthRepository) GetUserByEmail(ctx context.Context, email string) (core_domain.User, error) {
	if f.getUserByEmailFn == nil {
		return core_domain.User{}, errors.New("GetUserByEmail not stubbed")
	}
	return f.getUserByEmailFn(ctx, email)
}

func (f *fakeAuthRepository) GetUserByTelegramID(ctx context.Context, tgID int64) (core_domain.User, error) {
	if f.getUserByTelegramIDFn == nil {
		return core_domain.User{}, errors.New("GetUserByTelegramID not stubbed")
	}
	return f.getUserByTelegramIDFn(ctx, tgID)
}

func (f *fakeAuthRepository) GetUserByID(ctx context.Context, id uuid.UUID) (core_domain.User, error) {
	if f.getUserByIDFn == nil {
		return core_domain.User{}, errors.New("GetUserByID not stubbed")
	}
	return f.getUserByIDFn(ctx, id)
}

func (f *fakeAuthRepository) CreateSession(ctx context.Context, response core_domain.AuthResponse) error {
	f.createSessionCalls++
	f.lastCreateSession = response
	if f.createSessionFn == nil {
		return nil
	}
	return f.createSessionFn(ctx, response)
}

func (f *fakeAuthRepository) GetSession(ctx context.Context, oldToken string) (core_domain.RefreshToken, error) {
	if f.getSessionFn == nil {
		return core_domain.RefreshToken{}, errors.New("GetSession not stubbed")
	}
	return f.getSessionFn(ctx, oldToken)
}

func (f *fakeAuthRepository) DeleteSession(ctx context.Context, token string) error {
	f.deleteSessionCalls++
	if f.deleteSessionFn == nil {
		return nil
	}
	return f.deleteSessionFn(ctx, token)
}

func (f *fakeAuthRepository) DeleteAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	f.deleteAllCalls++
	if f.deleteAllUserSessionsFn == nil {
		return nil
	}
	return f.deleteAllUserSessionsFn(ctx, userID)
}

func (f *fakeAuthRepository) LinkTG(ctx context.Context, userID uuid.UUID, username string, tgID int64) error {
	if f.linkTGFn == nil {
		return nil
	}
	return f.linkTGFn(ctx, userID, username, tgID)
}

type fakeTxManager struct {
	withinTxFn func(ctx context.Context, fn func(ctx context.Context) error) error
	called     bool
}

func (f *fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	f.called = true
	if f.withinTxFn != nil {
		return f.withinTxFn(ctx, fn)
	}
	return fn(ctx)
}

type fakeStatsCache struct {
	called bool
	err    error
}

func (f *fakeStatsCache) InvalidateGlobalStats(ctx context.Context) error {
	f.called = true
	return f.err
}

type fakeHasher struct{}

func (h *fakeHasher) Hash(password string) string {
	return "hash:" + password
}

func (h *fakeHasher) ValidatePassword(password, hash string) bool {
	return hash == h.Hash(password)
}

type fakeJWTManager struct {
	generateFn func(userID uuid.UUID, roleID int, accessTTL time.Duration, refreshTTL time.Duration) (core_domain.AuthResponse, error)
}

func (f *fakeJWTManager) NewAccessToken(userID uuid.UUID, roleID int, ttl time.Duration) (string, error) {
	return "", nil
}

func (f *fakeJWTManager) Parse(accessToken string) (*auth_service.JWTClaims, error) {
	return nil, nil
}

func (f *fakeJWTManager) NewRefreshToken() (string, error) {
	return "", nil
}

func (f *fakeJWTManager) Generate(userID uuid.UUID, roleID int, accessTTL time.Duration, refreshTTL time.Duration) (core_domain.AuthResponse, error) {
	if f.generateFn == nil {
		return core_domain.AuthResponse{}, errors.New("Generate not stubbed")
	}
	return f.generateFn(userID, roleID, accessTTL, refreshTTL)
}

func newTestContext(t *testing.T) context.Context {
	t.Helper()
	log, err := core_logger.NewLogger(core_logger.LoggerConfig{
		LogLevel:  "DEBUG",
		LogFolder: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	t.Cleanup(func() {
		log.Close()
	})

	return context.WithValue(context.Background(), "log", log)
}
