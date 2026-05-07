package core_s3

import (
	"context"
	"fmt"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type S3Storage struct {
	client       *minio.Client
	bucketName   string
	minSize      int64
	maxSize      int64
	contentTypes []string
}

func NewS3Storage(log *core_logger.Logger, cfg Config) (*S3Storage, error) {
	client, err := minio.New(
		cfg.Endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				cfg.AccessKey,
				cfg.SecretKey,
				"",
			),
			Secure: cfg.UseSSL,
			Region: cfg.Region,
		})
	if err != nil {
		log.Debug("S3 params \n",
			zap.String("endpoint", cfg.Endpoint))
		return nil, fmt.Errorf("init s3 client: %w", err)
	}

	storage := &S3Storage{
		client:       client,
		bucketName:   cfg.BucketName,
		minSize:      cfg.UploadMinMB * 1024 * 1024,
		maxSize:      cfg.UploadMaxMB * 1024 * 1024,
		contentTypes: parseAllowedTypes(cfg.AllowedTypes),
	}

	exists, err := client.BucketExists(context.Background(), cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("check bucket exists: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(context.Background(), cfg.BucketName, minio.MakeBucketOptions{Region: cfg.Region}); err != nil {
			return nil, fmt.Errorf("create bucket: %w", err)
		}
	}

	return storage, nil
}
