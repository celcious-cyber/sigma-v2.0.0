package sigmacare

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/core/models"
)

type CareHandler struct {
	service CareService
}

func NewCareHandler(s CareService) *CareHandler {
	return &CareHandler{s}
}

func (h *CareHandler) GetDashboardStats(c fiber.Ctx) error {
	stats, err := h.service.GetClinicalDashboardStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stats)
}

func (h *CareHandler) CreateVisit(c fiber.Ctx) error {
	var visit models.MedicalVisit
	if err := c.Bind().JSON(&visit); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	// Set RecordedBy from JWT context
	if userID, ok := c.Locals("user_id").(float64); ok {
		visit.RecordedBy = uint(userID)
	} else if userID, ok := c.Locals("user_id").(uint); ok {
		visit.RecordedBy = userID
	}

	if err := h.service.CreateMedicalVisit(visit); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Visit recorded successfully"})
}

func (h *CareHandler) GetVisits(c fiber.Ctx) error {
	visits, err := h.service.GetMedicalVisits()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(visits)
}

func (h *CareHandler) RecordMCU(c fiber.Ctx) error {
	var record models.AnthropometryRecord
	if err := c.Bind().JSON(&record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	// Set RecordedBy from JWT context
	if userID, ok := c.Locals("user_id").(float64); ok {
		record.RecordedBy = uint(userID)
	} else if userID, ok := c.Locals("user_id").(uint); ok {
		record.RecordedBy = userID
	}

	if err := h.service.RecordAnthropometry(record); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "MCU record saved successfully"})
}

func (h *CareHandler) GetMCURecords(c fiber.Ctx) error {
	records, err := h.service.GetAnthropometryRecords()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(records)
}

func (h *CareHandler) GetMedicines(c fiber.Ctx) error {
	medicines, err := h.service.GetMedicines()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(medicines)
}

func (h *CareHandler) RecordMutation(c fiber.Ctx) error {
	var mutation models.MedicineMutation
	if err := c.Bind().JSON(&mutation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.RecordMedicineMutation(mutation); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Inventory updated successfully"})
}

func (h *CareHandler) CreateCertificate(c fiber.Ctx) error {
	var cert models.MedicalCertificate
	if err := c.Bind().JSON(&cert); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.CreateCertificate(cert); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Certificate issued successfully"})
}

func (h *CareHandler) CreateObservation(c fiber.Ctx) error {
	var obs models.ObservationRecord
	if err := c.Bind().JSON(&obs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.CreateObservation(obs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Observation started successfully"})
}

func (h *CareHandler) GetObservations(c fiber.Ctx) error {
	obs, err := h.service.GetActiveObservations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(obs)
}

func (h *CareHandler) UpdateObservation(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var updates map[string]interface{}
	if err := c.Bind().JSON(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.UpdateObservation(uint(id), updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Observation updated successfully"})
}

func (h *CareHandler) DischargeObservation(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.DischargeObservation(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": "Gagal memproses kepulangan santri",
			"details": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "Patient discharged successfully"})
}
