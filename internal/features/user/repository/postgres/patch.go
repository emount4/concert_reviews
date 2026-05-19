package user_postgres_repository

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
)

func (r *UserRepository) SubmitProfilePatch(ctx context.Context, current domain.User, patch domain.UserPatch) error {
	return r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		if patch.Username.Set {
			if err := insertProfileModeration(txCtx, exec, current, "username", &current.Username, patch.Username.Value); err != nil {
				return err
			}
		}
		if patch.Bio.Set {
			if err := insertProfileModeration(txCtx, exec, current, "bio", current.Bio, patch.Bio.Value); err != nil {
				return err
			}
		}
		if patch.AvatarKey.Set {
			if err := insertProfileModeration(txCtx, exec, current, "avatar_url", current.AvatarURL, patch.AvatarKey.Value); err != nil {
				return err
			}
		}
		if patch.BannerKey.Set {
			if err := insertProfileModeration(txCtx, exec, current, "banner_url", current.BannerURL, patch.BannerKey.Value); err != nil {
				return err
			}
		}

		return nil
	})
}

func insertProfileModeration(
	ctx context.Context,
	exec core_postgres_tx.Executor,
	user domain.User,
	fieldName string,
	oldValue *string,
	newValue *string,
) error {
	if equalStringPtr(oldValue, newValue) {
		return nil
	}

	query := `
		INSERT INTO profile_moderation (user_id, field_name, old_value, new_value, status)
		VALUES ($1, $2, $3, $4, 'pending')
	`
	if _, err := exec.Exec(ctx, query, user.ID, fieldName, oldValue, newValue); err != nil {
		return fmt.Errorf("insert profile moderation for %s: %w", fieldName, err)
	}
	return nil
}

func equalStringPtr(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
