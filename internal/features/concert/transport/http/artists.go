package concert_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
	"github.com/google/uuid"
)

// AddConcertArtist handles POST /admin/concerts/{id}/artists
func (h *ConcertHTTPHandler) AddConcertArtist(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	// Get concert ID from path
	concertIDStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot get concert id from path")
		return
	}

	concertID, err := uuid.Parse(concertIDStr)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid concert uuid format")
		return
	}

	// Decode request body
	var req AddConcertArtistRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate add artist request")
		return
	}

	// Add artist to concert
	if err := h.concertService.AddArtistToConcert(ctx, concertID, req.ArtistID, req.IsMain); err != nil {
		responseHandler.ErrorResponse(err, "failed to add artist to concert")
		return
	}

	h.logAdminAction(ctx, "concert_artist_added", concertID, map[string]any{
		"artist_id": req.ArtistID,
		"is_main":   req.IsMain,
	})

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}

// RemoveConcertArtist handles DELETE /admin/concerts/{id}/artists/{artist_id}
func (h *ConcertHTTPHandler) RemoveConcertArtist(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	// Get concert ID from path
	concertIDStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot get concert id from path")
		return
	}

	concertID, err := uuid.Parse(concertIDStr)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid concert uuid format")
		return
	}

	// Get artist ID from path
	artistID, err := core_http_request.GetIntPathValue(r, "artist_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot get artist id from path")
		return
	}

	// Remove artist from concert
	if err := h.concertService.RemoveArtistFromConcert(ctx, concertID, artistID); err != nil {
		responseHandler.ErrorResponse(err, "failed to remove artist from concert")
		return
	}

	h.logAdminAction(ctx, "concert_artist_removed", concertID, map[string]any{
		"artist_id": artistID,
	})

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}

// UpdateConcertArtistIsMain handles PATCH /admin/concerts/{id}/artists/{artist_id}
func (h *ConcertHTTPHandler) UpdateConcertArtistIsMain(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	// Get concert ID from path
	concertIDStr, err := core_http_request.GetStringPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot get concert id from path")
		return
	}

	concertID, err := uuid.Parse(concertIDStr)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid concert uuid format")
		return
	}

	// Get artist ID from path
	artistID, err := core_http_request.GetIntPathValue(r, "artist_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "cannot get artist id from path")
		return
	}

	// Decode request body
	var req UpdateConcertArtistRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate update artist request")
		return
	}

	// Update artist is_main flag
	if err := h.concertService.UpdateArtistMainStatus(ctx, concertID, artistID, req.IsMain); err != nil {
		responseHandler.ErrorResponse(err, "failed to update artist main status")
		return
	}

	h.logAdminAction(ctx, "concert_artist_main_updated", concertID, map[string]any{
		"artist_id": artistID,
		"is_main":   req.IsMain,
	})

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}
