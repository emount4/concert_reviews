package favorites_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *FavoritesRepository) AddFavorite(
	ctx context.Context,
	userID uuid.UUID,
	target domain.FavoriteTarget,
) (domain.Favorite, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var favorite domain.Favorite
	err := r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		query := `
			INSERT INTO favorites (user_id, target_id, target_type)
			VALUES ($1, $2, $3::target_type_enum)
			RETURNING favorite_id, user_id, target_type::text, target_id, created_at
		`

		var rec favoriteRecord
		if err := exec.QueryRow(txCtx, query, userID, target.ID, string(target.Type)).Scan(
			&rec.FavoriteID,
			&rec.UserID,
			&rec.TargetType,
			&rec.TargetID,
			&rec.CreatedAt,
		); err != nil {
			return mapFavoriteWriteError(err)
		}

		if err := updateFavoriteStats(txCtx, exec, target, 1); err != nil {
			return err
		}

		rec.TargetName = target.Name
		rec.ImageURL = target.ImageURL
		favorite = rec.MapToDomain()
		return nil
	})
	if err != nil {
		return domain.Favorite{}, err
	}

	return favorite, nil
}

func (r *FavoritesRepository) DeleteFavorite(
	ctx context.Context,
	userID uuid.UUID,
	target domain.FavoriteTarget,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	return r.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		exec, _ := core_postgres_tx.ExecutorFromContext(txCtx)

		tag, err := exec.Exec(
			txCtx,
			`DELETE FROM favorites WHERE user_id = $1 AND target_id = $2 AND target_type = $3::target_type_enum`,
			userID,
			target.ID,
			string(target.Type),
		)
		if err != nil {
			return fmt.Errorf("delete favorite: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("favorite not found: %w", core_errors.ErrNotFound)
		}

		if err := updateFavoriteStats(txCtx, exec, target, -1); err != nil {
			return err
		}

		return nil
	})
}

