package sigmaedu

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/core/models"
)

type EduHandler struct {
	service EduService
}

func NewEduHandler(s EduService) *EduHandler {
	return &EduHandler{s}
}

// GetAssessments handles listing student scores
func (h *EduHandler) GetAssessments(c fiber.Ctx) error {
	studentID := c.Locals("user_id").(uint) // If portal
	if c.Query("student_id") != "" {
		// Logic for admin can go here
	}
	
	assessments, err := h.service.GetAssessmentsByStudent(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(assessments)
}

// GetAttendances handles listing attendance logs
func (h *EduHandler) GetAttendances(c fiber.Ctx) error {
	studentID := c.Locals("user_id").(uint)
	attendances, err := h.service.GetAttendancesByStudent(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(attendances)
}

// CreateAttendance handles recording a new attendance event
func (h *EduHandler) CreateAttendance(c fiber.Ctx) error {
	var a models.Attendance
	if err := c.Bind().JSON(&a); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.service.RecordAttendance(a); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Attendance recorded successfully"})
}
