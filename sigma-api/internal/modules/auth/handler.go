package auth

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/core/models"
	"sigma-api/internal/shared/database"
	"sigma-api/internal/shared/utils"
	"os"
	"path/filepath"
	"fmt"
	"time"
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

// UpdateProfile handles updating user/student personal information
func UpdateProfile(c fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	loginType := c.Locals("type").(string)

	if loginType == "admin" {
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		type UpdateRequest struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		var req UpdateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		user.Name = req.Name
		user.Email = req.Email
		database.DB.Save(&user)
		return c.JSON(fiber.Map{"message": "Profile updated successfully", "data": user})
	}

	var student models.Student
	if err := database.DB.First(&student, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	type StudentUpdateRequest struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	var req StudentUpdateRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	student.Name = req.Name
	student.Phone = &req.Phone
	database.DB.Save(&student)
	return c.JSON(fiber.Map{"message": "Profile updated successfully", "data": student})
}

// UpdatePassword handles changing the user's password
func UpdatePassword(c fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	loginType := c.Locals("type").(string)

	type PasswordRequest struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	var req PasswordRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if loginType == "admin" {
		var user models.User
		database.DB.First(&user, userID)

		if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password lama tidak sesuai"})
		}

		hashed, _ := utils.HashPassword(req.NewPassword)
		user.Password = hashed
		database.DB.Save(&user)
	} else {
		var student models.Student
		database.DB.First(&student, userID)

		// Students might not have a password yet (fallback to BirthDate logic)
		currentPassword := ""
		if student.Password != nil {
			currentPassword = *student.Password
		}

		if currentPassword != "" {
			if !utils.CheckPasswordHash(req.OldPassword, currentPassword) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password lama tidak sesuai"})
			}
		} else if student.BirthDate != nil {
			if student.BirthDate.Format("2006-01-02") != req.OldPassword {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password lama tidak sesuai"})
			}
		}

		hashed, _ := utils.HashPassword(req.NewPassword)
		student.Password = &hashed
		database.DB.Save(&student)
	}

	return c.JSON(fiber.Map{"message": "Password berhasil diperbarui"})
}

// UploadAvatar handles uploading and saving a user's profile picture
func UploadAvatar(c fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	loginType := c.Locals("type").(string)

	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tidak ada file yang diunggah"})
	}

	// Create uploads directory if not exists
	uploadDir := "./uploads/avatars"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d_%d%s", loginType, userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan file"})
	}

	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)

	if loginType == "admin" {
		var user models.User
		database.DB.First(&user, userID)
		user.AvatarURL = &avatarURL
		database.DB.Save(&user)
	} else {
		var student models.Student
		database.DB.First(&student, userID)
		student.AvatarURL = &avatarURL
		database.DB.Save(&student)
	}

	return c.JSON(fiber.Map{
		"message": "Foto profil berhasil diperbarui",
		"url":     avatarURL,
	})
}
