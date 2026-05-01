package core_s3

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Endpoint     string `envconfig:"ENDPOINT" required:"true"`
	AccessKey    string `envconfig:"ACCESS_KEY" required:"true"`
	SecretKey    string `envconfig:"SECRET_KEY" required:"true"`
	BucketName   string `envconfig:"BUCKET_NAME" required:"true"`
	UseSSL       bool   `envconfig:"USE_SSL" required:"true"`
	Region       string `envconfig:"REGION" default:"us-east-1"`
	UploadMinMB  int64  `envconfig:"UPLOAD_MIN_MB" default:"0"`
	UploadMaxMB  int64  `envconfig:"UPLOAD_MAX_MB" default:"50"`
	AllowedTypes string `envconfig:"ALLOWED_TYPES" default:"image/jpeg,image/png,image/webp,image/gif"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process(
		"S3",
		&config,
	); err != nil {
		return Config{}, fmt.Errorf("s3 config: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get s3 config: %w", err)
		panic(err)
	}
	return config
}
