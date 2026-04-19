package main

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/gateways/admin"
	"sigma-api/internal/gateways/sigmagate"
	"sigma-api/internal/modules/auth"
	"sigma-api/internal/shared/config"
	"sigma-api/internal/shared/database"
	"sigma-api/internal/shared/logger"
	sigma_models "sigma-api/internal/core/models"
	fiberLogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	// 1. Initialize Core Infrastructure
	config.InitConfig()
	logger.InitLogger()
	database.InitDB(config.AppConfig.DBPath)

	// 2. Initialize Fiber v3 with Zero-Allocation JSON Encoder
	app := fiber.New(fiber.Config{
		AppName:     "SIGMA v2.0.0 Modular Monolith",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// 3. Register Gateways (BFF Strategy)
	api := app.Group("/api/v1")
	
	// Global Middlewares
	api.Use(recover.New())
	api.Use(fiberLogger.New()) 

	// 3.1 Register Auth Module (Shared)
	auth.RegisterRoutes(api)

	// 3.2 Register Gateways (BFF Strategy)
	// API Grouping to separate logic for different entry points
	
	// Admin Gateway (Backoffice)
	admin.SetupRouter(api, database.DB)
	
	// Sigmagate Gateway (Parent Portal)
	sigmagate.SetupRouter(api, database.DB)

	// Serve Static Files (Photos, etc.)
	app.Get("/uploads/*", static.New("./uploads"))

	// 4. Default Health Check
	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "online",
			"version": "v2.0.0",
			"gateways": []string{"admin", "sigmagate"},
		})
	})

	app.Get("/test-alumni", func(c fiber.Ctx) error {
		var al []sigma_models.Alumni
		if err := database.DB.Find(&al).Error; err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(al)
	})

	// 5. Start Server
	log.Printf("SIGMA API starting on :3000")
	log.Fatal(app.Listen(":3000"))
}
