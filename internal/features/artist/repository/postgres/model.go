package artist_postgres_repository

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
)

// ArtistRecord — представление строки в таблице artists
type ArtistRecord struct {
	ArtistID    int               `db:"artist_id"`
	Name        string            `db:"name"`
	Description *string           `db:"description"`
	PhotoURL    *string           `db:"photo_url"`
	SocialLinks map[string]string `db:"social_links"`
	Status      string            `db:"status"`
	CreatedAt   time.Time         `db:"created_at"`
	DeletedAt   *time.Time        `db:"deleted_at"`
}

// MapRecordToDomain — конвертация из БД-модели в Домен
func (r ArtistRecord) MapToDomain() domain.Artist {
	return domain.Artist{
		ArtistID:    r.ArtistID,
		Name:        r.Name,
		Description: r.Description,
		PhotoURL:    r.PhotoURL,
		SocialLinks: r.SocialLinks,
		Status:      domain.ContentStatus(r.Status),
		CreatedAt:   r.CreatedAt,
		DeletedAt:   r.DeletedAt,
	}
}
