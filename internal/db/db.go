package db

import (
	"context"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/jackc/pgx/v5"

	"github.com/bruno-nascimento/translation-api/internal/config"
	translationapisql "github.com/bruno-nascimento/translation-api/sql"
)

func Connect(cfg *config.Config) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), cfg.DB.URL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func RunMigrations(cfg *config.Config) error {
	datasource, err := url.Parse(cfg.DB.URL)
	if err != nil {
		return err
	}
	migrations := dbmate.New(datasource)
	migrations.FS = translationapisql.FS
	migrations.MigrationsDir = []string{"migrations"}

	err = migrations.CreateAndMigrate()
	if err != nil {
		return err
	}
	return nil
}
