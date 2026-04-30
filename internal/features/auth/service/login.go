package auth_service

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

type loginRequestValidation struct {
	Email    string `validate:"required,email,min=3,max=100"`
	Password string `validate:"required,min=4,max=50"`
}

func (s *AuthService) Login(ctx context.Context, email, password string) (core_domain.AuthResponse, error) {
	if s.authRepository == nil {
		return core_domain.AuthResponse{}, ErrAuthRepositoryNotConfigured
	}

	if err := s.validate.Struct(loginRequestValidation{
		Email:    email,
		Password: password,
	}); err != nil {
		return core_domain.AuthResponse{}, fmt.Errorf("validate login request: %w", err)
	}

	user, err := s.authRepository.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return core_domain.AuthResponse{}, fmt.Errorf("%w: invalid credentials", core_errors.ErrUnauthorized)
		}
		return core_domain.AuthResponse{}, fmt.Errorf("get user by email: %w", err)
	}

	if !s.hasher.ValidatePassword(password, user.PasswordHash) {
		return core_domain.AuthResponse{}, fmt.Errorf("%w: invalid credentials", core_errors.ErrUnauthorized)
	}

	tokens, err := s.jwt.Generate(user.ID, user.RoleID, s.config.AccessTokenTTL)
	if err != nil {
		return core_domain.AuthResponse{}, fmt.Errorf("generate tokens: %w", err)
	}
	tokens.User = user

	if err := s.authRepository.CreateSession(ctx, tokens); err != nil {
		return core_domain.AuthResponse{}, fmt.Errorf("create session: %w", err)
	}
	return tokens, nil
}
