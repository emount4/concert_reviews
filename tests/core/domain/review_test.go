package domain_test

import (
	"strings"
	"testing"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

func TestReviewCalculateRatingWithBaseCoefficient(t *testing.T) {
	review := domain.Review{
		UserID:    uuid.New(),
		ConcertID: uuid.New(),
		Title:     "Тестовая рецензия",
		Text:      strings.Repeat("a", 100),
		P1:        5,
		P2:        5,
		P3:        5,
		P4:        5,
		P5:        1,
	}

	got := review.CalculateRating()
	want := 28

	if got != want {
		t.Fatalf("CalculateRating() = %d, want %d", got, want)
	}
}

func TestReviewCalculateRatingMaxValue(t *testing.T) {
	review := domain.Review{
		UserID:    uuid.New(),
		ConcertID: uuid.New(),
		Title:     "Тестовая рецензия",
		Text:      strings.Repeat("a", 100),
		P1:        10,
		P2:        10,
		P3:        10,
		P4:        10,
		P5:        10,
	}

	got := review.CalculateRating()
	want := 90

	if got != want {
		t.Fatalf("CalculateRating() = %d, want %d", got, want)
	}
}

func TestReviewValidateSuccess(t *testing.T) {
	review := domain.Review{
		UserID:    uuid.New(),
		ConcertID: uuid.New(),
		Title:     "Хороший концерт",
		Text:      strings.Repeat("a", 100),
		P1:        8,
		P2:        7,
		P3:        9,
		P4:        8,
		P5:        10,
	}

	if err := review.Validate(); err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}
}

func TestReviewValidateEmptyTitle(t *testing.T) {
	review := domain.Review{
		UserID:    uuid.New(),
		ConcertID: uuid.New(),
		Title:     "",
		Text:      strings.Repeat("a", 100),
		P1:        8,
		P2:        7,
		P3:        9,
		P4:        8,
		P5:        10,
	}

	if err := review.Validate(); err == nil {
		t.Fatal("Validate() expected error for empty title, got nil")
	}
}

func TestReviewValidateShortText(t *testing.T) {
	review := domain.Review{
		UserID:    uuid.New(),
		ConcertID: uuid.New(),
		Title:     "Короткая рецензия",
		Text:      "too short",
		P1:        8,
		P2:        7,
		P3:        9,
		P4:        8,
		P5:        10,
	}

	if err := review.Validate(); err == nil {
		t.Fatal("Validate() expected error for short text, got nil")
	}
}

func TestReviewValidateInvalidScore(t *testing.T) {
	review := domain.Review{
		UserID:    uuid.New(),
		ConcertID: uuid.New(),
		Title:     "Рецензия с ошибкой",
		Text:      strings.Repeat("a", 100),
		P1:        11,
		P2:        7,
		P3:        9,
		P4:        8,
		P5:        10,
	}

	if err := review.Validate(); err == nil {
		t.Fatal("Validate() expected error for invalid score, got nil")
	}
}
