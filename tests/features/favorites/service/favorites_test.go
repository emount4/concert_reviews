package favorites_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	favorites_service "github.com/emount4/concert_reviews/internal/features/favorites/service"
	"github.com/google/uuid"
)

type fakeFavoritesRepository struct {
	addFavoriteFn            func(ctx context.Context, userID uuid.UUID, target domain.FavoriteTarget) (domain.Favorite, error)
	deleteFavoriteFn         func(ctx context.Context, userID uuid.UUID, target domain.FavoriteTarget) error
	countFavoritesByTypeFn   func(ctx context.Context, userID uuid.UUID, targetType domain.FavoriteTargetType) (int, error)
	getFavoriteTargetFn      func(ctx context.Context, targetType domain.FavoriteTargetType, targetID string) (domain.FavoriteTarget, error)
	getFavoritesByUsernameFn func(ctx context.Context, username string, targetType *domain.FavoriteTargetType) ([]domain.Favorite, error)

	addCalls       int
	deleteCalls    int
	countCalls     int
	getTargetCalls int
	listCalls      int

	lastUserID     uuid.UUID
	lastTarget     domain.FavoriteTarget
	lastTargetType domain.FavoriteTargetType
	lastTargetID   string
	lastUsername   string
	lastTypeFilter *domain.FavoriteTargetType
}

func (f *fakeFavoritesRepository) AddFavorite(ctx context.Context, userID uuid.UUID, target domain.FavoriteTarget) (domain.Favorite, error) {
	f.addCalls++
	f.lastUserID = userID
	f.lastTarget = target
	if f.addFavoriteFn != nil {
		return f.addFavoriteFn(ctx, userID, target)
	}
	return domain.Favorite{FavoriteID: 1, UserID: userID, Target: target, CreatedAt: time.Now()}, nil
}

func (f *fakeFavoritesRepository) DeleteFavorite(ctx context.Context, userID uuid.UUID, target domain.FavoriteTarget) error {
	f.deleteCalls++
	f.lastUserID = userID
	f.lastTarget = target
	if f.deleteFavoriteFn != nil {
		return f.deleteFavoriteFn(ctx, userID, target)
	}
	return nil
}

func (f *fakeFavoritesRepository) CountFavoritesByType(ctx context.Context, userID uuid.UUID, targetType domain.FavoriteTargetType) (int, error) {
	f.countCalls++
	f.lastUserID = userID
	f.lastTargetType = targetType
	if f.countFavoritesByTypeFn != nil {
		return f.countFavoritesByTypeFn(ctx, userID, targetType)
	}
	return 0, nil
}

func (f *fakeFavoritesRepository) GetFavoriteTarget(ctx context.Context, targetType domain.FavoriteTargetType, targetID string) (domain.FavoriteTarget, error) {
	f.getTargetCalls++
	f.lastTargetType = targetType
	f.lastTargetID = targetID
	if f.getFavoriteTargetFn != nil {
		return f.getFavoriteTargetFn(ctx, targetType, targetID)
	}
	return domain.FavoriteTarget{Type: targetType, ID: targetID, Name: "Target"}, nil
}

func (f *fakeFavoritesRepository) GetFavoritesByUsername(ctx context.Context, username string, targetType *domain.FavoriteTargetType) ([]domain.Favorite, error) {
	f.listCalls++
	f.lastUsername = username
	f.lastTypeFilter = targetType
	if f.getFavoritesByUsernameFn != nil {
		return f.getFavoritesByUsernameFn(ctx, username, targetType)
	}
	return []domain.Favorite{{FavoriteID: 1, Target: domain.FavoriteTarget{Type: domain.FavoriteTargetArtist, ID: "1", Name: "Artist"}}}, nil
}

func TestFavoritesServiceAddFavoriteSuccess(t *testing.T) {
	userID := uuid.New()
	imageURL := "artist.jpg"
	repo := &fakeFavoritesRepository{
		getFavoriteTargetFn: func(ctx context.Context, targetType domain.FavoriteTargetType, targetID string) (domain.FavoriteTarget, error) {
			return domain.FavoriteTarget{
				Type:     targetType,
				ID:       targetID,
				Name:     "Artist",
				ImageURL: &imageURL,
			}, nil
		},
	}
	service := favorites_service.NewFavoritesService(repo)

	favorite, err := service.AddFavorite(context.Background(), userID, " artist ", " 42 ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if favorite.UserID != userID {
		t.Fatalf("expected favorite user id %s, got %s", userID, favorite.UserID)
	}
	if favorite.Target.Type != domain.FavoriteTargetArtist || favorite.Target.ID != "42" {
		t.Fatalf("unexpected target: %+v", favorite.Target)
	}
	if favorite.Target.Name != "Artist" || favorite.Target.ImageURL == nil || *favorite.Target.ImageURL != imageURL {
		t.Fatalf("expected target presentation fields to be returned, got %+v", favorite.Target)
	}
	if repo.getTargetCalls != 1 || repo.countCalls != 1 || repo.addCalls != 1 {
		t.Fatalf("expected target lookup, count and add once, got target=%d count=%d add=%d", repo.getTargetCalls, repo.countCalls, repo.addCalls)
	}
}

