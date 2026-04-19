package auth

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/middleware"
)

// RegisterRoutes sets up the authentication endpoints
func RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	// Public routes
	auth.Post("/admin/login", AdminLogin)
	auth.Post("/student/login", StudentLogin)

	// Protected routes
	auth.Get("/me", middleware.Protected(), GetMe)
}
