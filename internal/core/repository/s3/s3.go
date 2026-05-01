package core_s3

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Storage struct {
	client       *minio.Client
	bucketName   string
	minSize      int64
	maxSize      int64
	contentTypes []string
}

func NewS3Storage(cfg Config) (*S3Storage, error) {
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

func (s *S3Storage) GetUploadURL(
	ctx context.Context,
	objectName string,
	expires time.Duration,
) (string, error) {
	presignedURL, err := s.client.PresignedPutObject(ctx, s.bucketName, objectName, expires)
	if err != nil {
		return "", fmt.Errorf("generate presigned put url: %w", err)
	}
	return presignedURL.String(), nil
}

func (s *S3Storage) GetUploadForm(
	ctx context.Context,
	objectName string,
	expires time.Duration,
) (string, map[string]string, error) {
	policy := minio.NewPostPolicy()
	if err := policy.SetBucket(s.bucketName); err != nil {
		return "", nil, fmt.Errorf("set bucket policy: %w", err)
	}
	if err := policy.SetKey(objectName); err != nil {
		return "", nil, fmt.Errorf("set object key policy: %w", err)
	}
	if err := policy.SetExpires(time.Now().Add(expires)); err != nil {
		return "", nil, fmt.Errorf("set policy expiry: %w", err)
	}
	if err := policy.SetContentLengthRange(s.minSize, s.maxSize); err != nil {
		return "", nil, fmt.Errorf("set size policy: %w", err)
	}
	if len(s.contentTypes) > 0 {
		contentTypeStartsWith := strings.Join(s.contentTypes, ",")
		if err := policy.SetContentTypeStartsWith(contentTypeStartsWith); err != nil {
			return "", nil, fmt.Errorf("set content type policy: %w", err)
		}
	}

	url, formData, err := s.client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return "", nil, fmt.Errorf("generate presigned post policy: %w", err)
	}

	result := make(map[string]string, len(formData)+1)
	for k, v := range formData {
		result[k] = v
	}
	result["url"] = url.String()

	return url.String(), result, nil
}

func parseAllowedTypes(raw string) []string {
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}
