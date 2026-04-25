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
	GetGlobalFinanceStats() (map[string]interface{}, error)
	GetAllInvoices() ([]map[string]interface{}, error)
	GetCategories() ([]models.PaymentCategory, error)
	CreateCategory(data map[string]interface{}) (*models.PaymentCategory, error)
	UpdateCategory(id uint, data map[string]interface{}) (*models.PaymentCategory, error)
	DeleteCategory(id uint) error
	CreateInvoice(data map[string]interface{}) (*models.Invoice, error)
	CreateBulkInvoices(data map[string]interface{}) (int, error)
	UpdateInvoice(id uint, data map[string]interface{}) (*models.Invoice, error)
	DeleteInvoice(id uint) error
	FindStudentByNIS(nis string) (*models.Student, error)
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
		invoice.Status = "Paid" 
		invoice.PaymentStatus = "PAID"
		return tx.Save(&invoice).Error
	})
}

func (s *financeService) GetGlobalFinanceStats() (map[string]interface{}, error) {
	var totalRevenue float64
	var pendingInvoices float64
	var totalInvoices int64

	// 1. Total Revenue (from paid payments)
	s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount_paid), 0)").Scan(&totalRevenue)

	// 2. Pending Invoices (total amount of unpaid/partial invoices)
	s.db.Model(&models.Invoice{}).Where("status != ?", "Paid").Select("COALESCE(SUM(amount), 0)").Scan(&pendingInvoices)

	// 3. Total Invoices Count
	s.db.Model(&models.Invoice{}).Count(&totalInvoices)

	// 4. Recent Transactions
	var recentPayments []models.Payment
	s.db.Preload("Invoice.Student").Preload("Invoice.PaymentCategory").Order("payment_date desc").Limit(10).Find(&recentPayments)

	formattedPayments := make([]map[string]interface{}, len(recentPayments))
	for i, p := range recentPayments {
		formattedPayments[i] = map[string]interface{}{
			"id":            p.ID,
			"student_name":  p.Invoice.Student.Name,
			"category":      p.Invoice.PaymentCategory.Name,
			"amount":        p.AmountPaid,
			"amount_formatted": utils.FormatRupiah(p.AmountPaid),
			"date":          p.PaymentDate.Format("02 Jan 2006, 15:04"),
			"method":        p.Method,
		}
	}

	return map[string]interface{}{
		"total_revenue":      totalRevenue,
		"revenue_formatted":  utils.FormatRupiah(totalRevenue),
		"pending_amount":     pendingInvoices,
		"pending_formatted":  utils.FormatRupiah(pendingInvoices),
		"total_invoices":     totalInvoices,
		"recent_payments":    formattedPayments,
	}, nil
}
func (s *financeService) GetAllInvoices() ([]map[string]interface{}, error) {
	var invoices []models.Invoice
	err := s.db.Preload("Student").Preload("PaymentCategory").Order("created_at desc").Find(&invoices).Error
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, inv := range invoices {
		studentName := "Unknown"
		if inv.Student != nil {
			studentName = inv.Student.Name
		}
		categoryName := "General"
		if inv.PaymentCategory != nil {
			categoryName = inv.PaymentCategory.Name
		}

		result = append(result, map[string]interface{}{
			"id":               inv.ID,
			"invoice_number":   inv.InvoiceNumber,
			"student_name":     studentName,
			"category":         categoryName,
			"amount":           inv.Amount,
			"amount_formatted": utils.FormatRupiah(inv.Amount),
			"status":           inv.Status, // Paid, Unpaid, Partial
			"due_date":         inv.DueDate.Format("02 Jan 2006"),
			"payment_url":      inv.PaymentURL,
		})
	}
	return result, nil
}

func (s *financeService) GetCategories() ([]models.PaymentCategory, error) {
	var cats []models.PaymentCategory
	err := s.db.Find(&cats).Error
	return cats, err
}

