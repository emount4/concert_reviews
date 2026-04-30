package auth_service

import (
	"context"
	"fmt"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	"go.uber.org/zap"
)

func (s *AuthService) Refresh(
	ctx context.Context,
	oldToken string,
) (core_domain.AuthResponse, error) {
	if s.authRepository == nil {
		return core_domain.AuthResponse{}, ErrAuthRepositoryNotConfigured
	}
	if s.txManager == nil {
		return core_domain.AuthResponse{}, ErrTxManagerNotConfigured
	}

	session, err := s.authRepository.GetSession(ctx, oldToken)
	if err != nil {
		return core_domain.AuthResponse{}, core_errors.ErrUnauthorized
	}

	user, err := s.authRepository.GetUserByID(ctx, session.UserID)

	if err != nil || !user.IsActive || user.IsBanned {
		err := s.authRepository.DeleteAllUserSessions(ctx, session.UserID)

		if err != nil {
			log := core_logger.FromContext(ctx)
			log.Error(
				"cannot delete all sessions",
				zap.Error(err),
			)
		}
		return core_domain.AuthResponse{}, core_errors.ErrUnauthorized
	}

	newTokens, err := s.jwt.Generate(
		user.ID,
		user.RoleID,
		s.config.AccessTokenTTL,
	)
	if err != nil {
		return core_domain.AuthResponse{}, core_errors.ErrUnauthorized
	}
	newTokens.User = user

	err = s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		if err := s.authRepository.DeleteSession(txCtx, oldToken); err != nil {
			return fmt.Errorf("delete old session: %w", err)
		}
		if err := s.authRepository.CreateSession(txCtx, newTokens); err != nil {
			return fmt.Errorf("create session: %w", err)
		}
		return nil
	})
	if err != nil {
		return core_domain.AuthResponse{}, core_errors.ErrUnauthorized
	}

	return newTokens, nil
}
