package main

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/bruno-nascimento/translation-api/internal/config"
	"github.com/bruno-nascimento/translation-api/internal/db"
	"github.com/bruno-nascimento/translation-api/internal/entrypoint/http"
	"github.com/bruno-nascimento/translation-api/internal/logger"
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

	translationAPI := http.NewTranslationAPI()

	s := http.NewServer(translationAPI)

	if err = s.Listen(net.JoinHostPort("0.0.0.0", "8080")); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
