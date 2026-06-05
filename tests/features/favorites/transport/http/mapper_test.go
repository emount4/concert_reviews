package favorites_transport_http_test

import (
	"testing"
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	favorites_http "github.com/emount4/concert_reviews/internal/features/favorites/transport/http"
	"github.com/google/uuid"
)

func TestMapDomainToResponseIncludesPresentationFields(t *testing.T) {
	imageURL := "artists/artist.jpg"
	createdAt := time.Date(2026, 5, 31, 9, 30, 0, 0, time.UTC)
	favorite := domain.Favorite{
		FavoriteID: 12,
		UserID:     uuid.New(),
		Target: domain.FavoriteTarget{
			Type:     domain.FavoriteTargetArtist,
			ID:       "42",
			Name:     "Artist Name",
			ImageURL: &imageURL,
		},
		CreatedAt: createdAt,
	}

	response := favorites_http.MapDomainToResponse(favorite)

	if response.ID != 12 {
		t.Fatalf("expected id 12, got %d", response.ID)
	}
	if response.TargetType != "artist" || response.TargetID != "42" {
		t.Fatalf("unexpected target fields: %+v", response)
	}
	if response.Name != "Artist Name" {
		t.Fatalf("expected target name, got %q", response.Name)
	}
	if response.ImageURL == nil || *response.ImageURL != imageURL {
		t.Fatalf("expected image url %q, got %v", imageURL, response.ImageURL)
	}
	if response.CreatedAt != createdAt.Format(time.RFC3339) {
		t.Fatalf("expected RFC3339 created_at, got %q", response.CreatedAt)
	}
}

func TestMapDomainListToResponseKeepsOrder(t *testing.T) {
	favorites := []domain.Favorite{
		{FavoriteID: 1, Target: domain.FavoriteTarget{Type: domain.FavoriteTargetArtist, ID: "1", Name: "Artist"}},
		{FavoriteID: 2, Target: domain.FavoriteTarget{Type: domain.FavoriteTargetVenue, ID: "2", Name: "Venue"}},
	}

	response := favorites_http.MapDomainListToResponse(favorites)

	if len(response.Items) != 2 {
		t.Fatalf("expected two items, got %d", len(response.Items))
	}
	if response.Items[0].ID != 1 || response.Items[1].ID != 2 {
		t.Fatalf("expected mapper to keep order, got %+v", response.Items)
	}
	if response.Items[0].Name != "Artist" || response.Items[1].TargetType != "venue" {
		t.Fatalf("unexpected mapped items: %+v", response.Items)
	}
}
