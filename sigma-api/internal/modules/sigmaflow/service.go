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
	GetAllPayments() ([]map[string]interface{}, error)
	ProcessManualPayment(invoiceID uint, amount float64, method string, adminName string) error
	GetAllCashFlows(filters map[string]string) ([]map[string]interface{}, error)
	CreateCashFlowEntry(data map[string]interface{}, adminName string) (*models.CashFlow, error)
	DeleteCashFlowEntry(id uint) error
	GetCashFlowSummary() (map[string]interface{}, error)
	SyncExistingPayments() (int, error)
	GetProfitLossReport(month, year int) (map[string]interface{}, error)
	GetArrearsReport() ([]map[string]interface{}, error)
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
	err := s.db.Preload("PaymentCategory").Where("student_id = ?", studentID).Order("created_at desc").Find(&invoices).Error
	if err != nil {
		return nil, err
	}

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

	if invoice.PaymentURL != nil && invoice.PaymentStatus == "PENDING" {
		return &invoice, nil
	}

	xenditReq := utils.XenditInvoiceRequest{
		ExternalID:      fmt.Sprintf("%d-%d", invoice.ID, time.Now().Unix()),
		Amount:          invoice.Amount,
		PayerEmail:      "pondok.sigma@example.com",
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

func (s *financeService) ProcessPaymentCallback(externalID string, amount float64, method string) error {
	var invoiceID uint
	fmt.Sscanf(externalID, "%d-", &invoiceID)

	return s.db.Transaction(func(tx *gorm.DB) error {
		var invoice models.Invoice
		if err := tx.Preload("Student").Preload("PaymentCategory").First(&invoice, invoiceID).Error; err != nil {
			return err
		}

		if invoice.Status == "Paid" {
			return nil
		}

		payment := models.Payment{
			InvoiceID:   invoice.ID,
			PaymentDate: time.Now(),
			AmountPaid:  amount,
			Method:      method,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		cashFlow := models.CashFlow{
			Type:        models.CashFlowIncome,
			Category:    invoice.PaymentCategory.Name,
			Amount:      amount,
			Date:        time.Now(),
			Description: fmt.Sprintf("Pembayaran %s - %s", invoice.PaymentCategory.Name, invoice.Student.Name),
			ReferenceID: &payment.ID,
			CreatedBy:   "SYSTEM/XENDIT",
		}
		if err := tx.Create(&cashFlow).Error; err != nil {
			return err
		}

		invoice.Status = "Paid" 
		invoice.PaymentStatus = "PAID"
		return tx.Save(&invoice).Error
	})
}

func (s *financeService) GetGlobalFinanceStats() (map[string]interface{}, error) {
	var totalIncome, totalExpense float64
	s.db.Model(&models.CashFlow{}).Where("type = ?", "income").Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)
	s.db.Model(&models.CashFlow{}).Where("type = ?", "expense").Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense)
	
	var totalInvoices int64
	s.db.Model(&models.Invoice{}).Count(&totalInvoices)
	
	var pendingAmount float64
	s.db.Model(&models.Invoice{}).Where("status != ?", "Paid").Select("COALESCE(SUM(amount), 0)").Scan(&pendingAmount)

	var paidCount, unpaidCount, partialCount int64
	s.db.Model(&models.Invoice{}).Where("status = ?", "Paid").Count(&paidCount)
	s.db.Model(&models.Invoice{}).Where("status = ?", "Unpaid").Count(&unpaidCount)
	s.db.Model(&models.Invoice{}).Where("status = ?", "Partial").Count(&partialCount)

	var recentPayments []models.Payment
	s.db.Preload("Invoice.Student").Preload("Invoice.PaymentCategory").Order("payment_date desc").Limit(5).Find(&recentPayments)

	// Calculate Growth
	now := time.Now()
	thisMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	lastMonthStart := thisMonthStart.AddDate(0, -1, 0)

	var thisMonthIncome, lastMonthIncome float64
	s.db.Model(&models.CashFlow{}).Where("type = ? AND date >= ?", models.CashFlowIncome, thisMonthStart).Select("COALESCE(SUM(amount), 0)").Scan(&thisMonthIncome)
	s.db.Model(&models.CashFlow{}).Where("type = ? AND date >= ? AND date < ?", models.CashFlowIncome, lastMonthStart, thisMonthStart).Select("COALESCE(SUM(amount), 0)").Scan(&lastMonthIncome)

	growth := 0.0
	if lastMonthIncome > 0 {
		growth = ((thisMonthIncome - lastMonthIncome) / lastMonthIncome) * 100
	} else if thisMonthIncome > 0 {
		growth = 100.0
	}

	formattedPayments := make([]map[string]interface{}, len(recentPayments))
	for i, p := range recentPayments {
		studentName, category := "Unknown", "General"
		if p.Invoice != nil {
			if p.Invoice.Student != nil { studentName = p.Invoice.Student.Name }
			if p.Invoice.PaymentCategory != nil { category = p.Invoice.PaymentCategory.Name }
		}
		formattedPayments[i] = map[string]interface{}{
			"id": p.ID, "student_name": studentName, "category": category, "amount": p.AmountPaid,
			"amount_formatted": utils.FormatRupiah(p.AmountPaid), "date": p.PaymentDate.Format("02 Jan 15:04"),
		}
	}

	return map[string]interface{}{
		"total_income":       totalIncome,
		"income_formatted":   utils.FormatRupiah(totalIncome),
		"income_growth":      fmt.Sprintf("%.1f%%", growth),
		"total_expense":      totalExpense,
		"expense_formatted":  utils.FormatRupiah(totalExpense),
		"current_balance":    totalIncome - totalExpense,
		"balance_formatted":  utils.FormatRupiah(totalIncome - totalExpense),
		"pending_amount":     pendingAmount,
		"pending_formatted":  utils.FormatRupiah(pendingAmount),
		"total_invoices":     totalInvoices,
		"status_distribution": map[string]int64{
			"paid":    paidCount,
			"unpaid":  unpaidCount,
			"partial": partialCount,
		},
		"recent_payments":    formattedPayments,
		"top_arrears":        s.getTopArrears(5),
	}, nil
}

