package core_ports

import (
	"context"
	"time"
)

type S3Provider interface {
	GetUploadURL(ctx context.Context, objectName string, expires time.Duration) (string, error)
	GetUploadForm(ctx context.Context, objectName string, expires time.Duration) (string, map[string]string, error)

	FileExists(ctx context.Context, objectName string) (int64, error)
	GetPublicURL(objectName string) string
	DeleteObject(ctx context.Context, objectName string) error
}
