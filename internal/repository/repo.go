package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bruno-nascimento/translation-api/internal/config"
	"github.com/bruno-nascimento/translation-api/internal/redis"
)

func NewRepository(cfg *config.Config, db DBTX) (Querier, error) {
	repo := New(db)
	if cfg.Cache.Enabled {
		conn, err := redis.Conn(cfg)
		if err != nil {
			return nil, err
		}
		return NewCachedRepository(cfg, conn, repo), nil
	}
	return repo, nil
}

const UniqueViolationErrCode = "23505"

func IsUniqueViolation(err error) bool {
	dbError := &pgconn.PgError{}
	if errors.As(err, &dbError) {
		return dbError.Code == UniqueViolationErrCode
	}
	return false
}
