package models

import "gorm.io/gorm"

// InstitutionalUnit represents a division within the Al-Hikmah organization (e.g., SMP, MA, KMI)
type InstitutionalUnit struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Code        string `gorm:"uniqueIndex;not null" json:"code"`
	Description string `gorm:"type:text" json:"description"`
	Status      string `gorm:"default:'Active'" json:"status"` // Active, Inactive
}