func (r *FavoritesRepository) CountFavoritesByType(
	ctx context.Context,
	userID uuid.UUID,
	targetType domain.FavoriteTargetType,
) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var count int
	err := r.pool.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM favorites WHERE user_id = $1 AND target_type = $2::target_type_enum`,
		userID,
		string(targetType),
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count favorites: %w", err)
	}

	return count, nil
}

func (r *FavoritesRepository) GetFavoriteTarget(
	ctx context.Context,
	targetType domain.FavoriteTargetType,
	targetID string,
) (domain.FavoriteTarget, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	target := domain.FavoriteTarget{
		Type: targetType,
		ID:   targetID,
	}

	var err error
	switch targetType {
	case domain.FavoriteTargetArtist:
		parsedID, parseErr := strconv.Atoi(targetID)
		if parseErr != nil {
			return domain.FavoriteTarget{}, fmt.Errorf("artist id must be integer: %w", core_errors.ErrInvalidArgument)
		}
		err = r.pool.QueryRow(
			ctx,
			`SELECT name, photo_url FROM artists WHERE artist_id = $1 AND deleted_at IS NULL AND status = 'active'`,
			parsedID,
		).Scan(&target.Name, &target.ImageURL)

	case domain.FavoriteTargetVenue:
		parsedID, parseErr := strconv.Atoi(targetID)
		if parseErr != nil {
			return domain.FavoriteTarget{}, fmt.Errorf("venue id must be integer: %w", core_errors.ErrInvalidArgument)
		}
		err = r.pool.QueryRow(
			ctx,
			`SELECT name, photo_url FROM venues WHERE venue_id = $1 AND deleted_at IS NULL AND status = 'active'`,
			parsedID,
		).Scan(&target.Name, &target.ImageURL)

	case domain.FavoriteTargetConcert:
		parsedID, parseErr := uuid.Parse(targetID)
		if parseErr != nil {
			return domain.FavoriteTarget{}, fmt.Errorf("concert id must be uuid: %w", core_errors.ErrInvalidArgument)
		}
		err = r.pool.QueryRow(
			ctx,
			`SELECT title, poster_url FROM concerts WHERE concert_id = $1 AND deleted_at IS NULL`,
			parsedID,
		).Scan(&target.Name, &target.ImageURL)

	default:
		return domain.FavoriteTarget{}, fmt.Errorf("unsupported favorite target type %s: %w", targetType, core_errors.ErrInvalidArgument)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.FavoriteTarget{}, fmt.Errorf("favorite target not found: %w", core_errors.ErrNotFound)
		}
		return domain.FavoriteTarget{}, fmt.Errorf("query favorite target: %w", err)
	}

	return target, nil
}

func (r *FavoritesRepository) GetFavoritesByUsername(
	ctx context.Context,
	username string,
	targetType *domain.FavoriteTargetType,
) ([]domain.Favorite, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var userExists bool
	if err := r.pool.QueryRow(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND is_active = true AND is_banned = false)`,
		username,
	).Scan(&userExists); err != nil {
		return nil, fmt.Errorf("check user by username: %w", err)
	}
	if !userExists {
		return nil, fmt.Errorf("user not found: %w", core_errors.ErrNotFound)
	}

	args := []any{username}
	typeFilter := ""
	if targetType != nil {
		args = append(args, string(*targetType))
		typeFilter = " AND f.target_type = $2::target_type_enum"
	}

	query := `
		SELECT
			f.favorite_id,
			f.user_id,
			f.target_type::text,
			f.target_id,
			CASE
				WHEN f.target_type = 'artist' THEN a.name
				WHEN f.target_type = 'venue' THEN v.name
				WHEN f.target_type = 'concert' THEN c.title
			END AS target_name,
			CASE
				WHEN f.target_type = 'artist' THEN a.photo_url
				WHEN f.target_type = 'venue' THEN v.photo_url
				WHEN f.target_type = 'concert' THEN c.poster_url
			END AS image_url,
			f.created_at
		FROM favorites f
		JOIN users u ON u.user_id = f.user_id
		LEFT JOIN artists a ON f.target_type = 'artist' AND f.target_id = a.artist_id::text AND a.deleted_at IS NULL AND a.status = 'active'
		LEFT JOIN venues v ON f.target_type = 'venue' AND f.target_id = v.venue_id::text AND v.deleted_at IS NULL AND v.status = 'active'
		LEFT JOIN concerts c ON f.target_type = 'concert' AND f.target_id = c.concert_id::text AND c.deleted_at IS NULL
		WHERE u.username = $1
			AND u.is_active = true
			AND u.is_banned = false
			AND (
				(f.target_type = 'artist' AND a.artist_id IS NOT NULL)
				OR (f.target_type = 'venue' AND v.venue_id IS NOT NULL)
				OR (f.target_type = 'concert' AND c.concert_id IS NOT NULL)
			)
	` + typeFilter + `
		ORDER BY f.target_type ASC, f.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query favorites by username: %w", err)
	}
	defer rows.Close()

	favorites := make([]domain.Favorite, 0)
	for rows.Next() {
		var rec favoriteRecord
		if err := rows.Scan(
			&rec.FavoriteID,
			&rec.UserID,
			&rec.TargetType,
			&rec.TargetID,
			&rec.TargetName,
			&rec.ImageURL,
			&rec.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan favorite: %w", err)
		}

		favorites = append(favorites, rec.MapToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate favorites rows: %w", err)
	}

	return favorites, nil
}

func updateFavoriteStats(
	ctx context.Context,
	exec core_postgres_tx.Executor,
	target domain.FavoriteTarget,
	delta int,
) error {
	switch target.Type {
	case domain.FavoriteTargetArtist:
		artistID, _ := strconv.Atoi(target.ID)
		if delta > 0 {
			_, err := exec.Exec(ctx, `
				INSERT INTO artist_stats (artist_id, favorites_count)
				VALUES ($1, 1)
				ON CONFLICT (artist_id) DO UPDATE
				SET favorites_count = COALESCE(artist_stats.favorites_count, 0) + 1
			`, artistID)
			if err != nil {
				return fmt.Errorf("increment artist favorites count: %w", err)
			}
			return nil
		}
		_, err := exec.Exec(ctx, `UPDATE artist_stats SET favorites_count = GREATEST(0, COALESCE(favorites_count, 0) - 1) WHERE artist_id = $1`, artistID)
		if err != nil {
			return fmt.Errorf("decrement artist favorites count: %w", err)
		}

	case domain.FavoriteTargetVenue:
		venueID, _ := strconv.Atoi(target.ID)
		if delta > 0 {
			_, err := exec.Exec(ctx, `
				INSERT INTO venue_stats (venue_id, favorites_count)
				VALUES ($1, 1)
				ON CONFLICT (venue_id) DO UPDATE
				SET favorites_count = COALESCE(venue_stats.favorites_count, 0) + 1
			`, venueID)
			if err != nil {
				return fmt.Errorf("increment venue favorites count: %w", err)
			}
			return nil
		}
		_, err := exec.Exec(ctx, `UPDATE venue_stats SET favorites_count = GREATEST(0, COALESCE(favorites_count, 0) - 1) WHERE venue_id = $1`, venueID)
		if err != nil {
			return fmt.Errorf("decrement venue favorites count: %w", err)
		}

	case domain.FavoriteTargetConcert:
		concertID, _ := uuid.Parse(target.ID)
		if delta > 0 {
			_, err := exec.Exec(ctx, `
				INSERT INTO concert_stats (concert_id, favorites_count)
				VALUES ($1, 1)
				ON CONFLICT (concert_id) DO UPDATE
				SET favorites_count = COALESCE(concert_stats.favorites_count, 0) + 1
			`, concertID)
			if err != nil {
				return fmt.Errorf("increment concert favorites count: %w", err)
			}
			return nil
		}
		_, err := exec.Exec(ctx, `UPDATE concert_stats SET favorites_count = GREATEST(0, COALESCE(favorites_count, 0) - 1) WHERE concert_id = $1`, concertID)
		if err != nil {
			return fmt.Errorf("decrement concert favorites count: %w", err)
		}
	}

	return nil
}

func mapFavoriteWriteError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return fmt.Errorf("write favorite: %w", err)
	}

	switch pgErr.Code {
	case "23505":
		return fmt.Errorf("favorite already exists: %w", core_errors.ErrConflict)
	case "23514":
		return fmt.Errorf("favorite limit exceeded: %w", core_errors.ErrConflict)
	default:
		return fmt.Errorf("write favorite: %w", err)
	}
}
