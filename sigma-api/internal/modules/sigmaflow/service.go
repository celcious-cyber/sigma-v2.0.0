package sigmaflow

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"sigma-api/internal/core/models"
	"sigma-api/internal/shared/utils"
)

type FinanceService interface {
	GetInvoicesByStudent(studentID uint) ([]map[string]interface{}, error)
	CreateXenditInvoice(invoiceID uint) (*models.Invoice, error)
	ProcessPaymentCallback(externalID string, amount float64, method string) error
}

type financeService struct {
	db *gorm.DB
}

func NewFinanceService(db *gorm.DB) FinanceService {
	return &financeService{db}
}

// GetInvoicesByStudent returns formatted invoices for the "Thin UI"
func (s *financeService) GetInvoicesByStudent(studentID uint) ([]map[string]interface{}, error) {
	var invoices []models.Invoice
	err := s.db.Preload("Category").Where("student_id = ?", studentID).Order("created_at desc").Find(&invoices).Error
	if err != nil {
		return nil, err
	}

	// Siksa backend phase: Format to ready-to-display JSON
	result := make([]map[string]interface{}, len(invoices))
	for i, inv := range invoices {
		result[i] = map[string]interface{}{
			"id":             inv.ID,
			"category_name":  inv.PaymentCategory.Name,
			"amount":         inv.Amount,
			"amount_display": utils.FormatRupiah(inv.Amount),
			"due_date":       inv.DueDate.Format("02 Jan 2006"),
			"status":         inv.Status,
			"payment_url":    inv.PaymentURL,
			"payment_status": inv.PaymentStatus,
		}
	}
	return result, nil
}

// CreateXenditInvoice generates an online payment link
func (s *financeService) CreateXenditInvoice(invoiceID uint) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := s.db.Preload("Student").Preload("PaymentCategory").First(&invoice, invoiceID).Error; err != nil {
		return nil, err
	}

	if invoice.Status == "Paid" {
		return nil, errors.New("invoice already paid")
	}

	// Reuse existing link if pending
	if invoice.PaymentURL != nil && invoice.PaymentStatus == "PENDING" {
		return &invoice, nil
	}

	xenditReq := utils.XenditInvoiceRequest{
		ExternalID:      fmt.Sprintf("%d-%d", invoice.ID, time.Now().Unix()),
		Amount:          invoice.Amount,
		PayerEmail:      "pondok.sigma@example.com", // Fallback email
		Description:     fmt.Sprintf("%s - %s", invoice.PaymentCategory.Name, invoice.Student.Name),
		InvoiceDuration: 86400,
	}

	resp, err := utils.CreateXenditInvoice(xenditReq)
	if err != nil {
		return nil, err
	}

	invoice.XenditID = &resp.ID
	invoice.PaymentURL = &resp.InvoiceURL
	invoice.PaymentStatus = "PENDING"
	s.db.Save(&invoice)

	return &invoice, nil
}

// ProcessPaymentCallback handles webhook data from Xendit
func (s *financeService) ProcessPaymentCallback(externalID string, amount float64, method string) error {
	// Extract internal ID from external_id (format: id-timestamp)
	var invoiceID uint
	fmt.Sscanf(externalID, "%d-", &invoiceID)

	return s.db.Transaction(func(tx *gorm.DB) error {
		var invoice models.Invoice
		if err := tx.First(&invoice, invoiceID).Error; err != nil {
			return err
		}

		if invoice.Status == "Paid" {
			return nil // Already processed
		}

		// 1. Create Payment record
		payment := models.Payment{
			InvoiceID:   invoice.ID,
			PaymentDate: time.Now(),
			AmountPaid:  amount,
			Method:      method,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		// 2. Update Invoice status
		invoice.Status = "Paid" // Simplified logic, could check partial
		invoice.PaymentStatus = "PAID"
		return tx.Save(&invoice).Error
	})
}
