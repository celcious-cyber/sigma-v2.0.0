package admin

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"sigma-api/internal/modules/sigmabase"
	"sigma-api/internal/modules/sigmaflow"
	"sigma-api/internal/modules/sigmaguard"
	"sigma-api/internal/modules/sigmacare"
	"sigma-api/internal/modules/sigmadesk"
	"sigma-api/internal/modules/sigmaedu"
	"sigma-api/internal/modules/sigmalit"
	"sigma-api/internal/modules/settings"
)

// SetupRouter initializes all admin backoffice routes
func SetupRouter(router fiber.Router, db *gorm.DB) {
	// 1. Initialize Module Services
	studentSvc := sigmabase.NewStudentService(db)
	financeSvc := sigmaflow.NewFinanceService(db)
	guardSvc := sigmaguard.NewGuardService(db, studentSvc)
	careSvc := sigmacare.NewCareService(db)
	deskSvc := sigmadesk.NewDeskService(db)
	eduSvc := sigmaedu.NewEduService(db)
	litSvc := sigmalit.NewLibraryService(db)
	settingsSvc := settings.NewSettingsService(db)

	admin := router.Group("/admin")

	admin.Get("/status", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"gateway": "admin",
			"status":  "active",
		})
	})

	// 2. Register Module Routes
	sigmabase.RegisterAdminRoutes(admin, studentSvc)
	sigmaflow.RegisterAdminRoutes(admin, financeSvc)
	sigmaguard.RegisterAdminRoutes(admin, guardSvc)
	sigmacare.RegisterAdminRoutes(admin, careSvc)
	sigmadesk.RegisterAdminRoutes(admin, deskSvc)
	sigmaedu.RegisterAdminRoutes(admin, eduSvc)
	sigmalit.RegisterAdminRoutes(admin, litSvc)
	settings.RegisterAdminRoutes(admin, settingsSvc)
}
