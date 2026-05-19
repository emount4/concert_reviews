package favorites_transport_http

import (
	"net/http"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *FavoritesHTTPHandler) AddFavorite(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	userID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	var req FavoriteRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		res.ErrorResponse(err, "invalid favorite data")
		return
	}

	favorite, err := h.favoritesService.AddFavorite(ctx, userID, req.TargetType, req.TargetID)
	if err != nil {
		res.ErrorResponse(err, "failed to add favorite")
		return
	}

	res.JSONResponse(MapDomainToResponse(favorite), http.StatusCreated)
}

func (h *FavoritesHTTPHandler) DeleteFavorite(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	userID, err := core_http_middleware.GetUserID(ctx)
	if err != nil {
		res.ErrorResponse(err, "unauthorized")
		return
	}

	targetType, err := core_http_request.GetStringPathValue(r, "target_type")
	if err != nil {
		res.ErrorResponse(err, "missing favorite target type")
		return
	}
	targetID, err := core_http_request.GetStringPathValue(r, "target_id")
	if err != nil {
		res.ErrorResponse(err, "missing favorite target id")
		return
	}

	if err := h.favoritesService.DeleteFavorite(ctx, userID, targetType, targetID); err != nil {
		res.ErrorResponse(err, "failed to delete favorite")
		return
	}

	res.JSONResponse(nil, http.StatusNoContent)
}

func (h *FavoritesHTTPHandler) GetFavoritesByUsername(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)

	username, err := core_http_request.GetStringPathValue(r, "username")
	if err != nil {
		res.ErrorResponse(err, "missing username")
		return
	}

	targetType := core_http_request.GetStringQueryParam(r, "type")
	favorites, err := h.favoritesService.GetFavoritesByUsername(ctx, username, targetType)
	if err != nil {
		res.ErrorResponse(err, "failed to get favorites")
		return
	}

	res.JSONResponse(MapDomainListToResponse(favorites), http.StatusOK)
}
