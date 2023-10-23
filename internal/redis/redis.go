package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/bruno-nascimento/translation-api/internal/config"
)

func Conn(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{Addr: cfg.Cache.RedisAddr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
