package artist_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func defaultArtistDirection(sort string) string {
	if sort == "name" {
		return "ASC"
	}
	return "DESC"
}

func allowedArtistSort(sort string) bool {
	switch sort {
	case "name", "rating", "reviews", "p1", "p2", "p3", "p4", "p5":
		return true
	default:
		return false
	}
}

func (h *ArtistHTTPHandler) GetArtists(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'limit/offset' query param")
		return
	}

	sort, direction, directionProvided := core_http_request.GetSortParamsDetailed(r, "name", "ASC")
	if !allowedArtistSort(sort) {
		sort = "name"
	}
	if !directionProvided {
		direction = defaultArtistDirection(sort)
	}

	search := core_http_request.GetStringQueryParam(r, "search")
	hasReviews := core_http_request.GetBoolQueryParam(r, "has_reviews")

	artistsDomains, total, err := h.artistService.GetArtists(ctx, *search, sort, direction, hasReviews, limit, offset)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get artists")
		return
	}

	artistsDTO := MapDomainListToResponse(artistsDomains)
	artistsDTO.PageCount = core_http_request.GetPageCount(total, limit)

	rw.Header().Set("Content-Type", "application/json")
	responseHandler.JSONResponse(artistsDTO, http.StatusOK)
}

func (h *ArtistHTTPHandler) GetArtist(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get artist id path value")
		return
	}

	artistDomain, err := h.artistService.GetArtistByID(ctx, id)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get artist")
		return
	}

	response := MapDomainToResponse(artistDomain)

	responseHandler.JSONResponse(response, http.StatusOK)
}
