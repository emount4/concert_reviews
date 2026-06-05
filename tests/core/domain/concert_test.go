package domain_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	core_types "github.com/emount4/concert_reviews/internal/core/types"
)

func TestConcertValidateSuccess(t *testing.T) {
	concert := domain.Concert{
		VenueID: 1,
		Title:   "Live show",
		Date:    time.Now().Add(24 * time.Hour),
	}

	if err := concert.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
}

func TestConcertValidateRejectsEmptyTitle(t *testing.T) {
	concert := domain.Concert{
		VenueID: 1,
		Title:   "   ",
		Date:    time.Now().Add(24 * time.Hour),
	}

	if err := concert.Validate(); err == nil {
		t.Fatal("Validate() expected error for empty title, got nil")
	}
}

func TestConcertValidateRejectsInvalidVenueID(t *testing.T) {
	concert := domain.Concert{
		VenueID: 0,
		Title:   "Show",
		Date:    time.Now().Add(24 * time.Hour),
	}

	if err := concert.Validate(); err == nil {
		t.Fatal("Validate() expected error for invalid venue_id, got nil")
	}
}

func TestConcertValidateRejectsZeroDate(t *testing.T) {
	concert := domain.Concert{
		VenueID: 1,
		Title:   "Show",
		Date:    time.Time{},
	}

	if err := concert.Validate(); err == nil {
		t.Fatal("Validate() expected error for zero date, got nil")
	}
}

func TestConcertValidateRejectsEmptyPosterKey(t *testing.T) {
	empty := "   "
	concert := domain.Concert{
		VenueID: 1,
		Title:   "Show",
		Date:    time.Now().Add(24 * time.Hour),
		PosterKey: &empty,
	}

	if err := concert.Validate(); err == nil {
		t.Fatal("Validate() expected error for empty poster_key, got nil")
	}
}

func TestConcertValidateRejectsLongPosterKey(t *testing.T) {
	longKey := strings.Repeat("a", 2049)
	concert := domain.Concert{
		VenueID: 1,
		Title:   "Show",
		Date:    time.Now().Add(24 * time.Hour),
		PosterKey: &longKey,
	}

	if err := concert.Validate(); err == nil {
		t.Fatal("Validate() expected error for long poster_key, got nil")
	}
}

func TestConcertPatchValidateRejectsEmptyTitle(t *testing.T) {
	empty := "   "
	patch := domain.ConcertPatch{
		Title: core_types.Nullable[string]{Value: &empty, Set: true},
	}

	if err := patch.Validate(); err == nil {
		t.Fatal("Validate() expected error for empty title patch, got nil")
	}
}

func TestConcertPatchValidateRejectsInvalidVenueID(t *testing.T) {
	zero := 0
	patch := domain.ConcertPatch{
		VenueID: core_types.Nullable[int]{Value: &zero, Set: true},
	}

	if err := patch.Validate(); err == nil {
		t.Fatal("Validate() expected error for invalid venue_id patch, got nil")
	}
	if err := patch.Validate(); !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("Validate() expected invalid argument, got %v", err)
	}
}

func TestConcertApplyPatchUpdatesFields(t *testing.T) {
	oldDate := time.Now().Add(24 * time.Hour)
	newDate := time.Now().Add(48 * time.Hour)
	newTitle := "New title"
	newVenueID := 2
	verified := true

	concert := domain.Concert{
		VenueID: 1,
		Title:   "Old title",
		Date:    oldDate,
	}

	patch := domain.ConcertPatch{
		VenueID:    core_types.Nullable[int]{Value: &newVenueID, Set: true},
		Title:      core_types.Nullable[string]{Value: &newTitle, Set: true},
		Date:       core_types.Nullable[time.Time]{Value: &newDate, Set: true},
		IsVerified: core_types.Nullable[bool]{Value: &verified, Set: true},
	}

	if err := concert.ApplyPatch(patch); err != nil {
		t.Fatalf("ApplyPatch() unexpected error: %v", err)
	}

	if concert.VenueID != newVenueID || concert.Title != newTitle || !concert.IsVerified {
		t.Fatalf("ApplyPatch() fields not updated: %+v", concert)
	}
	if !concert.Date.Equal(newDate) {
		t.Fatalf("ApplyPatch() date not updated: %v", concert.Date)
	}
}

func TestConcertStatsAvgRatingZeroCases(t *testing.T) {
	var stats *domain.ConcertStats
	if stats.AvgRating() != 0 {
		t.Fatal("AvgRating() expected 0 for nil stats")
	}

	stats = &domain.ConcertStats{}
	if stats.AvgRating() != 0 {
		t.Fatal("AvgRating() expected 0 for zero reviews")
	}
}

func TestConcertStatsAvgByParamInvalidParam(t *testing.T) {
	stats := &domain.ConcertStats{ReviewsCount: 2, SumP1: 10}
	if stats.AvgByParam(6) != 0 {
		t.Fatal("AvgByParam() expected 0 for invalid param")
	}
}
