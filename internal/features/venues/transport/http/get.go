package venue_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func defaultVenueDirection(sort string) string {
	if sort == "name" {
		return "ASC"
	}
	return "DESC"
}

func allowedVenueSort(sort string) bool {
	switch sort {
	case "name", "rating", "capacity", "p1", "p2", "p3", "p4", "p5":
		return true
	default:
		return false
	}
}

// GetVenues — публичный список с пагинацией и фильтрами
// GET /venues?city_id=1&search=stadium&limit=20&offset=0
func (h *VenueHTTPHandler) GetVenues(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get pagination params")
		return
	}

	sort, direction, directionProvided := core_http_request.GetSortParamsDetailed(r, "name", "ASC")
	if !allowedVenueSort(sort) {
		sort = "name"
	}
	if !directionProvided {
		direction = defaultVenueDirection(sort)
	}

	cityID, err := core_http_request.GetIntQueryParam(r, "city_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid city query param")
		return
	}

	search := core_http_request.GetStringQueryParam(r, "search")
	capacityFrom, err := core_http_request.GetIntQueryParam(r, "capacity_from")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid capacity_from query param")
		return
	}
	capacityTo, err := core_http_request.GetIntQueryParam(r, "capacity_to")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid capacity_to query param")
		return
	}

	venues, total, err := h.venueService.GetVenues(ctx, cityID, *search, sort, direction, capacityFrom, capacityTo, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get venues")
		return
	}

	response := MapDomainListToVenueResponse(venues)
	response.PageCount = core_http_request.GetPageCount(total, limit)

	rw.Header().Set("Content-Type", "application/json")
	responseHandler.JSONResponse(response, http.StatusOK)
}

// GetVenue — получение одной площадки по ID
// GET /venues/{id}
func (h *VenueHTTPHandler) GetVenue(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get venue id path value")
		return
	}

	venue, err := h.venueService.GetVenueByID(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get venue")
		return
	}

	response := MapDomainToVenueResponse(venue)
	responseHandler.JSONResponse(response, http.StatusOK)
}

// GetVenuesAdmin — список для админки
func (h *VenueHTTPHandler) GetVenuesAdmin(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get pagination params")
		return
	}

	sort, direction, directionProvided := core_http_request.GetSortParamsDetailed(r, "name", "ASC")
	if !allowedVenueSort(sort) {
		sort = "name"
	}
	if !directionProvided {
		direction = defaultVenueDirection(sort)
	}

	cityID, err := core_http_request.GetIntQueryParam(r, "city_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid city query param")
		return
	}

	search := core_http_request.DerefOr(core_http_request.GetStringQueryParam(r, "search"), "")
	status := core_http_request.DerefOr(core_http_request.GetStringQueryParam(r, "status"), "")
	includeDeleted := core_http_request.DerefOr(core_http_request.GetBoolQueryParam(r, "include_deleted"), false)
	capacityFrom, err := core_http_request.GetIntQueryParam(r, "capacity_from")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid capacity_from query param")
		return
	}
	capacityTo, err := core_http_request.GetIntQueryParam(r, "capacity_to")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid capacity_to query param")
		return
	}

	venues, total, err := h.venueService.GetVenuesAdmin(
		ctx,
		cityID,
		search,
		sort,
		direction,
		capacityFrom,
		capacityTo,
		limit,
		offset,
		includeDeleted,
		status,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get venues for admin")
		return
	}

	response := MapDomainListToVenueAdminResponse(venues)
	response.PageCount = core_http_request.GetPageCount(total, limit)

	rw.Header().Set("Content-Type", "application/json")
	responseHandler.JSONResponse(response, http.StatusOK)
}
