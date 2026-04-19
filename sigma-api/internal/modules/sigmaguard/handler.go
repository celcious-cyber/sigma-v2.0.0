package sigmaguard

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/core/models"
)

type GuardHandler struct {
	service GuardService
}

func NewGuardHandler(s GuardService) *GuardHandler {
	return &GuardHandler{s}
}

// GetViolations handles listing violations with optional gender filter
func (h *GuardHandler) GetViolations(c fiber.Ctx) error {
	gender := c.Query("gender") // L or P
	violations, err := h.service.GetViolations(gender)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(violations)
}

// GetPermits handles listing permits
func (h *GuardHandler) GetPermits(c fiber.Ctx) error {
	gender := c.Query("gender")
	permits, err := h.service.GetPermits(gender)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(permits)
}

// RecordViolation handles POST request for new violation
func (h *GuardHandler) RecordViolation(c fiber.Ctx) error {
	var v models.Violation
	if err := c.Bind().JSON(&v); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.service.RecordViolation(v); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Violation recorded and points refreshed"})
}
