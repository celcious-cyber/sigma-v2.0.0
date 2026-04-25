package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/shared/utils"
)

// Protected handles JWT validation for secured routes
func Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			fmt.Println("[AUTH ERROR] Missing authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			fmt.Printf("[AUTH ERROR] %v\n", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Store claims in context for later use in handlers
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("type", claims.Type)

		return c.Next()
	}
}

// AdminOnly restricts access to users with Administrator role
func AdminOnly() fiber.Handler {
	return func(c fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		fmt.Printf("DEBUG: AdminOnly middleware - role: '%s', ok: %v\n", role, ok)
		if !ok || (role != "Administrator" && role != "Admin") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: Administrator only",
			})
		}
		return c.Next()
	}
}
