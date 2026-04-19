package sigmacare

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/core/models"
)

type CareHandler struct {
	service CareService
}

func NewCareHandler(s CareService) *CareHandler {
	return &CareHandler{s}
}

// GetCareRecords handles listing medical records
func (h *CareHandler) GetCareRecords(c fiber.Ctx) error {
	gender := c.Query("gender")
	records, err := h.service.GetCareRecords(gender)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(records)
}

// GetFitnessRecords handles listing weight/height history
func (h *CareHandler) GetFitnessRecords(c fiber.Ctx) error {
	gender := c.Query("gender")
	records, err := h.service.GetFitnessRecords(gender)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(records)
}

// RecordTreatment handles recording a new medical event
func (h *CareHandler) RecordTreatment(c fiber.Ctx) error {
	var record models.CareRecord
	if err := c.Bind().JSON(&record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.service.RecordCare(record); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Treatment record saved successfully"})
}
