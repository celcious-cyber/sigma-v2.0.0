package sigmadesk

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/middleware"
)

// RegisterAdminRoutes sets up routes for Backoffice Secretariat management
func RegisterAdminRoutes(router fiber.Router, service DeskService) {
	h := NewDeskHandler(service)
	
	desk := router.Group("/desk", middleware.Protected())
	desk.Get("/announcements", h.GetAnnouncements)
	desk.Post("/announcements", h.CreateAnnouncement)
	desk.Get("/visitors", h.GetVisitors)
	
	// TODO: Letters (Incoming/Outgoing)
}

// RegisterPortalRoutes sets up routes for the Parent Portal (Sigmagate)
func RegisterPortalRoutes(router fiber.Router, service DeskService) {
	h := NewDeskHandler(service)
	
	portal := router.Group("/desk", middleware.Protected())
	portal.Get("/announcements", h.GetAnnouncements)
}
