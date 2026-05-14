package concert_service

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_ports "github.com/emount4/concert_reviews/internal/core/domain/ports"
	"github.com/google/uuid"
)

type ConcertService struct {
	concertRepository ConcertRepository
	s3                core_ports.S3Provider
}

type ConcertRepository interface {
	// Методы для основных концертов
	CreateConcert(ctx context.Context, concert domain.Concert, artists []domain.ConcertArtist) (domain.Concert, error)
	GetConcerts(
		ctx context.Context,
		cityID *int,
		artistID *int,
		search string,
		sort string,
		direction string,
		limit, offset *int,
	) ([]domain.Concert, int, error)
	GetConcertByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (domain.Concert, error)
	UpdateConcert(ctx context.Context, id uuid.UUID, patch domain.ConcertPatch) (domain.Concert, error)
	DeleteConcertHard(ctx context.Context, id uuid.UUID) error
	DeleteConcertSoft(ctx context.Context, id uuid.UUID) error
	RestoreConcert(ctx context.Context, id uuid.UUID) (domain.Concert, error)

	// Админские методы для концертов
	GetConcertsAdmin(
		ctx context.Context,
		cityID *int,
		artistID *int,
		search string,
		sort string,
		direction string,
		limit, offset *int,
		includeDeleted bool,
	) ([]domain.Concert, int, error)

	// Методы для работы с артистами в концертах
	AddConcertArtist(ctx context.Context, concertID uuid.UUID, artistID int, isMain bool) error
	RemoveConcertArtist(ctx context.Context, concertID uuid.UUID, artistID int) error
	UpdateConcertArtistIsMain(ctx context.Context, concertID uuid.UUID, artistID int, isMain bool) error

	// Методы для предложений (Suggestions)
	CreateSuggestion(ctx context.Context, suggestion domain.ConcertSuggestion) (domain.ConcertSuggestion, error)
	GetSuggestionsAdmin(ctx context.Context, limit, offset *int, status string) ([]domain.ConcertSuggestion, error)
	GetSuggestionByIDAdmin(ctx context.Context, id uuid.UUID) (domain.ConcertSuggestion, error)
	DeleteSuggestionAdmin(ctx context.Context, id uuid.UUID) error

	CountPendingSuggestions(ctx context.Context, userID uuid.UUID) (int, error)
}

func NewConcertService(
	concertRepository ConcertRepository,
	s3 core_ports.S3Provider,
) *ConcertService {
	return &ConcertService{
		concertRepository: concertRepository,
		s3:                s3,
	}
}
