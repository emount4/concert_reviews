package artist_transport_http

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_http_types "github.com/emount4/concert_reviews/internal/core/transport/http/types"
)

type CreateArtistRequest struct {
	Name        string            `json:"name" validate:"required,min=2,max=255"`
	Description string            `json:"description" validate:"max=2000"`
	PhotoKey    string            `json:"photo_key" validate:"omitempty,max=2048"`
	SocialLinks map[string]string `json:"social_links"`
}

type UpdateArtistRequest struct {
	Name        core_http_types.Nullable[string]        `json:"name"`
	Description core_http_types.Nullable[string]        `json:"description"`
	PhotoKey    core_http_types.Nullable[string]        `json:"photo_key"`
	SocialLinks core_http_types.NullableMapStringString `json:"social_links"`
}

// --- Responses (DTO для выхода) ---

type ArtistResponse struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	PhotoURL    string              `json:"photo_url,omitempty"`
	SocialLinks map[string]string   `json:"social_links,omitempty"`
	Stats       ArtistStatsResponse `json:"stats,omitempty"`
	Status      string              `json:"status"`
	CreatedAt   string              `json:"created_at"`
}

type ListArtistsResponse struct {
	Items     []ArtistResponse `json:"items"`
	PageCount int              `json:"page_count"`
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

	if a.Stats != nil {
		resp.Stats = ArtistStatsResponse{
			ReviewsCount:   a.Stats.ReviewsCount,
			SumRatingTotal: a.Stats.SumRatingTotal,
			ConcertsCount:  a.Stats.ConcertsCount,
			FavoritesCount: a.Stats.FavoritesCount,
			UpdatedAt:      a.Stats.UpdatedAt.Format(time.RFC3339),
		}
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

func MapPatchReqToDomain(req UpdateArtistRequest) domain.ArtistPatch {
	return domain.ArtistPatch{
		Name:        req.Name.ToDomain(),
		Description: req.Description.ToDomain(),
		PhotoKey:    req.PhotoKey.ToDomain(),
		SocialLinks: req.SocialLinks,
	}
}

type ArtistAdminResponse struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	PhotoURL    string              `json:"photo_url,omitempty"`
	SocialLinks map[string]string   `json:"social_links,omitempty"`
	Stats       ArtistStatsResponse `json:"stats,omitempty"`
	Status      string              `json:"status"`
	CreatedAt   string              `json:"created_at"`
	DeletedAt   *string             `json:"deleted_at,omitempty"` // показываем только в админке
	IsDeleted   bool                `json:"is_deleted"`
}

type ArtistStatsResponse struct {
	ReviewsCount   int    `json:"reviews_count"`
	SumRatingTotal int64  `json:"sum_rating_total"`
	ConcertsCount  int    `json:"concerts_count"`
	FavoritesCount int    `json:"favorites_count"`
	UpdatedAt      string `json:"updated_at"`
}

type ListArtistsAdminResponse struct {
	Items     []ArtistAdminResponse `json:"items"`
	PageCount int                  `json:"page_count"`
}

func MapDomainToAdminResponse(a domain.Artist) ArtistAdminResponse {
	resp := ArtistAdminResponse{
		ID:        a.ArtistID,
		Name:      a.Name,
		Status:    string(a.Status),
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
		IsDeleted: a.DeletedAt != nil,
	}
	if a.Stats != nil {
		resp.Stats = ArtistStatsResponse{
			ReviewsCount:   a.Stats.ReviewsCount,
			SumRatingTotal: a.Stats.SumRatingTotal,
			ConcertsCount:  a.Stats.ConcertsCount,
			FavoritesCount: a.Stats.FavoritesCount,
			UpdatedAt:      a.Stats.UpdatedAt.Format(time.RFC3339),
		}
	}
	if a.Description != nil {
		resp.Description = *a.Description
	}
	if a.PhotoURL != nil {
		resp.PhotoURL = *a.PhotoURL
	}
	if a.SocialLinks != nil {
		resp.SocialLinks = a.SocialLinks
	}
	if a.DeletedAt != nil {
		s := a.DeletedAt.Format(time.RFC3339)
		resp.DeletedAt = &s
	}
	return resp
}

func MapDomainListToAdminResponse(artists []domain.Artist) ListArtistsAdminResponse {
	items := make([]ArtistAdminResponse, len(artists))
	for i, a := range artists {
		items[i] = MapDomainToAdminResponse(a)
	}
	return ListArtistsAdminResponse{Items: items}
}
