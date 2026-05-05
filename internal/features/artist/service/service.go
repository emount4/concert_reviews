package artist_service

import (
	"context"

	"github.com/emount4/concert_reviews/internal/core/domain"
)

type ArtistService struct {
	artistRepository ArtistRepository
}

type ArtistRepository interface {
	CreateArtist(ctx context.Context, artist domain.Artist) (domain.Artist, error)
	GetArtists(ctx context.Context, search string, limit, offset *int) ([]domain.Artist, error)
	GetArtistByID(ctx context.Context, id int) (domain.Artist, error)
}

func NewArtistService(
	artistRepository ArtistRepository,
) *ArtistService {
	return &ArtistService{
		artistRepository: artistRepository,
	}
}
