package venue_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/jackc/pgx/v5"
)

func (r *VenueRepository) PatchVenue(
	ctx context.Context,
	id int,
	patch domain.VenuePatch,
) (domain.Venue, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exec := core_postgres_tx.Executor(r.pool)
	if txExec, ok := core_postgres_tx.ExecutorFromContext(ctx); ok {
		exec = txExec
	}

	// 1. Получаем текущую версию для применения патча
	var current VenueRecord
	err := exec.QueryRow(ctx, `
		SELECT venue_id, city_id, name, address, capacity, photo_url, description, social_links, status, created_at, deleted_at
		FROM venues
		WHERE venue_id = $1 AND deleted_at IS NULL
	`, id).Scan(
		&current.VenueID, &current.CityID, &current.Name, &current.Address, &current.Capacity,
		&current.PhotoURL, &current.Description, &current.SocialLinks, &current.Status,
		&current.CreatedAt, &current.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Venue{}, core_errors.ErrNotFound
		}
		return domain.Venue{}, fmt.Errorf("get venue for patch: %w", err)
	}

	// 2. Применяем патч к доменной модели (валидация внутри)
	venue := current.MapToDomain()
	if err := venue.ApplyPatch(patch); err != nil {
		return domain.Venue{}, fmt.Errorf("apply patch: %w", err)
	}

	// 3. Формируем динамический UPDATE-запрос (только изменённые поля)
	updates := make([]string, 0)
	args := make([]any, 0)
	argIdx := 1

	if patch.CityID.Set && patch.CityID.Value != nil {
		updates = append(updates, fmt.Sprintf("city_id = $%d", argIdx))
		args = append(args, venue.CityID)
		argIdx++
	}
	if patch.Name.Set && patch.Name.Value != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, venue.Name)
		argIdx++
	}
	if patch.Address.Set {
		updates = append(updates, fmt.Sprintf("address = $%d", argIdx))
		args = append(args, venue.Address)
		argIdx++
	}
	if patch.Capacity.Set {
		updates = append(updates, fmt.Sprintf("capacity = $%d", argIdx))
		args = append(args, venue.Capacity)
		argIdx++
	}
	if patch.PhotoKey.Set {
		updates = append(updates, fmt.Sprintf("photo_url = $%d", argIdx))
		args = append(args, venue.PhotoURL)
		argIdx++
	}
	if patch.Description.Set {
		updates = append(updates, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, venue.Description)
		argIdx++
	}
	if patch.SocialLinks.Set {
		updates = append(updates, fmt.Sprintf("social_links = $%d", argIdx))
		if patch.SocialLinks.Value != nil {
			args = append(args, *patch.SocialLinks.Value) // pgx: map → JSONB
		} else {
			args = append(args, nil) // явный NULL
		}
		argIdx++
	}
	if patch.Status.Set && patch.Status.Value != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, string(venue.Status))
		argIdx++
	}

	// Если ничего не меняем — возвращаем текущую версию
	if len(updates) == 0 {
		return venue, nil
	}

	// 4. Выполняем UPDATE с RETURNING
	query := fmt.Sprintf(`
		UPDATE venues
		SET %s
		WHERE venue_id = $%d AND deleted_at IS NULL
		RETURNING venue_id, city_id, name, address, capacity, photo_url, description, social_links, status, created_at, deleted_at
	`, strings.Join(updates, ", "), argIdx)

	args = append(args, id)

	var updated VenueRecord
	err = exec.QueryRow(ctx, query, args...).Scan(
		&updated.VenueID, &updated.CityID, &updated.Name, &updated.Address, &updated.Capacity,
		&updated.PhotoURL, &updated.Description, &updated.SocialLinks, &updated.Status,
		&updated.CreatedAt, &updated.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Venue{}, core_errors.ErrNotFound
		}
		return domain.Venue{}, fmt.Errorf("execute patch: %w", err)
	}

	return updated.MapToDomain(), nil
}