func (s *financeService) getTopArrears(limit int) []map[string]interface{} {
	var invoices []models.Invoice
	s.db.Preload("Student").Preload("PaymentCategory").Where("status != ?", "Paid").Limit(limit).Find(&invoices)

	result := make([]map[string]interface{}, 0)
	for _, inv := range invoices {
		var paid float64
		s.db.Model(&models.Payment{}).Where("invoice_id = ?", inv.ID).Select("COALESCE(SUM(amount_paid), 0)").Scan(&paid)
		
		arrears := inv.Amount - paid
		if arrears > 0 {
			result = append(result, map[string]interface{}{
				"student_name": inv.Student.Name,
				"category":     inv.PaymentCategory.Name,
				"arrears":      arrears,
				"formatted":    utils.FormatRupiah(arrears),
			})
		}
	}
	return result
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
		if inv.Student != nil { studentName = inv.Student.Name }
		categoryName := "General"
		if inv.PaymentCategory != nil { categoryName = inv.PaymentCategory.Name }

		var amountPaid float64
		s.db.Model(&models.Payment{}).Where("invoice_id = ?", inv.ID).Select("COALESCE(SUM(amount_paid), 0)").Scan(&amountPaid)

		result = append(result, map[string]interface{}{
			"id":                 inv.ID,
			"invoice_number":     inv.InvoiceNumber,
			"student_name":       studentName,
			"category":           categoryName,
			"amount":             inv.Amount,
			"amount_formatted":   utils.FormatRupiah(inv.Amount),
			"amount_paid":        amountPaid,
			"paid_formatted":     utils.FormatRupiah(amountPaid),
			"remaining":          inv.Amount - amountPaid,
			"remaining_formatted": utils.FormatRupiah(inv.Amount - amountPaid),
			"status":             inv.Status,
			"due_date":           inv.DueDate.Format("02 Jan 2006"),
			"payment_url":        inv.PaymentURL,
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
	cat := models.PaymentCategory{
		Name:        data["name"].(string),
		Amount:      data["amount"].(float64),
		IsRecurring: data["is_recurring"].(bool),
	}
	if err := s.db.Create(&cat).Error; err != nil { return nil, err }
	return &cat, nil
}

func (s *financeService) UpdateCategory(id uint, data map[string]interface{}) (*models.PaymentCategory, error) {
	var cat models.PaymentCategory
	if err := s.db.First(&cat, id).Error; err != nil { return nil, err }
	if name, ok := data["name"].(string); ok { cat.Name = name }
	if amount, ok := data["amount"].(float64); ok { cat.Amount = amount }
	if isRecurring, ok := data["is_recurring"].(bool); ok { cat.IsRecurring = isRecurring }
	if err := s.db.Save(&cat).Error; err != nil { return nil, err }
	return &cat, nil
}

func (s *financeService) DeleteCategory(id uint) error {
	return s.db.Delete(&models.PaymentCategory{}, id).Error
}

func (s *financeService) CreateInvoice(data map[string]interface{}) (*models.Invoice, error) {
	invNum := fmt.Sprintf("INV/%s/%04d", time.Now().Format("20060102"), time.Now().Nanosecond()%10000)
	dueDate, _ := time.Parse("2006-01-02", data["due_date"].(string))
	invoice := models.Invoice{
		StudentID:         getUint(data["student_id"]),
		PaymentCategoryID: getUint(data["payment_category_id"]),
		InvoiceNumber:     invNum,
		Amount:            getFloat(data["amount"]),
		DueDate:           dueDate,
		Status:            "Unpaid",
	}
	if err := s.db.Create(&invoice).Error; err != nil { return nil, err }
	return &invoice, nil
}

func (s *financeService) CreateBulkInvoices(data map[string]interface{}) (int, error) {
	cID := getUint(data["payment_category_id"])
	amount := getFloat(data["amount"])
	dueDate, _ := time.Parse("2006-01-02", data["due_date"].(string))
	target := data["target"].(string)
	
	var students []models.Student
	query := s.db.Model(&models.Student{})
	if target == "classroom" {
		query = query.Where("classroom_id = ?", getUint(data["classroom_id"]))
	}
	if err := query.Find(&students).Error; err != nil { return 0, err }

	count := 0
	now := time.Now()
	for i, student := range students {
		invNum := fmt.Sprintf("INV/%s/%d%d%d", now.Format("20060102"), student.ID, now.Unix()%1000, i)
		invoice := models.Invoice{
			StudentID:         student.ID,
			PaymentCategoryID: cID,
			InvoiceNumber:     invNum,
			Amount:            amount,
			DueDate:           dueDate,
			Status:            "Unpaid",
		}
		if err := s.db.Create(&invoice).Error; err == nil { count++ }
	}
	return count, nil
}

func (s *financeService) UpdateInvoice(id uint, data map[string]interface{}) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := s.db.First(&invoice, id).Error; err != nil { return nil, err }
	if val, ok := data["amount"]; ok { invoice.Amount = getFloat(val) }
	if val, ok := data["due_date"]; ok {
		if d, err := time.Parse("2006-01-02", val.(string)); err == nil { invoice.DueDate = d }
	}
	if val, ok := data["status"]; ok { invoice.Status = val.(string) }
	if err := s.db.Save(&invoice).Error; err != nil { return nil, err }
	return &invoice, nil
}

func (s *financeService) DeleteInvoice(id uint) error {
	return s.db.Delete(&models.Invoice{}, id).Error
}

func (s *financeService) FindStudentByNIS(nis string) (*models.Student, error) {
	var student models.Student
	if err := s.db.Where("nis = ?", nis).First(&student).Error; err != nil { return nil, err }
	return &student, nil
}

func (s *financeService) GetAllPayments() ([]map[string]interface{}, error) {
	var payments []models.Payment
	err := s.db.Preload("Invoice.Student").Preload("Invoice.PaymentCategory").Order("payment_date desc").Find(&payments).Error
	if err != nil { return nil, err }

	result := make([]map[string]interface{}, len(payments))
	for i, p := range payments {
		status := "SUCCESS"
		if p.Invoice != nil && p.Invoice.PaymentStatus != "" { status = p.Invoice.PaymentStatus }
		studentName, categoryName := "Unknown", "Unknown"
		if p.Invoice != nil {
			if p.Invoice.Student != nil { studentName = p.Invoice.Student.Name }
			if p.Invoice.PaymentCategory != nil { categoryName = p.Invoice.PaymentCategory.Name }
		}
		result[i] = map[string]interface{}{
			"id": p.ID, "student_name": studentName, "category": categoryName, "amount": p.AmountPaid,
			"amount_formatted": utils.FormatRupiah(p.AmountPaid), "date": p.PaymentDate.Format("02 Jan 2006, 15:04"),
			"method": p.Method, "status": status,
		}
	}
	return result, nil
}

func (s *financeService) ProcessManualPayment(invoiceID uint, amount float64, method string, adminName string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var invoice models.Invoice
		if err := tx.Preload("Payments").Preload("Student").Preload("PaymentCategory").First(&invoice, invoiceID).Error; err != nil { return err }
		if invoice.Status == "Paid" { return errors.New("faktur sudah lunas") }

		payment := models.Payment{ InvoiceID: invoice.ID, PaymentDate: time.Now(), AmountPaid: amount, Method: method }
		if err := tx.Create(&payment).Error; err != nil { return err }

		cashFlow := models.CashFlow{
			Type: models.CashFlowIncome, Category: invoice.PaymentCategory.Name, Amount: amount, Date: time.Now(),
			Description: fmt.Sprintf("Pembayaran %s - %s", invoice.PaymentCategory.Name, invoice.Student.Name),
			ReferenceID: &payment.ID, CreatedBy: adminName,
		}
		if err := tx.Create(&cashFlow).Error; err != nil { return err }

		var totalPaid float64
		for _, p := range invoice.Payments { totalPaid += p.AmountPaid }
		totalPaid += amount

		if totalPaid >= invoice.Amount {
			invoice.Status = "Paid"
			invoice.PaymentStatus = "PAID"
		} else {
			invoice.Status = "Partial"
			invoice.PaymentStatus = "PARTIAL"
		}
		return tx.Save(&invoice).Error
	})
}

