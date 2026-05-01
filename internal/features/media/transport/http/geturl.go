package media_transport_http

import (
	"net/http"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_http_request "github.com/emount4/concert_reviews/internal/core/transport/http/request"
	core_http_response "github.com/emount4/concert_reviews/internal/core/transport/http/response"
)

func (h *MediaHTTPHandler) GetPresignedURLs(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var req BatchUploadRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	domainParams := dtoToDomain(req)

	// Вызываем сервис, передавая доменные модели
	tickets, err := h.mediaService.PrepareBatchUpload(ctx, domainParams)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to prepare upload links")
		return
	}

	respItems := domainToDto(tickets)
	responseHandler.JSONResponse(BatchUploadResponse{Items: respItems}, http.StatusOK)
}

func dtoToDomain(dto BatchUploadRequest) []core_models.MediaUploadParams {
	domainParams := make([]core_models.MediaUploadParams, len(dto.Files))
	for i, f := range dto.Files {
		domainParams[i] = core_models.MediaUploadParams{
			FileName: f.FileName,
			FileSize: f.FileSize,
		}
	}
	return domainParams
}

func domainToDto(tickets []core_models.MediaUploadTicket) []UploadItemDTO {
	respItems := make([]UploadItemDTO, len(tickets))
	for i, t := range tickets {
		respItems[i] = UploadItemDTO{
			FileKey:    t.FileKey,
			UploadURL:  t.UploadURL,
			UploadForm: t.UploadForm,
		}
	}
	return respItems
}
