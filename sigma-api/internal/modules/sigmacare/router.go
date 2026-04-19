package sigmacare

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/middleware"
)

// RegisterAdminRoutes sets up routes for Backoffice UKS management
func RegisterAdminRoutes(router fiber.Router, service CareService) {
	h := NewCareHandler(service)
	
	care := router.Group("/care", middleware.Protected())
	care.Get("/records", h.GetCareRecords)
	care.Get("/fitness", h.GetFitnessRecords)
	care.Post("/records", h.RecordTreatment)
}

// RegisterPortalRoutes sets up routes for the Parent Portal (Sigmagate)
func RegisterPortalRoutes(router fiber.Router, service CareService) {
	h := NewCareHandler(service)
	
	portal := router.Group("/care", middleware.Protected())
	portal.Get("/records", h.GetCareRecords)
	portal.Get("/fitness", h.GetFitnessRecords)
}
