package sigmaguard

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/middleware"
)

// RegisterAdminRoutes sets up routes for Backoffice Discipline management
func RegisterAdminRoutes(router fiber.Router, service GuardService) {
	h := NewGuardHandler(service)
	
	guard := router.Group("/guard", middleware.Protected())
	
	// Violations
	guard.Get("/violations", h.GetViolations)
	guard.Post("/violations", h.RecordViolation)
	
	// Permits
	guard.Get("/permits", h.GetPermits)
	
	// TODO: Sanctions, Violation Dictionary
}

// RegisterPortalRoutes sets up routes for Parent Portal (Sigmagate)
func RegisterPortalRoutes(router fiber.Router, service GuardService) {
	h := NewGuardHandler(service)
	
	portal := router.Group("/guard", middleware.Protected())
	portal.Get("/violations", h.GetViolations) // Usually filtered by student_id in handler or logic
	portal.Get("/permits", h.GetPermits)
}