func (s *financeService) GetAllCashFlows(filters map[string]string) ([]map[string]interface{}, error) {
	var entries []models.CashFlow
	query := s.db.Order("date desc")
	if t, ok := filters["type"]; ok && t != "" { query = query.Where("type = ?", t) }
	if start, ok := filters["start_date"]; ok && start != "" { query = query.Where("date >= ?", start) }
	if end, ok := filters["end_date"]; ok && end != "" { query = query.Where("date <= ?", end) }
	if err := query.Find(&entries).Error; err != nil { return nil, err }

	result := make([]map[string]interface{}, len(entries))
	for i, e := range entries {
		result[i] = map[string]interface{}{
			"id": e.ID, "type": e.Type, "category": e.Category, "amount": e.Amount,
			"amount_formatted": utils.FormatRupiah(e.Amount), "date": e.Date.Format("2006-01-02 15:04"),
			"description": e.Description, "created_by": e.CreatedBy,
		}
	}
	return result, nil
}

func (s *financeService) CreateCashFlowEntry(data map[string]interface{}, adminName string) (*models.CashFlow, error) {
	date := time.Now()
	if dStr, ok := data["date"].(string); ok && dStr != "" {
		if d, err := time.Parse("2006-01-02", dStr); err == nil { date = d }
	}
	entry := models.CashFlow{
		Type: models.CashFlowType(data["type"].(string)), Category: data["category"].(string),
		Amount: getFloat(data["amount"]), Date: date, Description: data["description"].(string), CreatedBy: adminName,
	}
	if err := s.db.Create(&entry).Error; err != nil { return nil, err }
	return &entry, nil
}

