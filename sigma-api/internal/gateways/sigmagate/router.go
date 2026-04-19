package sigmagate

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"sigma-api/internal/middleware"
	"sigma-api/internal/modules/sigmabase"
	"sigma-api/internal/modules/sigmaflow"
	"sigma-api/internal/modules/sigmaguard"
	"sigma-api/internal/modules/sigmacare"
	"sigma-api/internal/modules/sigmadesk"
	"sigma-api/internal/modules/sigmaedu"
	"sigma-api/internal/modules/sigmalit"
)

// SetupRouter initializes all parent portal routes
func SetupRouter(router fiber.Router, db *gorm.DB) {
	// 1. Initialize Module Services
	studentSvc := sigmabase.NewStudentService(db)
	financeSvc := sigmaflow.NewFinanceService(db)
	guardSvc := sigmaguard.NewGuardService(db, studentSvc)
	careSvc := sigmacare.NewCareService(db)
	deskSvc := sigmadesk.NewDeskService(db)
	eduSvc := sigmaedu.NewEduService(db)
	litSvc := sigmalit.NewLibraryService(db)
	h := NewSigmagateHandler(studentSvc)

	portal := router.Group("/portal", middleware.Protected())

	portal.Get("/status", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"gateway": "sigmagate",
			"status":  "active",
		})
	})

	// Dashboard (Thin UI ready)
	portal.Get("/dashboard", h.GetDashboard)

	// Finance Portal Routes
	sigmaflow.RegisterPortalRoutes(portal, financeSvc)

	// Guard Portal Routes
	sigmaguard.RegisterPortalRoutes(portal, guardSvc)

	// Care Portal Routes
	sigmacare.RegisterPortalRoutes(portal, careSvc)

	// Desk Portal Routes
	sigmadesk.RegisterPortalRoutes(portal, deskSvc)

	// Edu Portal Routes
	sigmaedu.RegisterPortalRoutes(portal, eduSvc)

	// Library Portal Routes
	sigmalit.RegisterPortalRoutes(portal, litSvc)
}
