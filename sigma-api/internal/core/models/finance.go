package models

import (
	"time"

	"gorm.io/gorm"
)

// PaymentCategory represents a type of bill (e.g., SPP, Uang Makan)
type PaymentCategory struct {
	gorm.Model
	Name        string  `gorm:"type:string;not null" json:"name"`
	Amount      float64 `gorm:"type:decimal(12,2);not null" json:"amount"`
	IsRecurring bool    `gorm:"default:false" json:"is_recurring"`
}

// Invoice represents a billing record for a student
type Invoice struct {
	gorm.Model
	StudentID         uint            `gorm:"not null" json:"student_id"`
	Student           *Student        `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	PaymentCategoryID uint            `gorm:"not null" json:"payment_category_id"`
	PaymentCategory   *PaymentCategory `gorm:"foreignKey:PaymentCategoryID" json:"payment_category,omitempty"`
	InvoiceNumber     string          `gorm:"type:varchar(50);unique;not null" json:"invoice_number"`
	Amount            float64         `gorm:"type:decimal(12,2);not null" json:"amount"`
	DueDate           time.Time       `json:"due_date"`
	Status            string          `gorm:"type:string;default:'Unpaid'" json:"status"` // Unpaid, Partial, Paid
	XenditID          *string         `json:"xendit_id,omitempty"`
	PaymentURL        *string         `gorm:"type:text" json:"payment_url,omitempty"`
	PaymentStatus     string          `gorm:"type:string;default:'PENDING'" json:"payment_status"` // Xendit status
	Payments          []Payment       `gorm:"foreignKey:InvoiceID" json:"payments,omitempty"`
}

// Payment represents an actual payment transaction
type Payment struct {
	gorm.Model
	InvoiceID   uint       `gorm:"not null" json:"invoice_id"`
	PaymentDate time.Time  `json:"payment_date"`
	AmountPaid  float64    `gorm:"type:decimal(12,2);not null" json:"amount_paid"`
	Method      string     `json:"method"` // Transfer, Cash, Xendit
	ProofImage  *string    `json:"proof_image,omitempty"`
	Invoice     *Invoice   `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
}
