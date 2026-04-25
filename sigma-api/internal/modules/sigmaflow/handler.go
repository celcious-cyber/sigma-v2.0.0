package sigmaflow

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/shared/config"
	"sigma-api/internal/shared/utils"
)

type FinanceHandler struct {
	service FinanceService
}

func NewFinanceHandler(s FinanceService) *FinanceHandler {
	return &FinanceHandler{s}
}

// GetStudentInvoices returns invoices for the PWA dashboard
func (h *FinanceHandler) GetStudentInvoices(c fiber.Ctx) error {
	studentID := c.Locals("user_id").(uint)
	invoices, err := h.service.GetInvoicesByStudent(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(invoices)
}

// GetGlobalStats returns overall financial overview for admin
func (h *FinanceHandler) GetGlobalStats(c fiber.Ctx) error {
	stats, err := h.service.GetGlobalFinanceStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stats)
}

func (h *FinanceHandler) CreateBulkInvoices(c fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	count, err := h.service.CreateBulkInvoices(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Berhasil membuat %d faktur secara massal", count),
		"count":   count,
	})
}

func (h *FinanceHandler) GetAllInvoices(c fiber.Ctx) error {
	invoices, err := h.service.GetAllInvoices()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(invoices)
}

func (h *FinanceHandler) CreateInvoice(c fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	invoice, err := h.service.CreateInvoice(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(invoice)
}

func (h *FinanceHandler) GetCategories(c fiber.Ctx) error {
	cats, err := h.service.GetCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cats)
}

func (h *FinanceHandler) CreateCategory(c fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	cat, err := h.service.CreateCategory(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(cat)
}

func (h *FinanceHandler) UpdateCategory(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	var payload map[string]interface{}
	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	cat, err := h.service.UpdateCategory(id, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cat)
}

func (h *FinanceHandler) DeleteCategory(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	if err := h.service.DeleteCategory(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Category deleted"})
}

func (h *FinanceHandler) FindStudentByNIS(c fiber.Ctx) error {
	nis := c.Params("nis")
	student, err := h.service.FindStudentByNIS(nis)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Santri tidak ditemukan"})
	}
	return c.JSON(student)
}

func (h *FinanceHandler) UpdateInvoice(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	var payload map[string]interface{}
	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	inv, err := h.service.UpdateInvoice(id, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(inv)
}

func (h *FinanceHandler) DeleteInvoice(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	if err := h.service.DeleteInvoice(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Invoice deleted"})
}

// CreateXenditInvoice generates a payment link
func (h *FinanceHandler) CreateXenditInvoice(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	invoice, err := h.service.CreateXenditInvoice(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(invoice)
}

// XenditWebhook handles incoming payment notifications
func (h *FinanceHandler) XenditWebhook(c fiber.Ctx) error {
	// Security Check
	callbackToken := c.Get("X-CALLBACK-TOKEN")
	if callbackToken != config.AppConfig.XenditCallbackToken {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Invalid callback token"})
	}

	var payload struct {
		ID         string  `json:"id"`
		ExternalID string  `json:"external_id"`
		Status     string  `json:"status"`
		Amount     float64 `json:"amount"`
	}

	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if payload.Status != "PAID" {
		return c.SendStatus(fiber.StatusOK) // Ignore non-paid events but return 200
	}

	err := h.service.ProcessPaymentCallback(payload.ExternalID, payload.Amount, "Xendit/Online")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "success"})
}