func (s *financeService) DeleteCashFlowEntry(id uint) error {
	return s.db.Delete(&models.CashFlow{}, id).Error
}

func (s *financeService) GetCashFlowSummary() (map[string]interface{}, error) {
	var totalIncome, totalExpense float64
	s.db.Model(&models.CashFlow{}).Where("type = ?", "income").Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)
	s.db.Model(&models.CashFlow{}).Where("type = ?", "expense").Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense)
	balance := totalIncome - totalExpense
	return map[string]interface{}{
		"total_income": totalIncome, "income_formatted": utils.FormatRupiah(totalIncome),
		"total_expense": totalExpense, "expense_formatted": utils.FormatRupiah(totalExpense),
		"balance": balance, "balance_formatted": utils.FormatRupiah(balance),
	}, nil
}

func (s *financeService) GetProfitLossReport(month, year int) (map[string]interface{}, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0)

	var incomeEntries []models.CashFlow
	var expenseEntries []models.CashFlow

	s.db.Where("type = ? AND date >= ? AND date < ?", models.CashFlowIncome, startDate, endDate).Find(&incomeEntries)
	s.db.Where("type = ? AND date >= ? AND date < ?", models.CashFlowExpense, startDate, endDate).Find(&expenseEntries)

	incomeByCategory := make(map[string]float64)
	totalIncome := 0.0
	for _, e := range incomeEntries {
		incomeByCategory[e.Category] += e.Amount
		totalIncome += e.Amount
	}

	expenseByCategory := make(map[string]float64)
	totalExpense := 0.0
	for _, e := range expenseEntries {
		expenseByCategory[e.Category] += e.Amount
		totalExpense += e.Amount
	}

	return map[string]interface{}{
		"income_categories":  incomeByCategory,
		"expense_categories": expenseByCategory,
		"total_income":       totalIncome,
		"income_formatted":   utils.FormatRupiah(totalIncome),
		"total_expense":      totalExpense,
		"expense_formatted":  utils.FormatRupiah(totalExpense),
		"net_profit":         totalIncome - totalExpense,
		"profit_formatted":   utils.FormatRupiah(totalIncome - totalExpense),
		"month":              startDate.Format("January"),
		"year":               year,
	}, nil
}

