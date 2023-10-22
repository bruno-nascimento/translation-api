package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bruno-nascimento/translation-api/internal/entrypoint/http/api"
)

type TranslationAPI struct {
}

func NewTranslationAPI() *TranslationAPI {
	return &TranslationAPI{}
}

func (t TranslationAPI) FindTranslation(c *fiber.Ctx, params api.FindTranslationParams) error {
	return c.JSON(api.Translation{Result: "teste"})
}

func (t TranslationAPI) AddTranslation(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusCreated)
}
