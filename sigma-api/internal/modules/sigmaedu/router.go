package sigmaedu

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/middleware"
)

// RegisterAdminRoutes sets up routes for Backoffice Academic management
func RegisterAdminRoutes(router fiber.Router, service EduService) {
	h := NewEduHandler(service)
	
	edu := router.Group("/edu", middleware.Protected())
	edu.Get("/assessments", h.GetAssessments)
	edu.Get("/attendances", h.GetAttendances)
	edu.Post("/attendances", h.CreateAttendance)
	
	// TODO: Subjects, Schedules, Academic Calendar
}

// RegisterPortalRoutes sets up routes for the Parent Portal (Sigmagate)
func RegisterPortalRoutes(router fiber.Router, service EduService) {
	h := NewEduHandler(service)
	
	portal := router.Group("/edu", middleware.Protected())
	portal.Get("/assessments", h.GetAssessments)
	portal.Get("/attendances", h.GetAttendances)
}
