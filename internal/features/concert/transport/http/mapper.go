package concert_transport_http

import (
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/google/uuid"
)

func MapCreateConcertDTOToDomain(dto CreateConcertRequest) (domain.Concert, []domain.ConcertArtist) {
	var poster *string
	if dto.PosterKey != "" {
		poster = &dto.PosterKey
	}

	concert := domain.Concert{
		VenueID:    dto.VenueID,
		Title:      dto.Title,
		Date:       dto.Date,
		PosterKey:  poster, // Сохраняем KEY как есть
		IsVerified: true,   // Создано админом
	}

	artists := make([]domain.ConcertArtist, len(dto.Artists))
	for i, a := range dto.Artists {
		artists[i] = domain.ConcertArtist{
			ArtistID: a.ArtistID,
			IsMain:   a.IsMain,
		}
	}

	return concert, artists
}

func MapPatchConcertReqToDomain(req UpdateConcertRequest) domain.ConcertPatch {
	return domain.ConcertPatch{
		VenueID:    req.VenueID.ToDomain(),
		Title:      req.Title.ToDomain(),
		Date:       req.Date.ToDomain(),
		PosterKey:  req.PosterKey.ToDomain(),
		IsVerified: req.IsVerified.ToDomain(),
	}
}

func MapDomainToConcertResponse(c domain.Concert) ConcertResponse {
	resp := ConcertResponse{
		ID:         c.ConcertID.String(),
		Title:      c.Title,
		Date:       c.Date.Format(time.RFC3339),
		IsVerified: c.IsVerified,
		CreatedAt:  c.CreatedAt.Format(time.RFC3339),
	}

	// Прокидываем KEY из базы как URL (фронт сам подставит домен S3)
	if c.PosterKey != nil {
		resp.PosterURL = *c.PosterKey
	}

	if c.Venue != nil {
		resp.Venue = VenueBriefResponse{
			ID:   c.Venue.VenueID,
			Name: c.Venue.Name,
			City: c.Venue.City.Name,
		}
	}

	if len(c.Artists) > 0 {
		resp.Artists = make([]ArtistBriefResponse, len(c.Artists))
		for i, a := range c.Artists {
			resp.Artists[i] = ArtistBriefResponse{
				ID:     a.ArtistID,
				Name:   a.Name,
				IsMain: a.IsMain,
			}
		}
	}

	if c.Stats != nil {
		resp.Stats = &ConcertStatsResponse{
			ReviewsCount:   c.Stats.ReviewsCount,
			AvgRatingTotal: c.Stats.AvgRating(),
			AvgP1:          c.Stats.AvgByParam(1),
			AvgP2:          c.Stats.AvgByParam(2),
			AvgP3:          c.Stats.AvgByParam(3),
			AvgP4:          c.Stats.AvgByParam(4),
			AvgP5:          c.Stats.AvgByParam(5),
			UpdatedAt:      c.Stats.UpdatedAt.Format(time.RFC3339),
		}
	}

	return resp
}

func MapDomainListToConcertResponse(concerts []domain.Concert) ListConcertsResponse {
	items := make([]ConcertResponse, len(concerts))
	for i, c := range concerts {
		items[i] = MapDomainToConcertResponse(c)
	}
	return ListConcertsResponse{Items: items}
}

// --- Mappings: Suggestion Request -> Domain ---

// MapCreateSuggestionDTOToDomain принимает DTO и ID пользователя из контекста
func MapCreateSuggestionDTOToDomain(dto CreateSuggestionRequest, userID uuid.UUID) domain.ConcertSuggestion {
	var infoPtr *string
	if dto.Info != "" {
		infoPtr = &dto.Info
	}

	return domain.ConcertSuggestion{
		SuggestionID:  uuid.New(), // Генерируем новый ID, если репозиторий не делает этого сам
		UserID:        userID,
		RawArtistName: dto.ArtistName,
		RawVenueName:  dto.VenueName,
		ConcertDate:   dto.Date,
		Info:          infoPtr,
	}
}

// --- Mappings: Suggestion Domain -> Response ---

func MapDomainToSuggestionResponse(s domain.ConcertSuggestion) SuggestionResponse {
	return SuggestionResponse{
		ID:         s.SuggestionID.String(),
		UserID:     s.UserID.String(),
		ArtistName: s.RawArtistName,
		VenueName:  s.RawVenueName,
		Date:       s.ConcertDate,
		Info:       s.Info,
		CreatedAt:  s.CreatedAt.Format(time.RFC3339),
	}
}

func MapDomainListToSuggestionResponse(suggestions []domain.ConcertSuggestion) ListSuggestionsResponse {
	items := make([]SuggestionResponse, len(suggestions))
	for i, s := range suggestions {
		items[i] = MapDomainToSuggestionResponse(s)
	}
	return ListSuggestionsResponse{Items: items}
}

func MapDomainToConcertAdminResponse(c domain.Concert) ConcertResponseAdmin {
	resp := ConcertResponseAdmin{
		ID:         c.ConcertID.String(),
		VenueID:    c.VenueID,
		Title:      c.Title,
		Date:       c.Date.Format(time.RFC3339),
		IsVerified: c.IsVerified,
		CreatedAt:  c.CreatedAt.Format(time.RFC3339),
	}

	if c.PosterKey != nil {
		resp.PosterURL = *c.PosterKey
	}

	if c.CreatedByUserID != nil {
		uid := c.CreatedByUserID.String()
		resp.CreatedByUserID = &uid
	}

	if c.DeletedAt != nil {
		da := c.DeletedAt.Format(time.RFC3339)
		resp.DeletedAt = &da
	}

	// Маппинг вложенных структур (Venue, Artists, Stats)
	// используется та же логика, что и в обычном ответе
	if c.Venue != nil {
		resp.Venue = VenueBriefResponse{
			ID:   c.Venue.VenueID,
			Name: c.Venue.Name,
			City: c.Venue.City.Name,
		}
	}

	if len(c.Artists) > 0 {
		resp.Artists = make([]ArtistBriefResponse, len(c.Artists))
		for i, a := range c.Artists {
			resp.Artists[i] = ArtistBriefResponse{
				ID:     a.ArtistID,
				Name:   a.Name,
				IsMain: a.IsMain,
			}
		}
	}

	if c.Stats != nil {
		resp.Stats = &ConcertStatsResponse{
			ReviewsCount:   c.Stats.ReviewsCount,
			AvgRatingTotal: c.Stats.AvgRating(),
			AvgP1:          c.Stats.AvgByParam(1),
			AvgP2:          c.Stats.AvgByParam(2),
			AvgP3:          c.Stats.AvgByParam(3),
			AvgP4:          c.Stats.AvgByParam(4),
			AvgP5:          c.Stats.AvgByParam(5),
			UpdatedAt:      c.Stats.UpdatedAt.Format(time.RFC3339),
		}
	}

	return resp
}

func MapDomainListToConcertAdminResponse(concerts []domain.Concert) ListConcertsAdminResponse {
	items := make([]ConcertResponseAdmin, len(concerts))
	for i, c := range concerts {
		items[i] = MapDomainToConcertAdminResponse(c)
	}
	return ListConcertsAdminResponse{Items: items}
}
