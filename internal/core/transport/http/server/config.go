package core_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr                 string        `envconfig:"ADDR"                  required:"true"`
	ShutdownTimeout      time.Duration `envconfig:"SHUTDOWN_TIMEOUT"      default:"30s"`
	CORSAllowedOrigins   string        `envconfig:"CORS_ALLOWED_ORIGINS"  default:""`
	CORSAllowCredentials bool          `envconfig:"CORS_ALLOW_CREDENTIALS" default:"true"`
	CORSMaxAgeSeconds    int           `envconfig:"CORS_MAX_AGE_SECONDS"  default:"600"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("HTTP", &config); err != nil {
		return Config{}, fmt.Errorf("http server config process err: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()

	if err != nil {
		err = fmt.Errorf("get HTTP process config: %w", err)
		panic(err)
	}
	return config
}
