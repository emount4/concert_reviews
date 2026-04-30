package auth_service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	"github.com/google/uuid"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"go.uber.org/zap"
)

func (s *AuthService) LinkTG(
	ctx context.Context,
	user core_domain.User,
	initData string,
) error {
	if s.authRepository == nil {
		return ErrAuthRepositoryNotConfigured
	}
	if user.ID == uuid.Nil {
		return fmt.Errorf("%w: user_id is required", core_errors.ErrInvalidArgument)
	}
	data, err := s.getData(ctx, initData)
	if err != nil {
		return fmt.Errorf("init data: %w", err)
	}

	err = s.authRepository.LinkTG(ctx, user.ID, data.User.Username, data.User.ID)

	if err != nil {
		return fmt.Errorf("cannot link tg: %w", err)
	}

	return nil
}

func (s *AuthService) LoginTG(
	ctx context.Context,
	initData string,
) (core_domain.AuthResponse, error) {
	if s.authRepository == nil {
		return core_domain.AuthResponse{}, ErrAuthRepositoryNotConfigured
	}
	data, err := s.getData(ctx, initData)
	if err != nil {
		return core_domain.AuthResponse{}, fmt.Errorf("init data: %w", err)
	}

	user, err := s.authRepository.GetUserByTelegramID(ctx, data.User.ID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return core_domain.AuthResponse{}, fmt.Errorf("%w: unknown telegram account", core_errors.ErrUnauthorized)
		}
		return core_domain.AuthResponse{}, fmt.Errorf("get user by telegram id: %w", err)
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

func (s *AuthService) getData(ctx context.Context, initData string) (initdata.InitData, error) {
	log := core_logger.FromContext(ctx)
	if strings.TrimSpace(initData) == "" {
		return initdata.InitData{}, fmt.Errorf("%w: init_data is required", core_errors.ErrInvalidArgument)
	}

	botToken := s.config.BotToken

	expIn := 24 * time.Hour

	err := initdata.Validate(initData, botToken, expIn)

	if err != nil {
		log.Debug("init data validation", zap.Error(err))
		return initdata.InitData{}, fmt.Errorf("%w: invalid telegram credentials: %v", core_errors.ErrUnauthorized, err)
	}

	data, err := initdata.Parse(initData)

	if err != nil {
		log.Debug("init data parse", zap.Error(err))
		return initdata.InitData{}, fmt.Errorf("%w: init data parse: %v", core_errors.ErrInvalidArgument, err)
	}

	return data, nil
}
