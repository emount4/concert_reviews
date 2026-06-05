package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_http_types "github.com/emount4/concert_reviews/internal/core/transport/http/types"
	core_types "github.com/emount4/concert_reviews/internal/core/types"
)

func TestVenueValidateSuccess(t *testing.T) {
	venue := domain.Venue{CityID: 1, Name: "Venue"}
	if err := venue.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
}

func TestVenueValidateRejectsInvalidCityID(t *testing.T) {
	venue := domain.Venue{CityID: 0, Name: "Venue"}
	if err := venue.Validate(); err == nil {
		t.Fatal("Validate() expected error for invalid city_id, got nil")
	}
}

func TestVenueValidateRejectsEmptyAddress(t *testing.T) {
	addr := "  "
	venue := domain.Venue{CityID: 1, Name: "Venue", Address: &addr}
	if err := venue.Validate(); err == nil {
		t.Fatal("Validate() expected error for empty address, got nil")
	}
}

func TestVenueValidateRejectsInvalidCapacity(t *testing.T) {
	cap := 0
	venue := domain.Venue{CityID: 1, Name: "Venue", Capacity: &cap}
	if err := venue.Validate(); err == nil {
		t.Fatal("Validate() expected error for invalid capacity, got nil")
	}
}

func TestVenueValidateRejectsEmptyPhotoURL(t *testing.T) {
	photo := " "
	venue := domain.Venue{CityID: 1, Name: "Venue", PhotoURL: &photo}
	if err := venue.Validate(); err == nil {
		t.Fatal("Validate() expected error for empty photo_url, got nil")
	}
}

func TestVenuePatchValidateRejectsInvalidStatus(t *testing.T) {
	status := "bad"
	patch := domain.VenuePatch{Status: core_types.Nullable[string]{Value: &status, Set: true}}
	err := patch.Validate()
	if err == nil {
		t.Fatal("Validate() expected error for invalid status, got nil")
	}
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatal("Validate() expected invalid argument error")
	}
}

func TestVenueApplyPatchUpdatesFields(t *testing.T) {
	venue := domain.Venue{CityID: 1, Name: "Venue"}
	links := map[string]string{"site": "example.com"}
	cityID := 2
	name := "New Venue"
	status := string(domain.StatusHidden)

	patch := domain.VenuePatch{
		CityID:      core_types.Nullable[int]{Value: &cityID, Set: true},
		Name:        core_types.Nullable[string]{Value: &name, Set: true},
		SocialLinks: core_http_types.NewNullableMapStringString(links, true),
		Status:      core_types.Nullable[string]{Value: &status, Set: true},
	}

	if err := venue.ApplyPatch(patch); err != nil {
		t.Fatalf("ApplyPatch() unexpected error: %v", err)
	}
	if venue.CityID != 2 || venue.Name != "New Venue" || venue.Status != domain.StatusHidden {
		t.Fatalf("ApplyPatch() unexpected venue state: %+v", venue)
	}
	if venue.SocialLinks == nil || venue.SocialLinks["site"] != "example.com" {
		t.Fatalf("ApplyPatch() expected social links to be set, got %+v", venue.SocialLinks)
	}

	patch = domain.VenuePatch{SocialLinks: core_http_types.NullableMapStringStringNull()}
	if err := venue.ApplyPatch(patch); err != nil {
		t.Fatalf("ApplyPatch() unexpected error: %v", err)
	}
	if venue.SocialLinks != nil {
		t.Fatalf("ApplyPatch() expected social links cleared, got %+v", venue.SocialLinks)
	}
}

func TestVenueValidateRejectsLongDescription(t *testing.T) {
	desc := strings.Repeat("a", 2001)
	venue := domain.Venue{CityID: 1, Name: "Venue", Description: &desc}
	if err := venue.Validate(); err == nil {
		t.Fatal("Validate() expected error for long description, got nil")
	}
}