func (s *financeService) CreateCategory(data map[string]interface{}) (*models.PaymentCategory, error) {
	name, _ := data["name"].(string)
	amount, _ := data["amount"].(float64)
	isRecurring, _ := data["is_recurring"].(bool)

	cat := models.PaymentCategory{
		Name:        name,
		Amount:      amount,
		IsRecurring: isRecurring,
	}

	if err := s.db.Create(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (s *financeService) UpdateCategory(id uint, data map[string]interface{}) (*models.PaymentCategory, error) {
	var cat models.PaymentCategory
	if err := s.db.First(&cat, id).Error; err != nil {
		return nil, err
	}

	if name, ok := data["name"].(string); ok {
		cat.Name = name
	}
	if amount, ok := data["amount"].(float64); ok {
		cat.Amount = amount
	}
	if isRecurring, ok := data["is_recurring"].(bool); ok {
		cat.IsRecurring = isRecurring
	}

	if err := s.db.Save(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (s *financeService) DeleteCategory(id uint) error {
	return s.db.Delete(&models.PaymentCategory{}, id).Error
}

func (s *financeService) CreateInvoice(data map[string]interface{}) (*models.Invoice, error) {
	// Generate Invoice Number: INV/YYYYMMDD/RAND
	invNum := fmt.Sprintf("INV/%s/%04d", time.Now().Format("20060102"), time.Now().Nanosecond()%10000)

	dueDateStr, ok := data["due_date"].(string)
	if !ok || dueDateStr == "" {
		return nil, errors.New("jatuh tempo harus diisi")
	}
	dueDate, _ := time.Parse("2006-01-02", dueDateStr)

	sID := getUint(data["student_id"])
	cID := getUint(data["payment_category_id"])

	// Validate Student exists
	var studentCount int64
	s.db.Model(&models.Student{}).Where("id = ?", sID).Count(&studentCount)
	if studentCount == 0 {
		return nil, fmt.Errorf("santri dengan ID %d tidak ditemukan", sID)
	}

	invoice := models.Invoice{
		StudentID:         sID,
		PaymentCategoryID: cID,
		InvoiceNumber:     invNum,
		Amount:            getFloat(data["amount"]),
		DueDate:           dueDate,
		Status:            "Unpaid",
	}

	if err := s.db.Create(&invoice).Error; err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (s *financeService) CreateBulkInvoices(data map[string]interface{}) (int, error) {
	cID := getUint(data["payment_category_id"])
	amount := getFloat(data["amount"])
	dueDateStr := data["due_date"].(string)
	dueDate, _ := time.Parse("2006-01-02", dueDateStr)
	target := data["target"].(string) // "all" or "classroom"
	
	var students []models.Student
	query := s.db.Model(&models.Student{})
	
	if target == "classroom" {
		classID := getUint(data["classroom_id"])
		query = query.Where("classroom_id = ?", classID)
	}
	
	if err := query.Find(&students).Error; err != nil {
		return 0, err
	}

	if len(students) == 0 {
		return 0, errors.New("tidak ada santri yang sesuai kriteria")
	}

	count := 0
	now := time.Now()
	for i, student := range students {
		// Generate Unique Invoice Number using timestamp and index
		invNum := fmt.Sprintf("INV/%s/%d%d%d", now.Format("20060102"), student.ID, now.Unix()%1000, i)
		
		invoice := models.Invoice{
			StudentID:         student.ID,
			PaymentCategoryID: cID,
			InvoiceNumber:     invNum,
			Amount:            amount,
			DueDate:           dueDate,
			Status:            "Unpaid",
		}
		
		if err := s.db.Create(&invoice).Error; err == nil {
			count++
		} else {
			fmt.Printf("Error creating invoice for student %d: %v\n", student.ID, err)
		}
	}

	return count, nil
}

func (s *financeService) UpdateInvoice(id uint, data map[string]interface{}) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := s.db.First(&invoice, id).Error; err != nil {
		return nil, err
	}

	// Update fields if present in data
	if val, ok := data["amount"]; ok { invoice.Amount = getFloat(val) }
	if val, ok := data["due_date"]; ok {
		if d, err := time.Parse("2006-01-02", val.(string)); err == nil {
			invoice.DueDate = d
		}
	}
	if val, ok := data["status"]; ok { invoice.Status = val.(string) }

	if err := s.db.Save(&invoice).Error; err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (s *financeService) DeleteInvoice(id uint) error {
	return s.db.Delete(&models.Invoice{}, id).Error
}

func (s *financeService) FindStudentByNIS(nis string) (*models.Student, error) {
	var student models.Student
	if err := s.db.Where("nis = ?", nis).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

// Helpers to get types from map interface (could be float64 from JSON or string)
func getUint(v interface{}) uint {
	switch t := v.(type) {
	case float64:
		return uint(t)
	case string:
		var u uint
		fmt.Sscanf(t, "%d", &u)
		return u
	}
	return 0
}

func getFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case string:
		var f float64
		fmt.Sscanf(t, "%f", &f)
		return f
	}
	return 0
}
