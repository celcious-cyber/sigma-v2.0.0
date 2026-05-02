package sigmaflow

import (
	"fmt"
	"time"
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

func (h *FinanceHandler) GetStudentInvoices(c fiber.Ctx) error {
	studentID := c.Locals("user_id").(uint)
	invoices, err := h.service.GetInvoicesByStudent(studentID)
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(invoices)
}

func (h *FinanceHandler) GetGlobalStats(c fiber.Ctx) error {
	stats, err := h.service.GetGlobalFinanceStats()
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(stats)
}

func (h *FinanceHandler) CreateBulkInvoices(c fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.Bind().JSON(&payload); err != nil { return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"}) }
	count, err := h.service.CreateBulkInvoices(payload)
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(fiber.Map{"message": fmt.Sprintf("Berhasil membuat %d faktur", count), "count": count})
}

func (h *FinanceHandler) GetAllInvoices(c fiber.Ctx) error {
	invoices, err := h.service.GetAllInvoices()
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(invoices)
}

func (h *FinanceHandler) GetAllPayments(c fiber.Ctx) error {
	payments, err := h.service.GetAllPayments()
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(payments)
}

func (h *FinanceHandler) ProcessManualPayment(c fiber.Ctx) error {
	var payload struct {
		InvoiceID uint    `json:"invoice_id"`
		Amount    float64 `json:"amount"`
		Method    string  `json:"method"`
		AdminName string  `json:"admin_name"`
	}
	if err := c.Bind().JSON(&payload); err != nil { return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"}) }
	err := h.service.ProcessManualPayment(payload.InvoiceID, payload.Amount, payload.Method, payload.AdminName)
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(fiber.Map{"message": "Pembayaran berhasil dicatat"})
}

func (h *FinanceHandler) GetCashFlow(c fiber.Ctx) error {
	filters := map[string]string{ "type": c.Query("type"), "start_date": c.Query("start_date"), "end_date": c.Query("end_date") }
	entries, err := h.service.GetAllCashFlows(filters)
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(entries)
}

func (h *FinanceHandler) CreateCashFlow(c fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.Bind().JSON(&payload); err != nil { return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"}) }
	adminName, _ := payload["admin_name"].(string)
	if adminName == "" { adminName = "Admin" }
	entry, err := h.service.CreateCashFlowEntry(payload, adminName)
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.Status(201).JSON(entry)
}

func (h *FinanceHandler) DeleteCashFlow(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	if err := h.service.DeleteCashFlowEntry(id); err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(fiber.Map{"message": "Entry deleted"})
}

func (h *FinanceHandler) GetCashFlowSummary(c fiber.Ctx) error {
	summary, err := h.service.GetCashFlowSummary()
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(summary)
}

func (h *FinanceHandler) SyncCashFlow(c fiber.Ctx) error {
	count, err := h.service.SyncExistingPayments()
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(fiber.Map{"message": fmt.Sprintf("Berhasil menyinkronkan %d transaksi", count), "count": count})
}

func (h *FinanceHandler) CreateInvoice(c fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.Bind().JSON(&payload); err != nil { return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"}) }
	invoice, err := h.service.CreateInvoice(payload)
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.Status(201).JSON(invoice)
}

func (h *FinanceHandler) GetCategories(c fiber.Ctx) error {
	cats, err := h.service.GetCategories()
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(cats)
}

func (h *FinanceHandler) CreateCategory(c fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.Bind().JSON(&payload); err != nil { return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"}) }
	cat, err := h.service.CreateCategory(payload)
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.Status(201).JSON(cat)
}

func (h *FinanceHandler) UpdateCategory(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	var payload map[string]interface{}
	if err := c.Bind().JSON(&payload); err != nil { return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"}) }
	cat, err := h.service.UpdateCategory(id, payload)
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(cat)
}

func (h *FinanceHandler) DeleteCategory(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	if err := h.service.DeleteCategory(id); err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(fiber.Map{"message": "Category deleted"})
}

func (h *FinanceHandler) FindStudentByNIS(c fiber.Ctx) error {
	student, err := h.service.FindStudentByNIS(c.Params("nis"))
	if err != nil { return c.Status(404).JSON(fiber.Map{"error": "Santri tidak ditemukan"}) }
	return c.JSON(student)
}

func (h *FinanceHandler) UpdateInvoice(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	var payload map[string]interface{}
	if err := c.Bind().JSON(&payload); err != nil { return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"}) }
	inv, err := h.service.UpdateInvoice(id, payload)
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(inv)
}

func (h *FinanceHandler) DeleteInvoice(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	if err := h.service.DeleteInvoice(id); err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(fiber.Map{"message": "Invoice deleted"})
}

func (h *FinanceHandler) CreateXenditInvoice(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	invoice, err := h.service.CreateXenditInvoice(id)
	if err != nil { return c.Status(400).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(invoice)
}

func (h *FinanceHandler) XenditWebhook(c fiber.Ctx) error {
	if c.Get("X-CALLBACK-TOKEN") != config.AppConfig.XenditCallbackToken {
		return c.Status(403).JSON(fiber.Map{"error": "Invalid token"})
	}
	var payload struct { ID string `json:"id"`; ExternalID string `json:"external_id"`; Status string `json:"status"`; Amount float64 `json:"amount"` }
	if err := c.Bind().JSON(&payload); err != nil { return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"}) }
	if payload.Status != "PAID" { return c.SendStatus(200) }
	err := h.service.ProcessPaymentCallback(payload.ExternalID, payload.Amount, "Xendit/Online")
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(fiber.Map{"message": "success"})
}
func (h *FinanceHandler) GetProfitLoss(c fiber.Ctx) error {
	month := utils.StringToInt(c.Query("month"))
	year := utils.StringToInt(c.Query("year"))
	if month == 0 { month = int(time.Now().Month()) }
	if year == 0 { year = time.Now().Year() }

	report, err := h.service.GetProfitLossReport(month, year)
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(report)
}

func (h *FinanceHandler) GetArrearsReport(c fiber.Ctx) error {
	report, err := h.service.GetArrearsReport()
	if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(report)
}
