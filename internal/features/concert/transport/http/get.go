package concert_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func defaultConcertDirection(sort string) string {
	switch sort {
	case "title":
		return "ASC"
	default:
		return "DESC"
	}
}

func allowedConcertSort(sort string) bool {
	switch sort {
	case "date", "rating", "reviews", "title", "p1", "p2", "p3", "p4", "p5":
		return true
	default:
		return false
	}
}

// GetConcerts — Получение публичного списка концертов с фильтрами
func (h *ConcertHTTPHandler) GetConcerts(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit/offset")
		return
	}

	sort, direction, directionProvided := core_http_request.GetSortParamsDetailed(r, "date", "DESC")
	if !allowedConcertSort(sort) {
		sort = "date"
	}
	if !directionProvided {
		direction = defaultConcertDirection(sort)
	}

	cityID, _ := core_http_request.GetIntQueryParam(r, "city_id")
	artistID, _ := core_http_request.GetIntQueryParam(r, "artist_id")
	search := r.URL.Query().Get("search")

	concerts, total, err := h.concertService.GetConcerts(ctx, cityID, artistID, search, sort, direction, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get concerts")
		return
	}

	response := MapDomainListToConcertResponse(concerts)
	response.PageCount = core_http_request.GetPageCount(total, limit)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func (h *ConcertHTTPHandler) GetConcert(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	idStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "missing concert id")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid concert uuid format")
		return
	}

	concert, err := h.concertService.GetConcertByID(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get concert")
		return
	}

	responseHandler.JSONResponse(MapDomainToConcertResponse(concert), http.StatusOK)
}

func (h *ConcertHTTPHandler) GetConcertsAdmin(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit/offset")
		return
	}

	sort, direction, directionProvided := core_http_request.GetSortParamsDetailed(r, "date", "DESC")
	if !allowedConcertSort(sort) {
		sort = "date"
	}
	if !directionProvided {
		direction = defaultConcertDirection(sort)
	}
	cityID, _ := core_http_request.GetIntQueryParam(r, "city_id")
	artistID, _ := core_http_request.GetIntQueryParam(r, "artist_id")
	search := r.URL.Query().Get("search")
	includeDeleted := false
	if value := core_http_request.GetBoolQueryParam(r, "include_deleted"); value != nil {
		includeDeleted = *value
	}

	concerts, total, err := h.concertService.GetConcertsAdmin(ctx, cityID, artistID, search, sort, direction, limit, offset, includeDeleted)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to fetch admin concerts list")
		return
	}

	response := MapDomainListToConcertAdminResponse(concerts)
	response.PageCount = core_http_request.GetPageCount(total, limit)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func (h *ConcertHTTPHandler) GetSuggestionsAdmin(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, _ := core_http_request.GetLimitOffsetByQueryParam(r)
	status := ""
	if value := core_http_request.GetStringQueryParam(r, "status"); value != nil {
		status = *value
	}

	suggestions, err := h.concertService.GetSuggestionsAdmin(ctx, limit, offset, status)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to fetch suggestions")
		return
	}

	responseHandler.JSONResponse(MapDomainListToSuggestionResponse(suggestions), http.StatusOK)
}

func (h *ConcertHTTPHandler) GetSuggestionByIDAdmin(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	idStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "missing suggestion id")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid suggestion uuid")
		return
	}

	suggestion, err := h.concertService.GetSuggestionByIDAdmin(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get suggestion")
		return
	}

	responseHandler.JSONResponse(MapDomainToSuggestionResponse(suggestion), http.StatusOK)
}
