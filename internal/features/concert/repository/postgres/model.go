package concert_postgres_repository

import (
	"encoding/json"
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

type concertRecord struct {
	ConcertID       uuid.UUID  `db:"concert_id"`
	VenueID         int        `db:"venue_id"`
	Title           string     `db:"title"`
	Date            time.Time  `db:"date"`
	PosterKey       *string    `db:"poster_url"`
	IsVerified      bool       `db:"is_verified"`
	CreatedByUserID *uuid.UUID `db:"created_by_user_id"`
	CreatedAt       time.Time  `db:"created_at"`
	DeletedAt       *time.Time `db:"deleted_at"`

	// Вложенные данные Площадки
	VenueName string `db:"venue_name"`
	CityID    int    `db:"city_id"`
	CityName  string `db:"city_name"`
	CitySlug  string `db:"city_slug"`

	// Данные статистики (LEFT JOIN может вернуть NULL)
	ReviewsCount   *int       `db:"reviews_count"`
	SumP1          *int       `db:"sum_p1"`
	SumP2          *int       `db:"sum_p2"`
	SumP3          *int       `db:"sum_p3"`
	SumP4          *int       `db:"sum_p4"`
	SumP5          *int       `db:"sum_p5"`
	SumRatingTotal *int64     `db:"sum_rating_total"`
	FavoritesCount int        `db:"favorites_count"`
	StatsUpdatedAt *time.Time `db:"stats_updated_at"`

	// Список артистов в формате JSON
	ArtistsJSON []byte `db:"artists_json"`
}

func (r concertRecord) MapToDomain() domain.Concert {
	c := domain.Concert{
		ConcertID:       r.ConcertID,
		VenueID:         r.VenueID,
		Title:           r.Title,
		Date:            r.Date,
		PosterKey:       r.PosterKey,
		IsVerified:      r.IsVerified,
		CreatedByUserID: r.CreatedByUserID,
		CreatedAt:       r.CreatedAt,
		DeletedAt:       r.DeletedAt,
		// Собираем вложенную площадку
		Venue: &domain.Venue{
			VenueID: r.VenueID,
			Name:    r.VenueName,
			City: &domain.City{
				CityID: r.CityID,
				Name:   r.CityName,
				Slug:   r.CitySlug,
			},
		},
	}

	// Маппим статистику, если она есть
	if r.ReviewsCount != nil && r.SumP1 != nil && r.SumP2 != nil && r.SumP3 != nil && r.SumP4 != nil && r.SumP5 != nil && r.SumRatingTotal != nil && r.StatsUpdatedAt != nil {
		c.Stats = &domain.ConcertStats{
			ReviewsCount:   *r.ReviewsCount,
			SumP1:          *r.SumP1,
			SumP2:          *r.SumP2,
			SumP3:          *r.SumP3,
			SumP4:          *r.SumP4,
			SumP5:          *r.SumP5,
			SumRatingTotal: *r.SumRatingTotal,
			FavoritesCount: r.FavoritesCount,
			UpdatedAt:      *r.StatsUpdatedAt,
		}
	}

	// Распаковываем артистов
	if len(r.ArtistsJSON) > 0 {
		var artists []domain.ConcertArtist
		if err := json.Unmarshal(r.ArtistsJSON, &artists); err == nil {
			c.Artists = artists
		}
	}

	return c
}
