package config

import (
	"context"
	"github.com/creasty/defaults"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Runtime `env:",prefix=RUNTIME_"`
	Proxy   `env:",prefix=PROXY_"`
}

func New(ctx context.Context, cfg *Config) error {
	godotenv.Load()

	err := envconfig.Process(ctx, cfg)
	if err != nil {
		return err
	}

	err = defaults.Set(cfg)
	if err != nil {
		return err
	}

	return nil
}
