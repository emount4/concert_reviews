package artist_transport_http

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
)

type CreateArtistRequest struct {
	Name        string            `json:"name" validate:"required,min=2,max=255"`
	Description string            `json:"description" validate:"max=2000"`
	PhotoKey    string            `json:"photo_key" validate:"omitempty,max=2048"`
	SocialLinks map[string]string `json:"social_links"`
}

// type UpdateArtistRequest struct {
// 	Name        *string           `json:"name" validate:"omitempty,min=2,max=255"`
// 	Description *string           `json:"description" validate:"omitempty,max=2000"`
// 	PhotoKey    *string           `json:"photo_key" validate:"omitempty,max=2048"`
// 	SocialLinks map[string]string `json:"social_links"`
// }

// --- Responses (DTO для выхода) ---

type ArtistResponse struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	PhotoURL    string            `json:"photo_url,omitempty"`
	SocialLinks map[string]string `json:"social_links,omitempty"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"created_at"`
}

type ListArtistsResponse struct {
	Items []ArtistResponse `json:"items"`
}

// --- Mappers (Конвертация) ---

// MapCreateDTOToDomain преобразует данные запроса в чистую доменную модель
func MapCreateDTOToDomain(dto CreateArtistRequest) domain.Artist {
	var desc *string
	if dto.Description != "" {
		desc = &dto.Description
	}

	var photo *string
	if dto.PhotoKey != "" {
		photo = &dto.PhotoKey
	}

	return domain.Artist{
		Name:        dto.Name,
		Description: desc,
		PhotoURL:    photo,
		SocialLinks: dto.SocialLinks,
		Status:      domain.StatusActive,
	}
}

// MapDomainToResponse преобразует доменную модель в формат JSON ответа
func MapDomainToResponse(a domain.Artist) ArtistResponse {
	resp := ArtistResponse{
		ID:        a.ArtistID,
		Name:      a.Name,
		Status:    string(a.Status),
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	}

	// Безопасно разыменовываем указатели
	if a.Description != nil {
		resp.Description = *a.Description
	}

	if a.PhotoURL != nil {
		resp.PhotoURL = *a.PhotoURL
	}

	if a.SocialLinks != nil {
		resp.SocialLinks = a.SocialLinks
	}

	return resp
}

// MapDomainListToResponse преобразует срез доменных моделей в DTO ответа
func MapDomainListToResponse(artists []domain.Artist) ListArtistsResponse {
	items := make([]ArtistResponse, len(artists))
	for i, a := range artists {
		items[i] = MapDomainToResponse(a)
	}
	return ListArtistsResponse{Items: items}
}
