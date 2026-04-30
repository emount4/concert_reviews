package auth_service

import (
	"context"
	"fmt"
	"strings"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

func (s *AuthService) CreateUser(
	ctx context.Context,
	user core_domain.User,
	password string,
) (core_domain.User, error) {
	if s.authRepository == nil {
		return core_domain.User{}, ErrAuthRepositoryNotConfigured
	}

	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.ID = uuid.New()
	user.PasswordHash = s.hasher.Hash(password)

	if err := user.Validate(); err != nil {
		return core_domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err := s.authRepository.CreateUser(ctx, user)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
