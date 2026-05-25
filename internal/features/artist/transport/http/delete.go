package artist_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *ArtistHTTPHandler) DeleteArtistHard(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid artist id in url")
		return
	}

	if err := h.artistService.DeleteArtistHard(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, "failed to hard delete artist")
		return
	}

	h.logAdminAction(ctx, "artist_hard_deleted", id, nil)

	rw.WriteHeader(http.StatusNoContent)
}

func (h *ArtistHTTPHandler) DeleteArtistSoft(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid artist id in url")
		return
	}

	if err := h.artistService.DeleteArtistSoft(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, "failed to soft delete artist")
		return
	}

	h.logAdminAction(ctx, "artist_soft_deleted", id, nil)

	rw.WriteHeader(http.StatusNoContent)
}

func (h *ArtistHTTPHandler) RestoreArtist(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid artist id in url")
		return
	}

	restoredArtist, err := h.artistService.RestoreArtist(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to restore artist")
		return
	}

	h.logAdminAction(ctx, "artist_restored", id, map[string]any{
		"name": restoredArtist.Name,
	})

	rw.Header().Set("Content-Type", "application/json")
	response := MapDomainToResponse(restoredArtist)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func (h *ArtistHTTPHandler) GetArtistsAdmin(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	// Пагинация
	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get pagination params")
		return
	}

	sort, direction, directionProvided := core_http_request.GetSortParamsDetailed(r, "name", "ASC")
	if !allowedArtistSort(sort) {
		sort = "name"
	}
	if !directionProvided {
		direction = defaultArtistDirection(sort)
	}

	// Фильтры
	includeDeletedPtr := core_http_request.GetBoolQueryParam(r, "include_deleted")
	statusPtr := core_http_request.GetStringQueryParam(r, "status")
	searchPtr := core_http_request.GetStringQueryParam(r, "search")
	hasReviews := core_http_request.GetBoolQueryParam(r, "has_reviews")

	includeDeleted := false
	if includeDeletedPtr != nil {
		includeDeleted = *includeDeletedPtr
	}

	status := ""
	if statusPtr != nil {
		status = *statusPtr
	}

	search := ""
	if searchPtr != nil {
		search = *searchPtr
	}

	// Вызов сервиса
	artists, total, err := h.artistService.GetArtistsAdmin(
		ctx,
		search,
		sort,
		direction,
		hasReviews,
		limit,
		offset,
		includeDeleted,
		status,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to fetch artists for admin")
		return
	}

	// Маппинг → ответ
	// Используем тот же MapDomainListToResponse, либо создаём отдельный для админки
	// Если нужно показывать deleted_at — расширьте ArtistResponse или сделайте Admin-версию
	response := MapDomainListToAdminResponse(artists)
	response.PageCount = core_http_request.GetPageCount(total, limit)

	rw.Header().Set("Content-Type", "application/json")
	responseHandler.JSONResponse(response, http.StatusOK)
}
