package config

import (
	"context"
	"fmt"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type (
	Config struct {
		Component string `env:"COMPONENT,default=translation-api"`
		Cache     struct {
			Enabled   bool          `env:"CACHE_ENABLED,default=true"`
			RedisAddr string        `env:"CACHE_REDIS_ADDR,default=localhost:6379"`
			TTL       time.Duration `env:"CACHE_TTL,default=30s"`
		}
		Log struct {
			Level       string `env:"LOG_LEVEL,default=debug"`
			PrettyPrint bool   `env:"LOG_PRETTY_PRINT,default=true"`
		}
		DB struct {
			URL string `env:"DB_URL,default=postgres://postgres:root@127.0.0.1:5432/postgres?sslmode=disable"`
		}
		Suggestions struct {
			Enabled bool `env:"SUGGESTIONS_ENABLED,default=true"`
			Limit   int  `env:"SUGGESTIONS_LIMIT,default=5"`
		}
		OpenApiSpec struct {
			Path string `env:"OPENAPI_SPEC_PATH,default=/Users/bruno.nascimento/dev/code/tmp/translation-api/openapi/translation-api.yml"`
		}
		Server struct {
			Port string `env:"SERVER_PORT,default=8080"`
			Host string `env:"SERVER_HOST,default=127.0.0.1"`
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
