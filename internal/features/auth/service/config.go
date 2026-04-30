package auth_service

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PasswordSalt   string        `envconfig:"PASSWORD_SALT" required:"true"`
	JWTSigningKey  string        `envconfig:"JWT_SIGNING_KEY" required:"true"`
	AccessTokenTTL time.Duration `envconfig:"ACCESS_TOKEN_TTL" required:"true"`
	BotToken       string        `envconfig:"BOT_TOKEN" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("AUTH", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Auth config: %w", err)
		panic(err)
	}
	return config
}
