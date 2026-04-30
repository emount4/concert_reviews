package auth_service

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
)

type registerRequestValidation struct {
	Username string `validate:"required,min=3,max=50"`
	Email    string `validate:"required,email,min=3,max=255"`
	Password string `validate:"required,min=4,max=50"`
}

var ErrLoginNotImplemented = errors.New("login is not implemented yet")

func (s *AuthService) Register(ctx context.Context, user core_domain.User, password string) (core_domain.AuthResponse, error) {
	if s.authRepository == nil {
		return core_domain.AuthResponse{}, ErrAuthRepositoryNotConfigured
	}
	if s.txManager == nil {
		return core_domain.AuthResponse{}, ErrTxManagerNotConfigured
	}

	if err := s.validate.Struct(registerRequestValidation{
		Username: user.Username,
		Email:    user.Email,
		Password: password,
	}); err != nil {
		return core_domain.AuthResponse{}, fmt.Errorf("validate register request: %w", err)
	}
	//TODO: ПРИВЯЗКА TG, ЕСЛИ ЕСТЬ INIT DATA
	var tokens core_domain.AuthResponse
	if err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		createdUser, err := s.CreateUser(txCtx, user, password)
		if err != nil {
			return err
		}

		generated, err := s.jwt.Generate(createdUser.ID, createdUser.RoleID, s.config.AccessTokenTTL)
		if err != nil {
			return fmt.Errorf("generate tokens: %w", err)
		}
		generated.User = createdUser

		if err := s.authRepository.CreateSession(txCtx, generated); err != nil {
			return fmt.Errorf("create session: %w", err)
		}

		tokens = generated
		return nil
	}); err != nil {
		return core_domain.AuthResponse{}, err
	}

	return tokens, nil
}
