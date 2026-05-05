package core_s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

// DeleteObject физически удаляет файл из облака.
// Нужен при отклонении фото модератором или смене аватарки.
func (s *S3Storage) DeleteObject(ctx context.Context, objectName string) error {
	opts := minio.RemoveObjectOptions{}
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, opts)
	if err != nil {
		return fmt.Errorf("s3 delete object: %w", err)
	}
	return nil
}

// GetPublicURL возвращает готовую ссылку для тега <img src="...">.
// Твой бэкенд будет использовать это в мапперах (из Domain в DTO).
func (s *S3Storage) GetPublicURL(objectName string) string {
	if objectName == "" {
		return ""
	}
	// Если бакет публичный (рекомендуется для рецензий),
	// ссылка формируется как: https://endpoint/bucket/object

	// ВАЖНО: В проде это может быть адрес CDN (например, https://cdn.mysite.ru/uuid.jpg)
	// Для простоты возвращаем прямую ссылку на облако
	return fmt.Sprintf("%s/%s/%s", s.client.EndpointURL().String(), s.bucketName, objectName)
}
