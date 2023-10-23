package http

import (
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	middleware "github.com/oapi-codegen/fiber-middleware"
	"github.com/rs/zerolog/log"

	"github.com/bruno-nascimento/translation-api/internal/config"
	"github.com/bruno-nascimento/translation-api/internal/entrypoint/http/api"
)

func NewServer(cfg *config.Config, translation *TranslationAPI) *fiber.App {

	swggr, err := api.GetSwagger()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading swagger spec")
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
		ExposeHeaders:    "*",
	}))

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/docs")
	})

	app.Use(swagger.New(swagger.Config{
		FilePath: cfg.OpenApiSpec.Path,
	}))

	app.Use(logger.New())

	api.RegisterHandlers(app, translation)

	app.Use(middleware.OapiRequestValidator(swggr))

	return app
}
