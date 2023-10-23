package main

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/bruno-nascimento/translation-api/internal/config"
	"github.com/bruno-nascimento/translation-api/internal/db"
	"github.com/bruno-nascimento/translation-api/internal/entrypoint/http"
	"github.com/bruno-nascimento/translation-api/internal/logger"
	"github.com/bruno-nascimento/translation-api/internal/repository"
	"github.com/bruno-nascimento/translation-api/internal/service"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	logger.SetupLogger(cfg)

	err = db.RunMigrations(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	dbConn, err := db.Connect(cfg)
	if err != nil {
		return
	}
	defer dbConn.Close(ctx)

	repo, err := repository.NewRepository(cfg, dbConn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create repository")
	}

	translationAPI := http.NewTranslationAPI(service.NewTranslation(cfg, repo))

	s := http.NewServer(cfg, translationAPI)
	defer s.Shutdown()

	if err = s.Listen(net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}

}
