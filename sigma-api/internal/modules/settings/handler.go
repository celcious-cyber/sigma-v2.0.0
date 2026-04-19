package settings

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/core/models"
)

type SettingsHandler struct {
	service SettingsService
}

func NewSettingsHandler(s SettingsService) *SettingsHandler {
	return &SettingsHandler{s}
}

func (h *SettingsHandler) GetGlobalSettings(c fiber.Ctx) error {
	settings, err := h.service.GetSettings()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(settings)
}

func (h *SettingsHandler) UpdateGlobalSettings(c fiber.Ctx) error {
	var s models.Setting
	if err := c.Bind().JSON(&s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.UpdateSettings(s); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Settings updated"})
}
