package config

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

type (
	Config struct {
		Component string `env:"COMPONENT,default=translation-api"`
		Cache     struct {
			Enabled bool `env:"CACHE_ENABLED,default=true"`
		}
		Log struct {
			Level       string `env:"LOG_LEVEL,default=debug"`
			PrettyPrint bool   `env:"LOG_PRETTY_PRINT,default=true"`
		}
		DB struct {
			URL string `env:"DB_URL,default=postgres://postgres:root@127.0.0.1:5432/postgres?sslmode=disable"`
		}
	}
)

func New(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		fmt.Printf("error loading configuration: %s", err.Error())
		return nil, err
	}
	return &cfg, nil
}
