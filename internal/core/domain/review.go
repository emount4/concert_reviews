package domain

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ModerationStatus string

const (
	StatusPending  ModerationStatus = "pending"
	StatusApproved ModerationStatus = "approved"
	StatusRejected ModerationStatus = "rejected"
)

// Review представляет доменную модель рецензии
type Review struct {
	ReviewID     uuid.UUID
	UserID       uuid.UUID
	ConcertID    uuid.UUID
	Title        string
	Text         string
	OriginalText *string

	// Оценки параметров (1-10)
	P1 int // Звук
	P2 int // Свет
	P3 int // Исполнение
	P4 int // Атмосфера
	P5 int // Вайб (Коэффициент)

	RatingTotal       int
	Status            ModerationStatus
	RejectionReason   *string
	ModeratedByUserID *uuid.UUID
	IsVisible         bool
	CreatedAt         time.Time
	DeletedAt         *time.Time

	// Поля обогащения (заполняются репозиторием через JOIN)
	AuthorName       string
	AuthorAvatar     *string
	ConcertTitle     string
	ConcertPosterURL *string
	ConcertArtists   []ConcertArtist
	Media            []ReviewMedia
	LikesCount       int
	IsLikedByMe      bool
}

// ReviewMedia представляет файл, прикрепленный к рецензии
type ReviewMedia struct {
	MediaID   uuid.UUID
	ReviewID  uuid.UUID
	MediaURL  string // Ключ S3
	MediaType string // image, video
	FileSize  *int64
	Status    ModerationStatus
	CreatedAt time.Time
}

// (p1+p2+p3+p4) * 1.4 * K, где K зависит от p5 (линейный прирост)
func (r *Review) CalculateRating() int {
	// 1. Коэффициент K для P5.
	// При P5=1, K=1.0. При P5=10, K=1.6072.
	// Формула K: 1.0 + (p5 - 1) * ((1.6072 - 1.0) / 9)
	step := (1.6072 - 1.0) / 9.0
	k := 1.0 + float64(r.P5-1)*step

	// 2. Основная часть формулы
	sumBase := float64(r.P1 + r.P2 + r.P3 + r.P4)
	result := sumBase * 1.4 * k

	// 3. Округляем до ближайшего целого и ограничиваем максимумом 90
	finalScore := math.Round(result)
	if finalScore > 90 {
		finalScore = 90
	}

	return int(finalScore)
}

// Validate проверяет рецензию на соответствие правилам системы и БД
func (r *Review) Validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if len(r.Title) > 255 {
		return fmt.Errorf("title too long (max 255)")
	}

	// Проверка длины текста (согласно твоей миграции 100-8000)
	textLen := len(r.Text)
	if textLen < 100 || textLen > 8000 {
		return fmt.Errorf("text length must be between 100 and 8000 symbols (current: %d)", textLen)
	}

	// Проверка оценок
	params := []struct {
		val  int
		name string
	}{
		{r.P1, "p1"}, {r.P2, "p2"}, {r.P3, "p3"}, {r.P4, "p4"}, {r.P5, "p5"},
	}

	for _, p := range params {
		if p.val < 1 || p.val > 10 {
			return fmt.Errorf("parameter %s must be between 1 and 10", p.name)
		}
	}

	return nil
}

func (r *Review) IsApproved() bool { return r.Status == StatusApproved }
func (r *Review) IsPending() bool  { return r.Status == StatusPending }
func (r *Review) IsRejected() bool { return r.Status == StatusRejected }
