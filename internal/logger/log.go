package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/bruno-nascimento/translation-api/internal/config"
)

func SetupLogger(cfg *config.Config) {
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.WarnLevel
	}

	zerolog.SetGlobalLevel(level)
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	log.Logger = log.Logger.With().Fields([]interface{}{"COMPONENT", cfg.Component}).Logger()
	log.Logger.Info().Msg("Log configuration: done")
}