func (s *financeService) GetArrearsReport() ([]map[string]interface{}, error) {
	var invoices []models.Invoice
	err := s.db.Preload("Student").Preload("PaymentCategory").Where("status != ?", "Paid").Find(&invoices).Error
	if err != nil { return nil, err }

	result := make([]map[string]interface{}, 0)
	for _, inv := range invoices {
		var paid float64
		s.db.Model(&models.Payment{}).Where("invoice_id = ?", inv.ID).Select("COALESCE(SUM(amount_paid), 0)").Scan(&paid)

		result = append(result, map[string]interface{}{
			"id":             inv.ID,
			"student_name":   inv.Student.Name,
			"nis":            inv.Student.NIS,
			"category":       inv.PaymentCategory.Name,
			"total_amount":   inv.Amount,
			"total_formatted": utils.FormatRupiah(inv.Amount),
			"paid":           paid,
			"paid_formatted":  utils.FormatRupiah(paid),
			"arrears":        inv.Amount - paid,
			"arrears_formatted": utils.FormatRupiah(inv.Amount - paid),
			"due_date":       inv.DueDate.Format("02 Jan 2006"),
			"status":         inv.Status,
		})
	}
	return result, nil
}

func (s *financeService) SyncExistingPayments() (int, error) {
	var payments []models.Payment
	err := s.db.Preload("Invoice.Student").Preload("Invoice.PaymentCategory").Find(&payments).Error
	if err != nil { return 0, err }

	count := 0
	for _, p := range payments {
		var exists int64
		s.db.Model(&models.CashFlow{}).Where("reference_id = ?", p.ID).Count(&exists)
		if exists == 0 {
			studentName, categoryName := "Unknown", "General"
			if p.Invoice != nil {
				if p.Invoice.Student != nil { studentName = p.Invoice.Student.Name }
				if p.Invoice.PaymentCategory != nil { categoryName = p.Invoice.PaymentCategory.Name }
			}
			
			entry := models.CashFlow{
				Type: models.CashFlowIncome, Category: categoryName, Amount: p.AmountPaid,
				Date: p.PaymentDate, Description: fmt.Sprintf("Pembayaran %s - %s (Sync)", categoryName, studentName),
				ReferenceID: &p.ID, CreatedBy: "SYSTEM/SYNC",
			}
			if err := s.db.Create(&entry).Error; err == nil { count++ }
		}
	}
	return count, nil
}

func getUint(v interface{}) uint {
	switch t := v.(type) {
	case float64: return uint(t)
	case string:
		var u uint
		fmt.Sscanf(t, "%d", &u)
		return u
	}
	return 0
}

func getFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64: return t
	case string:
		var f float64
		fmt.Sscanf(t, "%f", &f)
		return f
	}
	return 0
}
