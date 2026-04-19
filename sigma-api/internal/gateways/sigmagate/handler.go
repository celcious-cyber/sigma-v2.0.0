package sigmagate

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/core/models"
	"sigma-api/internal/modules/sigmabase"
	"sigma-api/internal/shared/database"
)

type SigmagateHandler struct {
	studentSvc sigmabase.StudentService
}

func NewSigmagateHandler(s sigmabase.StudentService) *SigmagateHandler {
	return &SigmagateHandler{s}
}

// GetDashboard returns the pre-aggregated dashboard for the Parent Portal
func (h *SigmagateHandler) GetDashboard(c fiber.Ctx) error {
	studentID := c.Locals("user_id").(uint)
	
	// Siksa backend phase: Aggregating all data into one premium payload
	student, err := h.studentSvc.GetStudentByID(studentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	// 1. Attendance Stats (Recent 5)
	var recentAttendance []models.Attendance
	database.DB.Where("student_id = ?", studentID).Order("date desc").Limit(5).Find(&recentAttendance)

	// 2. Permit Quotas (Pre-calculated in service)
	isBlocked, permitStats, _ := h.studentSvc.IsPermitBlocked(studentID)

	// 3. Health Stats (Recent 8)
	var fitnessHistory []models.FitnessRecord
	database.DB.Where("student_id = ?", studentID).Order("date desc").Limit(8).Find(&fitnessHistory)

	// 4. Global Announcements
	var announcements []models.Announcement
	database.DB.Where("is_active = ?", true).Order("created_at desc").Limit(5).Find(&announcements)

	// Ready-to-Display JSON (Thin UI approach)
	return c.JSON(fiber.Map{
		"student":      student,
		"attendance":   recentAttendance,
		"permit_stats": permitStats,
		"is_blocked":   isBlocked,
		"fitness":      fitnessHistory,
		"announcements": announcements,
	})
}
