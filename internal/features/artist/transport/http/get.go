package artist_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *ArtistHTTPHandler) GetArtists(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetByQueryParam(r)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'limit/offset' query param")
		return
	}

	search := core_http_request.GetStringQueryParam(r, "search")

	artistsDomains, err := h.artistService.GetArtists(ctx, *search, limit, offset)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get artists")
		return
	}

	artistsDTO := MapDomainListToResponse(artistsDomains)

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
