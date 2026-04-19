package settings

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/middleware"
)

func RegisterAdminRoutes(router fiber.Router, service SettingsService) {
	h := NewSettingsHandler(service)
	
	settings := router.Group("/settings", middleware.Protected(), middleware.AdminOnly())
	settings.Get("/", h.GetGlobalSettings)
	settings.Post("/", h.UpdateGlobalSettings)
}
