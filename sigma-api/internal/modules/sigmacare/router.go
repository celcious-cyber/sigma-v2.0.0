package sigmacare

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/middleware"
)

// RegisterAdminRoutes sets up routes for Backoffice UKS management
func RegisterAdminRoutes(router fiber.Router, service CareService) {
	handler := NewCareHandler(service)

	care := router.Group("/care", middleware.Protected())
	care.Get("/stats", handler.GetDashboardStats)
	care.Get("/visits", handler.GetVisits)
	care.Post("/visits", handler.CreateVisit)
	care.Get("/observations", handler.GetObservations)
	care.Post("/observations", handler.CreateObservation)
	care.Patch("/observations/:id", handler.UpdateObservation)
	care.Delete("/observations/:id", handler.DischargeObservation)
	care.Post("/certificates", handler.CreateCertificate)
	care.Get("/mcu", handler.GetMCURecords)
	care.Post("/mcu", handler.RecordMCU)
	care.Get("/medicines", handler.GetMedicines)
	care.Post("/medicines/mutation", handler.RecordMutation)
}

// RegisterPortalRoutes sets up routes for the Parent Portal (Sigmagate)
func RegisterPortalRoutes(router fiber.Router, service CareService) {
	// h := NewCareHandler(service)
	// portal := router.Group("/care", middleware.Protected())
	// portal.Get("/records", h.GetVisits)
}
