package core_redis

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr        string        `envconfig:"ADDR" required:"true"`
	Password    string        `envconfig:"PASSWORD"`
	User        string        `envconfig:"USER"`
	DB          int           `envconfig:"DB" required:"true"`
	MaxRetries  int           `envconfig:"MAX_RETRIES" required:"true"`
	DialTimeout time.Duration `envconfig:"DIAL_TIMEOUT" required:"true"`
	Timeout     time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("REDIS", &config); err != nil {
		return Config{}, fmt.Errorf("redis config: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Redis config: %w", err)
		panic(err)
	}
	return config
}
