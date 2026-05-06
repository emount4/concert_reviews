package artist_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *ArtistHTTPHandler) Create(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var req CreateArtistRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r,
		&req,
	); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode artist post request")
		return
	}

	artist := MapCreateDTOToDomain(req)

	createdArtist, err := h.artistService.CreateArtist(ctx, artist)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create artist")
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	response := MapDomainToResponse(createdArtist)

	responseHandler.JSONResponse(response, http.StatusCreated)
}
