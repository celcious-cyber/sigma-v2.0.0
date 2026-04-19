package sigmalit

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/middleware"
)

func RegisterAdminRoutes(router fiber.Router, service LibraryService) {
	h := NewLibraryHandler(service)
	
	lit := router.Group("/lit", middleware.Protected())
	lit.Get("/books", h.ListBooks)
}

func RegisterPortalRoutes(router fiber.Router, service LibraryService) {
	h := NewLibraryHandler(service)
	
	portal := router.Group("/lit", middleware.Protected())
	portal.Get("/books", h.ListBooks)
}
