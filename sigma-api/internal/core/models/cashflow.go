package models

import (
	"time"

	"gorm.io/gorm"
)

// CashFlowType defines if the transaction is income or expense
type CashFlowType string

const (
	CashFlowIncome  CashFlowType = "income"
	CashFlowExpense CashFlowType = "expense"
)

// CashFlow represents a single financial transaction record
type CashFlow struct {
	gorm.Model
	Type        CashFlowType `gorm:"type:string;not null" json:"type"`
	Category    string       `gorm:"type:string;not null" json:"category"`
	Amount      float64      `gorm:"type:decimal(12,2);not null" json:"amount"`
	Date        time.Time    `json:"date"`
	Description string       `gorm:"type:text" json:"description"`
	ReferenceID *uint        `json:"reference_id,omitempty"` // Link to Payment ID or other entities
	CreatedBy   string       `json:"created_by"`
}
