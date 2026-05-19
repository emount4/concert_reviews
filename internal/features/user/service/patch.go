package user_service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *UserService) UpdateMe(ctx context.Context, userID uuid.UUID, patch domain.UserPatch) (domain.User, error) {
	if userID == uuid.Nil {
		return domain.User{}, fmt.Errorf("user id is empty: %w", core_errors.ErrUnauthorized)
	}

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("fetch profile for patch: %w", err)
	}

	if !user.IsActive || user.IsBanned {
		return domain.User{}, fmt.Errorf("account is not available: %w", core_errors.ErrForbidden)
	}

	if patch.IsEmpty() {
		return user, nil
	}

	current := user
	normalizeProfilePatch(&patch)
	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("apply profile patch: %w", err)
	}

	if err := s.validateUsername(ctx, current, patch); err != nil {
		return domain.User{}, err
	}

	if err := s.validateMediaKeys(ctx, patch); err != nil {
		return domain.User{}, err
	}

	if err := s.userRepository.SubmitProfilePatch(ctx, current, patch); err != nil {
		return domain.User{}, fmt.Errorf("submit profile moderation: %w", err)
	}

	return s.userRepository.GetByID(ctx, userID)
}

func normalizeProfilePatch(patch *domain.UserPatch) {
	if patch.Bio.Set && patch.Bio.Value != nil && strings.TrimSpace(*patch.Bio.Value) == "" {
		patch.Bio.Value = nil
	}
	if patch.Username.Set && patch.Username.Value != nil {
		trimmed := strings.TrimSpace(*patch.Username.Value)
		patch.Username.Value = &trimmed
	}
	if patch.AvatarKey.Set && patch.AvatarKey.Value != nil {
		trimmed := strings.TrimSpace(*patch.AvatarKey.Value)
		patch.AvatarKey.Value = &trimmed
	}
	if patch.BannerKey.Set && patch.BannerKey.Value != nil {
		trimmed := strings.TrimSpace(*patch.BannerKey.Value)
		patch.BannerKey.Value = &trimmed
	}
}

func (s *UserService) validateUsername(ctx context.Context, current domain.User, patch domain.UserPatch) error {
	if !patch.Username.Set || patch.Username.Value == nil || *patch.Username.Value == current.Username {
		return nil
	}

	existing, err := s.userRepository.GetByUsername(ctx, *patch.Username.Value)
	if err == nil && existing.ID != current.ID {
		return fmt.Errorf("username already exists: %w", core_errors.ErrConflict)
	}
	if err != nil && !errors.Is(err, core_errors.ErrNotFound) {
		return fmt.Errorf("check username availability: %w", err)
	}
	return nil
}

func (s *UserService) validateMediaKeys(ctx context.Context, patch domain.UserPatch) error {
	if patch.AvatarKey.Set && patch.AvatarKey.Value != nil {
		if _, err := s.s3.FileExists(ctx, *patch.AvatarKey.Value); err != nil {
			return fmt.Errorf("avatar file not found in S3: %w", err)
		}
	}
	if patch.BannerKey.Set && patch.BannerKey.Value != nil {
		if _, err := s.s3.FileExists(ctx, *patch.BannerKey.Value); err != nil {
			return fmt.Errorf("banner file not found in S3: %w", err)
		}
	}
	return nil
}
