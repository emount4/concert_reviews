// venue_postgres_repository/record.go

package venue_postgres_repository

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
)

// VenueRecord — представление строки таблицы venues в БД
type VenueRecord struct {
	VenueID     int               `db:"venue_id"`
	CityID      int               `db:"city_id"`
	CityName    string            `db:"city_name"`
	CitySlug    string            `db:"city_slug"`
	Name        string            `db:"name"`
	Address     *string           `db:"address"`
	Capacity    *int              `db:"capacity"`
	SocialLinks map[string]string `db:"social_links"`
	PhotoURL    *string           `db:"photo_url"`
	Description *string           `db:"description"`
	Status      string            `db:"status"`
	CreatedAt   time.Time         `db:"created_at"`
	DeletedAt   *time.Time        `db:"deleted_at"`
}

// MapToDomain конвертирует запись БД в доменную модель
func (r VenueRecord) MapToDomain() domain.Venue {
	return domain.Venue{
		VenueID: r.VenueID,
		CityID:  r.CityID,
		City: &domain.City{
			CityID: r.CityID,
			Name:   r.CityName,
			Slug:   r.CitySlug,
		},
		Name:        r.Name,
		Address:     r.Address,
		Capacity:    r.Capacity,
		SocialLinks: r.SocialLinks,
		PhotoURL:    r.PhotoURL,
		Description: r.Description,
		Status:      domain.ContentStatus(r.Status),
		CreatedAt:   r.CreatedAt,
		DeletedAt:   r.DeletedAt,
	}
}
