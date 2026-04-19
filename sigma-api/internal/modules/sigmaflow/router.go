package sigmaflow

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/middleware"
)

// RegisterAdminRoutes sets up routes for the Backoffice Finance management
func RegisterAdminRoutes(router fiber.Router, service FinanceService) {
	h := NewFinanceHandler(service)
	
	flow := router.Group("/flow", middleware.Protected(), middleware.AdminOnly())
	
	flow.Get("/invoices", func(c fiber.Ctx) error {
		// Logic to list all invoices (Admin only)
		return c.SendString("List of all invoices")
	})
	
	// Webhook is public but validated inside handler
	router.Post("/xendit/callback", h.XenditWebhook)
}

// RegisterPortalRoutes sets up routes for the Parent Portal (Sigmagate)
func RegisterPortalRoutes(router fiber.Router, service FinanceService) {
	h := NewFinanceHandler(service)
	
	portal := router.Group("/finance", middleware.Protected())
	portal.Get("/invoices", h.GetStudentInvoices)
	portal.Post("/invoices/:id/xendit", h.CreateXenditInvoice)
}
