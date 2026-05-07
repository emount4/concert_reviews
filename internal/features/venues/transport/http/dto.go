package venue_transport_http

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_http_types "github.com/emount4/concert_reviews/internal/core/transport/http/types"
)

// --- Requests ---

type CreateVenueRequest struct {
	CityID      int               `json:"city_id" validate:"required,min=1"`
	Name        string            `json:"name" validate:"required,min=2,max=255"`
	Address     string            `json:"address" validate:"omitempty,max=500"`
	Capacity    *int              `json:"capacity" validate:"omitempty,min=1,max=500000"`
	PhotoKey    string            `json:"photo_key" validate:"omitempty,max=2048"`
	Description string            `json:"description" validate:"omitempty,max=2000"`
	SocialLinks map[string]string `json:"social_links"`
}

type UpdateVenueRequest struct {
	CityID      core_http_types.Nullable[int]           `json:"city_id"`
	Name        core_http_types.Nullable[string]        `json:"name"`
	Address     core_http_types.Nullable[string]        `json:"address"`
	Capacity    core_http_types.Nullable[int]           `json:"capacity"`
	PhotoKey    core_http_types.Nullable[string]        `json:"photo_key"`
	Description core_http_types.Nullable[string]        `json:"description"`
	SocialLinks core_http_types.NullableMapStringString `json:"social_links"`
	Status      core_http_types.Nullable[string]        `json:"status"`
}

// --- Responses (DTO для выхода) ---

type CityBriefResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type VenueStatsResponse struct {
	ReviewsCount   int    `json:"reviews_count"`
	SumRatingTotal int64  `json:"sum_rating_total"`
	ConcertsCount  int    `json:"concerts_count"`
	FavoritesCount int    `json:"favorites_count"`
	UpdatedAt      string `json:"updated_at"`
}

type VenueResponse struct {
	ID          int                `json:"id"`
	City        CityBriefResponse  `json:"city"` // <-- вложенный объект
	Name        string             `json:"name"`
	Address     string             `json:"address,omitempty"`
	Capacity    *int               `json:"capacity,omitempty"`
	PhotoURL    string             `json:"photo_url,omitempty"`
	Description string             `json:"description,omitempty"`
	SocialLinks map[string]string  `json:"social_links,omitempty"`
	Stats       VenueStatsResponse `json:"stats,omitempty"`
	Status      string             `json:"status"`
	CreatedAt   string             `json:"created_at"`
}

type VenueResponseAdmin struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	City        CityBriefResponse  `json:"city"`
	Address     string             `json:"address,omitempty"`
	Capacity    *int               `json:"capacity,omitempty"`
	PhotoURL    string             `json:"photo_url,omitempty"`
	Description string             `json:"description,omitempty"`
	SocialLinks map[string]string  `json:"social_links,omitempty"`
	Stats       VenueStatsResponse `json:"stats,omitempty"`
	Status      string             `json:"status"`
	CreatedAt   string             `json:"created_at"`
	DeletedAt   string             `json:"deleted_at,omitempty"`
}

type ListVenuesResponse struct {
	Items []VenueResponse `json:"items"`
}

type ListVenuesAdminResponse struct {
	Items []VenueResponseAdmin `json:"items"`
}

func MapCreateVenueDTOToDomain(dto CreateVenueRequest) domain.Venue {
	var address *string
	if dto.Address != "" {
		address = &dto.Address
	}

	var photo *string
	if dto.PhotoKey != "" {
		photo = &dto.PhotoKey
	}

	var desc *string
	if dto.Description != "" {
		desc = &dto.Description
	}

	return domain.Venue{
		CityID:      dto.CityID,
		Name:        dto.Name,
		Address:     address,
		Capacity:    dto.Capacity,
		PhotoURL:    photo,
		Description: desc,
		SocialLinks: dto.SocialLinks,
		Status:      domain.StatusActive,
	}
}

