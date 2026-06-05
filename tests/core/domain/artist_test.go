package domain_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_http_types "github.com/emount4/concert_reviews/internal/core/transport/http/types"
	core_types "github.com/emount4/concert_reviews/internal/core/types"
)

func TestArtistValidateSuccess(t *testing.T) {
	artist := domain.Artist{
		Name: "Artist",
	}

	if err := artist.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
}

func TestArtistValidateRejectsEmptyName(t *testing.T) {
	artist := domain.Artist{Name: "  "}
	if err := artist.Validate(); err == nil {
		t.Fatal("Validate() expected error for empty name, got nil")
	}
}

func TestArtistValidateRejectsLongDescription(t *testing.T) {
	desc := strings.Repeat("a", 2001)
	artist := domain.Artist{Name: "Artist", Description: &desc}
	if err := artist.Validate(); err == nil {
		t.Fatal("Validate() expected error for long description, got nil")
	}
}

func TestArtistValidateRejectsEmptyPhotoURL(t *testing.T) {
	photo := "  "
	artist := domain.Artist{Name: "Artist", PhotoURL: &photo}
	if err := artist.Validate(); err == nil {
		t.Fatal("Validate() expected error for empty photo_url, got nil")
	}
}

func TestArtistIsDeleted(t *testing.T) {
	now := time.Now()
	artist := domain.Artist{DeletedAt: &now}
	if !artist.IsDeleted() {
		t.Fatal("IsDeleted() expected true")
	}
}

func TestArtistApplyPatchRejectsNullName(t *testing.T) {
	artist := domain.Artist{Name: "Artist"}
	patch := domain.ArtistPatch{
		Name: core_types.Nullable[string]{Value: nil, Set: true},
	}
	err := artist.ApplyPatch(patch)
	if err == nil {
		t.Fatal("ApplyPatch() expected error for null name patch, got nil")
	}
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatal("ApplyPatch() expected invalid argument error")
	}
}

func TestArtistApplyPatchUpdatesSocialLinks(t *testing.T) {
	artist := domain.Artist{Name: "Artist"}
	links := map[string]string{"site": "example.com"}
	patch := domain.ArtistPatch{
		SocialLinks: core_http_types.NewNullableMapStringString(links, true),
	}

	if err := artist.ApplyPatch(patch); err != nil {
		t.Fatalf("ApplyPatch() unexpected error: %v", err)
	}
	if artist.SocialLinks == nil || artist.SocialLinks["site"] != "example.com" {
		t.Fatalf("ApplyPatch() expected social links to be set, got %+v", artist.SocialLinks)
	}

	patch = domain.ArtistPatch{
		SocialLinks: core_http_types.NullableMapStringStringNull(),
	}
	if err := artist.ApplyPatch(patch); err != nil {
		t.Fatalf("ApplyPatch() unexpected error: %v", err)
	}
	if artist.SocialLinks != nil {
		t.Fatalf("ApplyPatch() expected social links cleared, got %+v", artist.SocialLinks)
	}
}
