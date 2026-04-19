package auth

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/core/models"
	"sigma-api/internal/shared/database"
	"sigma-api/internal/shared/utils"
)

type LoginRequest struct {
	Identifier string `json:"identifier"` // Email or NIS
	Password   string `json:"password"`   // Password or BirthDate
}

// AdminLogin handles authentication for teachers and administrators
func AdminLogin(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Identifier).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
	}

	token, err := utils.GenerateToken(user.ID, user.Role, "admin")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// StudentLogin handles authentication for the Sigmagate PWA (Parents/Students)
func StudentLogin(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var student models.Student
	if err := database.DB.Where("nis = ? OR email = ?", req.Identifier, req.Identifier).First(&student).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Student not found"})
	}

	// Logic from legacy GateController: 
	// If password is set, check it. If not, fallback to BirthDate check.
	isValid := false
	if student.Password != nil && *student.Password != "" {
		isValid = utils.CheckPasswordHash(req.Password, *student.Password)
	} else if student.BirthDate != nil {
		// Fallback check: Compare birth date string (YYYY-MM-DD)
		storedDate := student.BirthDate.Format("2006-01-02")
		isValid = (storedDate == req.Password)
	}

	if !isValid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := utils.GenerateToken(student.ID, "Student", "student")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"student": fiber.Map{
			"id":   student.ID,
			"name": student.Name,
			"nis":  student.NIS,
		},
	})
}

// GetMe returns the currently authenticated user/student data
func GetMe(c fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	loginType := c.Locals("type").(string)

	if loginType == "admin" {
		var user models.User
		database.DB.First(&user, userID)
		return c.JSON(fiber.Map{"type": "admin", "data": user})
	}

	var student models.Student
	database.DB.Preload("Classroom").First(&student, userID)
	return c.JSON(fiber.Map{"type": "student", "data": student})
}