func MapDomainToVenueResponse(v domain.Venue) VenueResponse {
	resp := VenueResponse{
		ID: v.VenueID,

		Name:      v.Name,
		Status:    string(v.Status),
		CreatedAt: v.CreatedAt.Format(time.RFC3339),
		Capacity:  v.Capacity, // *int -> *int
	}

	// Безопасно разыменовываем указатели

	if v.City != nil {
		resp.City = CityBriefResponse{
			ID:   v.City.CityID,
			Name: v.City.Name,
			Slug: v.City.Slug,
		}
	}
	if v.Address != nil {
		resp.Address = *v.Address
	}

	if v.PhotoURL != nil {
		resp.PhotoURL = *v.PhotoURL
	}

	if v.Description != nil {
		resp.Description = *v.Description
	}

	if v.SocialLinks != nil {
		resp.SocialLinks = v.SocialLinks
	}

	if v.Stats != nil {
		resp.Stats = VenueStatsResponse{
			ReviewsCount:   v.Stats.ReviewsCount,
			SumRatingTotal: v.Stats.SumRatingTotal,
			ConcertsCount:  v.Stats.ConcertsCount,
			FavoritesCount: v.Stats.FavoritesCount,
			UpdatedAt:      v.Stats.UpdatedAt.Format(time.RFC3339),
		}
	}

	return resp
}

func MapDomainToVenueAdminResponse(v domain.Venue) VenueResponseAdmin {
	resp := VenueResponseAdmin{
		ID:        v.VenueID,
		Name:      v.Name,
		Status:    string(v.Status),
		CreatedAt: v.CreatedAt.Format(time.RFC3339),
		Capacity:  v.Capacity,
	}

	if v.City != nil {
		resp.City = CityBriefResponse{
			ID:   v.City.CityID,
			Name: v.City.Name,
			Slug: v.City.Slug,
		}
	}

	if v.DeletedAt != nil {
		resp.DeletedAt = v.DeletedAt.Format(time.RFC3339)
	}
	// Безопасно разыменовываем указатели
	if v.Address != nil {
		resp.Address = *v.Address
	}

	if v.PhotoURL != nil {
		resp.PhotoURL = *v.PhotoURL
	}

	if v.Description != nil {
		resp.Description = *v.Description
	}

	if v.SocialLinks != nil {
		resp.SocialLinks = v.SocialLinks
	}

	if v.Stats != nil {
		resp.Stats = VenueStatsResponse{
			ReviewsCount:   v.Stats.ReviewsCount,
			SumRatingTotal: v.Stats.SumRatingTotal,
			ConcertsCount:  v.Stats.ConcertsCount,
			FavoritesCount: v.Stats.FavoritesCount,
			UpdatedAt:      v.Stats.UpdatedAt.Format(time.RFC3339),
		}
	}

	return resp
}

// MapDomainListToVenueResponse преобразует срез доменных моделей в DTO ответа
func MapDomainListToVenueResponse(venues []domain.Venue) ListVenuesResponse {
	items := make([]VenueResponse, len(venues))
	for i, v := range venues {
		items[i] = MapDomainToVenueResponse(v)
	}
	return ListVenuesResponse{Items: items}
}
func MapDomainListToVenueAdminResponse(venues []domain.Venue) ListVenuesAdminResponse {
	items := make([]VenueResponseAdmin, len(venues))
	for i, v := range venues {
		items[i] = MapDomainToVenueAdminResponse(v)
	}
	return ListVenuesAdminResponse{Items: items}
}

// MapPatchVenueReqToDomain преобразует запрос обновления в доменный патч
func MapPatchVenueReqToDomain(req UpdateVenueRequest) domain.VenuePatch {
	return domain.VenuePatch{
		CityID:      req.CityID.ToDomain(),
		Name:        req.Name.ToDomain(),
		Address:     req.Address.ToDomain(),
		Capacity:    req.Capacity.ToDomain(),
		PhotoKey:    req.PhotoKey.ToDomain(),
		Description: req.Description.ToDomain(),
		SocialLinks: req.SocialLinks,
		Status:      req.Status.ToDomain(),
	}
}
