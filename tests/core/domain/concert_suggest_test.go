package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

func TestConcertSuggestionValidateSuccess(t *testing.T) {
	suggestion := domain.ConcertSuggestion{
		SuggestionID:  uuid.New(),
		UserID:        uuid.New(),
		RawArtistName: "Artist",
		ConcertDate:   time.Now().Add(24 * time.Hour),
	}

	if err := suggestion.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
}

func TestConcertSuggestionValidateRejectsMissingDate(t *testing.T) {
	suggestion := domain.ConcertSuggestion{RawArtistName: "Artist"}
	if err := suggestion.Validate(); err == nil {
		t.Fatal("Validate() expected error for missing date, got nil")
	}
}

func TestConcertSuggestionValidateRequiresArtistOrVenue(t *testing.T) {
	suggestion := domain.ConcertSuggestion{ConcertDate: time.Now().Add(24 * time.Hour)}
	if err := suggestion.Validate(); err == nil {
		t.Fatal("Validate() expected error for missing names, got nil")
	}
}

func TestConcertSuggestionValidateRejectsLongFields(t *testing.T) {
	long := strings.Repeat("a", 256)
	suggestion := domain.ConcertSuggestion{
		RawArtistName: long,
		ConcertDate:   time.Now().Add(24 * time.Hour),
	}
	if err := suggestion.Validate(); err == nil {
		t.Fatal("Validate() expected error for long artist name, got nil")
	}

	info := strings.Repeat("b", 2001)
	suggestion = domain.ConcertSuggestion{
		RawVenueName: "Venue",
		ConcertDate:  time.Now().Add(24 * time.Hour),
		Info:         &info,
	}
	if err := suggestion.Validate(); err == nil {
		t.Fatal("Validate() expected error for long info, got nil")
	}
}
