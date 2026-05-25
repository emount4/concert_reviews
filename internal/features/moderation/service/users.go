package moderation_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ModerationService) GetUsers(
	ctx context.Context,
	search string,
	limit, offset *int,
) ([]domain.User, int, error) {
	if limit != nil && *limit < 0 {
		return nil, 0, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, 0, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	users, total, err := s.moderationRepository.GetUsers(ctx, strings.TrimSpace(search), limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get users: %w", err)
	}

	return users, total, nil
}

func (s *ModerationService) SetUserBan(
	ctx context.Context,
	moderatorID uuid.UUID,
	moderatorRoleID int,
	targetID uuid.UUID,
	isBanned bool,
) (domain.User, error) {
	if err := validateModerator(moderatorID, moderatorRoleID); err != nil {
		return domain.User{}, err
	}
	if targetID == uuid.Nil {
		return domain.User{}, fmt.Errorf("target user id is empty: %w", core_errors.ErrInvalidArgument)
	}
	if moderatorID == targetID {
		return domain.User{}, fmt.Errorf("moderator cannot ban or unban self: %w", core_errors.ErrForbidden)
	}

	target, err := s.moderationRepository.GetUserByID(ctx, targetID)
	if err != nil {
		return domain.User{}, fmt.Errorf("get target user: %w", err)
	}
	if moderatorRoleID == domain.RoleAdminID && target.RoleID != domain.RoleUserID {
		return domain.User{}, fmt.Errorf("admin can ban only regular users: %w", core_errors.ErrForbidden)
	}

	updated, err := s.moderationRepository.SetUserBan(ctx, moderatorID, targetID, isBanned)
	if err != nil {
		return domain.User{}, fmt.Errorf("set user ban: %w", err)
	}

	return updated, nil
}

func (s *ModerationService) SetUserRole(
	ctx context.Context,
	moderatorID uuid.UUID,
	targetID uuid.UUID,
	roleID int,
) (domain.User, error) {
	if moderatorID == uuid.Nil {
		return domain.User{}, fmt.Errorf("moderator id is empty: %w", core_errors.ErrUnauthorized)
	}
	if targetID == uuid.Nil {
		return domain.User{}, fmt.Errorf("target user id is empty: %w", core_errors.ErrInvalidArgument)
	}
	if roleID < domain.RoleUserID || roleID > domain.RoleSuperAdminID {
		return domain.User{}, fmt.Errorf("role_id must be 1, 2 or 3: %w", core_errors.ErrInvalidArgument)
	}

	updated, err := s.moderationRepository.SetUserRole(ctx, moderatorID, targetID, roleID)
	if err != nil {
		return domain.User{}, fmt.Errorf("set user role: %w", err)
	}

	return updated, nil
}

func (s *ModerationService) AnonymizeUser(
	ctx context.Context,
	moderatorID uuid.UUID,
	targetID uuid.UUID,
) (domain.User, error) {
	if moderatorID == uuid.Nil {
		return domain.User{}, fmt.Errorf("moderator id is empty: %w", core_errors.ErrUnauthorized)
	}
	if targetID == uuid.Nil {
		return domain.User{}, fmt.Errorf("target user id is empty: %w", core_errors.ErrInvalidArgument)
	}
	if moderatorID == targetID {
		return domain.User{}, fmt.Errorf("moderator cannot anonymize self: %w", core_errors.ErrForbidden)
	}

	updated, err := s.moderationRepository.AnonymizeUser(ctx, moderatorID, targetID)
	if err != nil {
		return domain.User{}, fmt.Errorf("anonymize user: %w", err)
	}

	if s.statsCache != nil {
		_ = s.statsCache.InvalidateGlobalStats(ctx)
	}

	return updated, nil
}

func validateModerator(moderatorID uuid.UUID, moderatorRoleID int) error {
	if moderatorID == uuid.Nil {
		return fmt.Errorf("moderator id is empty: %w", core_errors.ErrUnauthorized)
	}
	if moderatorRoleID < domain.RoleAdminID {
		return fmt.Errorf("moderator role is not allowed: %w", core_errors.ErrForbidden)
	}
	return nil
}
