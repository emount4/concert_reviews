package domain_test

import (
	"errors"
	"testing"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
)

func TestParseFavoriteTargetTypeSuccess(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  domain.FavoriteTargetType
	}{
		{"artist", "artist", domain.FavoriteTargetArtist},
		{"venue", "venue", domain.FavoriteTargetVenue},
		{"concert", "concert", domain.FavoriteTargetConcert},
	}

	for _, tc := range cases {
		got, err := domain.ParseFavoriteTargetType(tc.input)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.name, err)
		}
		if got != tc.want {
			t.Fatalf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestParseFavoriteTargetTypeRejectsInvalid(t *testing.T) {
	_, err := domain.ParseFavoriteTargetType("bad")
	if err == nil {
		t.Fatal("expected error for invalid target type, got nil")
	}
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}
