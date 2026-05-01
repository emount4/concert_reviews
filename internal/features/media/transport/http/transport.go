package media_transport_http

import (
	"context"
	"net/http"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
)

// Service теперь работает только с Доменными моделями
type Service interface {
	PrepareBatchUpload(ctx context.Context, files []core_models.MediaUploadParams) ([]core_models.MediaUploadTicket, error)
}

type MediaHTTPHandler struct {
	mediaService Service
}

func NewMediaHTTPHandler(mediaService Service) *MediaHTTPHandler {
	return &MediaHTTPHandler{mediaService: mediaService}
}

func (h *MediaHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/media/presign",
			Handler: http.HandlerFunc(h.GetPresignedURLs),
		},
	}
}
