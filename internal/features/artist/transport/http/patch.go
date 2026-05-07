package artist_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *ArtistHTTPHandler) PatchArtist(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot get artist id path value")
		return
	}

	var req UpdateArtistRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r,
		&req,
	); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode patch artist request")
		return
	}

	artistPatch := MapPatchReqToDomain(req)

	domain, err := h.artistService.PatchArtist(ctx, id, artistPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot patch artist")
		return
	}

	resp := MapDomainToResponse(domain)
	responseHandler.JSONResponse(resp, http.StatusOK)
}
