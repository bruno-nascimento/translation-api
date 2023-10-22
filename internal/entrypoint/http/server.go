package http

import (
	"fmt"
	"os"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	middleware "github.com/oapi-codegen/fiber-middleware"

	"github.com/bruno-nascimento/translation-api/internal/entrypoint/http/api"
)

func NewServer(translation *TranslationAPI) *fiber.App {

	swggr, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

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

	cfg := swagger.Config{
		FilePath: "/Users/bruno.nascimento/dev/code/tmp/translation-api/openapi/translation-api.yml",
	}

	app.Use(swagger.New(cfg))

	app.Use(logger.New())

	api.RegisterHandlers(app, translation)

	app.Use(middleware.OapiRequestValidator(swggr))

	return app
}
