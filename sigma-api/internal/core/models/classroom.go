package models

import (
	"gorm.io/gorm"
)

// Classroom represents a class entity
type Classroom struct {
	gorm.Model
	Name        string `gorm:"type:string;uniqueIndex;not null" json:"name"`
	Level       string `json:"level"`
	Type        string `gorm:"default:'KMI'" json:"type"`
	Gender      string `gorm:"default:'Umum'" json:"gender"`
	WaliKelasID *uint  `json:"wali_kelas_id"`
	WaliKelas   *User  `gorm:"foreignKey:WaliKelasID" json:"wali_kelas,omitempty"`
}