func TestFavoritesServiceAddFavoriteRejectsEmptyUserID(t *testing.T) {
	repo := &fakeFavoritesRepository{}
	service := favorites_service.NewFavoritesService(repo)

	_, err := service.AddFavorite(context.Background(), uuid.Nil, "artist", "1")
	if !errors.Is(err, core_errors.ErrUnauthorized) {
		t.Fatalf("expected unauthorized, got %v", err)
	}
	if repo.getTargetCalls != 0 || repo.countCalls != 0 || repo.addCalls != 0 {
		t.Fatalf("expected no repository calls, got target=%d count=%d add=%d", repo.getTargetCalls, repo.countCalls, repo.addCalls)
	}
}

func TestFavoritesServiceAddFavoriteRejectsInvalidTargetID(t *testing.T) {
	repo := &fakeFavoritesRepository{}
	service := favorites_service.NewFavoritesService(repo)

	_, err := service.AddFavorite(context.Background(), uuid.New(), "venue", "not-int")
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
	if repo.getTargetCalls != 0 || repo.countCalls != 0 || repo.addCalls != 0 {
		t.Fatalf("expected no repository calls after validation error, got target=%d count=%d add=%d", repo.getTargetCalls, repo.countCalls, repo.addCalls)
	}
}

func TestFavoritesServiceAddFavoriteRejectsMaxPerType(t *testing.T) {
	repo := &fakeFavoritesRepository{
		countFavoritesByTypeFn: func(ctx context.Context, userID uuid.UUID, targetType domain.FavoriteTargetType) (int, error) {
			return favorites_service.MaxFavoritesPerType, nil
		},
	}
	service := favorites_service.NewFavoritesService(repo)

	_, err := service.AddFavorite(context.Background(), uuid.New(), "artist", "1")
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
	if repo.getTargetCalls != 1 || repo.countCalls != 1 || repo.addCalls != 0 {
		t.Fatalf("expected lookup and count without add, got target=%d count=%d add=%d", repo.getTargetCalls, repo.countCalls, repo.addCalls)
	}
}

func TestFavoritesServiceDeleteFavoriteTrimsAndValidatesTarget(t *testing.T) {
	userID := uuid.New()
	concertID := uuid.New().String()
	repo := &fakeFavoritesRepository{}
	service := favorites_service.NewFavoritesService(repo)

	if err := service.DeleteFavorite(context.Background(), userID, " concert ", " "+concertID+" "); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.deleteCalls != 1 {
		t.Fatalf("expected DeleteFavorite once, got %d", repo.deleteCalls)
	}
	if repo.lastTarget.Type != domain.FavoriteTargetConcert || repo.lastTarget.ID != concertID {
		t.Fatalf("unexpected delete target: %+v", repo.lastTarget)
	}
}

func TestFavoritesServiceGetFavoritesByUsernameTrimsAndFiltersType(t *testing.T) {
	repo := &fakeFavoritesRepository{}
	service := favorites_service.NewFavoritesService(repo)
	targetType := " venue "

	favorites, err := service.GetFavoritesByUsername(context.Background(), " user_name ", &targetType)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(favorites) != 1 {
		t.Fatalf("expected one favorite, got %d", len(favorites))
	}
	if repo.lastUsername != "user_name" {
		t.Fatalf("expected trimmed username, got %q", repo.lastUsername)
	}
	if repo.lastTypeFilter == nil || *repo.lastTypeFilter != domain.FavoriteTargetVenue {
		t.Fatalf("expected venue filter, got %v", repo.lastTypeFilter)
	}
}

func TestFavoritesServiceGetFavoritesByUsernameRejectsEmptyUsername(t *testing.T) {
	repo := &fakeFavoritesRepository{}
	service := favorites_service.NewFavoritesService(repo)

	_, err := service.GetFavoritesByUsername(context.Background(), "   ", nil)
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
	if repo.listCalls != 0 {
		t.Fatalf("expected no list repository calls, got %d", repo.listCalls)
	}
}
