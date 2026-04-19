package models

import (
	"time"

	"gorm.io/gorm"
)

// Medicine represents UKS/Clinic medicine inventory
type Medicine struct {
	gorm.Model
	Name      string     `gorm:"not null" json:"name"`
	Stock     int        `gorm:"default:0" json:"stock"`
	Unit      string     `gorm:"default:'Pcs'" json:"unit"` // Tablet, Botol, Pcs
	ExpiredAt *time.Time `json:"expired_at"`
	Description *string  `gorm:"type:text" json:"description"`
}

// Disease represents a catalog of diseases/illnesses
type Disease struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name"`
	Description *string `gorm:"type:text" json:"description"`
}

// FitnessRecord represents student health metrics (BMIs, etc.)
type FitnessRecord struct {
	gorm.Model
	StudentID         uint     `gorm:"not null" json:"student_id"`
	Student           *Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Height            *float64 `json:"height"`
	Weight            *float64 `json:"weight"`
	Date              time.Time `json:"date"`
	FitnessPercentage int      `gorm:"default:0" json:"fitness_percentage"`
	Status            string   `gorm:"default:'Cukup'" json:"status"` // Rajin, Cukup, Kurang, Malas
	Notes             *string  `gorm:"type:text" json:"notes"`
	RecordedBy        *uint    `json:"recorded_by"`
	Author            *User    `gorm:"foreignKey:RecordedBy" json:"author,omitempty"`
}
