package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/bruno-nascimento/translation-api/internal/entrypoint/http/api"
	"github.com/bruno-nascimento/translation-api/internal/service"
)

type TranslationAPI struct {
	service service.Translation
}

func NewTranslationAPI(service service.Translation) *TranslationAPI {
	return &TranslationAPI{service: service}
}

func (t TranslationAPI) FindTranslation(c *fiber.Ctx, params api.FindTranslationParams) error {
	findParams, err := service.NewFindTranslationParamFromAPI(params)
	if err != nil {
		return sendError(c, http.StatusBadRequest, err.Error())
	}
	translations, err := t.service.FindTranslations(c.Context(), findParams)
	if err != nil {
		return sendError(c, http.StatusBadRequest, err.Error())
	}

	if translations.Result != nil {
		return c.JSON(translations.Result)
	}
	if translations.Suggestion != nil {
		return c.Status(fiber.StatusNotFound).JSON(translations.Suggestion)
	}
	return c.SendStatus(fiber.StatusNoContent)

}

func (t TranslationAPI) AddTranslation(c *fiber.Ctx) error {

	var newTranslation api.NewTranslation

	if err := c.BodyParser(&newTranslation); err != nil {
		return sendError(c, http.StatusBadRequest, "Invalid format for new translation params")
	}

	serviceParam, err := service.NewTranslationParamFromAPI(newTranslation)
	if err != nil {
		return sendError(c, http.StatusBadRequest, err.Error())
	}

	err = t.service.CreateTranslation(c.Context(), serviceParam)
	if err != nil {
		return sendError(c, http.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func sendError(c *fiber.Ctx, code int, message string) error {

	petErr := api.Error{
		Code:    int32(code),
		Message: message,
	}

	return c.Status(code).JSON(petErr)
}
