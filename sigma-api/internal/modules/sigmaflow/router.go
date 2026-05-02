package sigmaflow

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/middleware"
)

// RegisterAdminRoutes sets up routes for the Backoffice Finance management
func RegisterAdminRoutes(router fiber.Router, service FinanceService) {
	h := NewFinanceHandler(service)
	
	// Group under /admin/flow
	flow := router.Group("/flow", middleware.Protected(), middleware.AdminOnly())
	
	flow.Get("/stats", h.GetGlobalStats)
	flow.Get("/invoices", h.GetAllInvoices)
	flow.Get("/payments", h.GetAllPayments)
	flow.Post("/payments/manual", h.ProcessManualPayment)
	flow.Post("/invoices/bulk", h.CreateBulkInvoices)
	flow.Post("/invoices", h.CreateInvoice)
	flow.Put("/invoices/:id", h.UpdateInvoice)
	flow.Delete("/invoices/:id", h.DeleteInvoice)
	flow.Get("/categories", h.GetCategories)
	flow.Post("/categories", h.CreateCategory)
	flow.Put("/categories/:id", h.UpdateCategory)
	flow.Delete("/categories/:id", h.DeleteCategory)
	flow.Get("/students/nis/:nis", h.FindStudentByNIS)
	
	// Arus Kas
	flow.Get("/cashflow", h.GetCashFlow)
	flow.Post("/cashflow", h.CreateCashFlow)
	flow.Delete("/cashflow/:id", h.DeleteCashFlow)
	flow.Get("/cashflow/summary", h.GetCashFlowSummary)
	flow.Post("/cashflow/sync", h.SyncCashFlow)
	
	// Laporan
	flow.Get("/reports/profit-loss", h.GetProfitLoss)
	flow.Get("/reports/arrears", h.GetArrearsReport)
	
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
