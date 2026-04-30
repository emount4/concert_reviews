package auth_service

import (
	"context"
	"errors"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var (
	ErrAuthRepositoryNotConfigured = errors.New("auth repository is not configured")
	ErrTxManagerNotConfigured      = errors.New("transaction manager is not configured")
)

type AuthService struct {
	authRepository AuthRepository
	txManager      TxManager

	validate *validator.Validate
	hasher   PasswordHasher
	jwt      JWTManager
	config   Config
}

type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type AuthRepository interface {
	CreateUser(ctx context.Context, user core_domain.User) (core_domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (core_domain.User, error)
	GetUserByTelegramID(ctx context.Context, tgID int64) (core_domain.User, error)
	CreateSession(ctx context.Context, response core_domain.AuthResponse) error

	LinkTG(ctx context.Context, userID uuid.UUID, username string, tgID int64) error
}

func NewAuthService(
	authRepository AuthRepository,
	txManager TxManager,
	config Config,
	hasher PasswordHasher,
	jwt JWTManager,
) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		txManager:      txManager,
		validate:       validator.New(),
		hasher:         hasher,
		jwt:            jwt,
		config:         config,
	}
}

func NewService(config Config, hasher PasswordHasher, jwt JWTManager) *AuthService {
	return NewAuthService(nil, nil, config, hasher, jwt)
}
