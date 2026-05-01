package media_service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
	"github.com/google/uuid"
)

// S3Provider — интерфейс, который должен реализовать наш пакет инфраструктуры (core_s3)
type S3Provider interface {
	GetUploadURL(ctx context.Context, objectName string, expires time.Duration) (string, error)
	GetUploadForm(ctx context.Context, objectName string, expires time.Duration) (string, map[string]string, error)
}

type MediaService struct {
	s3            S3Provider
	allowedExt    map[string]bool
	minUploadSize int64
	maxUploadSize int64
}

func NewMediaService(s3 S3Provider, allowedExtensions map[string]bool, minUploadSize, maxUploadSize int64) *MediaService {
	return &MediaService{
		s3:            s3,
		allowedExt:    allowedExtensions,
		minUploadSize: minUploadSize,
		maxUploadSize: maxUploadSize,
	}
}

func (s *MediaService) PrepareBatchUpload(
	ctx context.Context,
	files []core_models.MediaUploadParams,
) ([]core_models.MediaUploadTicket, error) {

	tickets := make([]core_models.MediaUploadTicket, 0, len(files))

	for _, f := range files {
		// 1. Извлекаем и нормализуем расширение
		ext := strings.ToLower(filepath.Ext(f.FileName))
		if ext == "" {
			ext = ".jpg" // дефолт, если расширения нет
		}

		// 2. Валидация расширения
		if !s.allowedExt[ext] {
			return nil, fmt.Errorf("file type %s is not allowed", ext)
		}

		if f.FileSize < s.minUploadSize || f.FileSize > s.maxUploadSize {
			return nil, fmt.Errorf("file size %d is out of allowed range", f.FileSize)
		}

		// 3. Генерируем уникальный ключ объекта для S3 (uuid + расширение)
		// Сохраняем в папку reviews/
		fileKey := fmt.Sprintf("reviews/%s%s", uuid.New().String(), ext)

		// 4. Запрашиваем Pre-signed POST policy у S3 провайдера (время жизни 15 минут)
		uploadURL, formData, err := s.s3.GetUploadForm(ctx, fileKey, 15*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("failed to get presigned post policy for %s: %w", f.FileName, err)
		}

		// 5. Добавляем "тикет" в результат
		tickets = append(tickets, core_models.MediaUploadTicket{
			FileKey:    fileKey,
			UploadURL:  uploadURL,
			UploadForm: formData,
		})
	}

	return tickets, nil
}
