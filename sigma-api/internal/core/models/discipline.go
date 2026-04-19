package models

import (
	"time"

	"gorm.io/gorm"
)

// ViolationDictionary represents the master list of possible rules broken
type ViolationDictionary struct {
	gorm.Model
	Name                string `gorm:"not null" json:"name"`
	Severity            string `json:"severity"` // Ringan, Sedang, Berat
	Gender              string `gorm:"default:'L'" json:"gender"` // L = Putra, P = Putri
	Points              int    `json:"points"`
	RecommendedSanction string `json:"recommended_sanction"`
}

// Violation represents an actual rule break by a student
type Violation struct {
	gorm.Model
	StudentID     uint      `gorm:"not null" json:"student_id"`
	Student       *Student  `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Class         string    `json:"class"`
	ViolationType string    `json:"violation_type"`
	Points        int       `gorm:"default:0" json:"points"`
	Date          time.Time `json:"date"`
}

// Sanction represents a punishment given for violations
type Sanction struct {
	gorm.Model
	StudentID     uint     `gorm:"not null" json:"student_id"`
	Student       *Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Class         string   `json:"class"`
	ViolationType string   `json:"violation_type"`
	SanctionType  string   `json:"sanction_type"`
	Duration      *string  `json:"duration"`
	Status        string   `gorm:"default:'Pending'" json:"status"` // Pending, Proses, Selesai
}

// Permit represents a student leaving the school/boarding area
type Permit struct {
	gorm.Model
	StudentID   uint       `gorm:"not null" json:"student_id"`
	Student     *Student   `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Class       string     `json:"class"`
	Purpose     string     `json:"purpose"`
	Type        string     `gorm:"default:'Keluar'" json:"type"` // Pulang, Keluar
	IsEmergency bool       `gorm:"default:false" json:"is_emergency"`
	OutTime     time.Time  `json:"out_time"`
	ReturnTime  *time.Time `json:"return_time"`
	Status      string     `gorm:"default:'Diluar'" json:"status"` // Diluar, Kembali, Terlambat
}

// CareRecord represents a behavioral or guidance counseling session or medical visit
type CareRecord struct {
	gorm.Model
	StudentID      uint     `gorm:"not null" json:"student_id"`
	Student        *Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	DiseaseID      *uint    `json:"disease_id"`
	Disease        *Disease `gorm:"foreignKey:DiseaseID" json:"disease,omitempty"`
	Complaint      *string  `json:"complaint"`
	Diagnosis      *string  `json:"diagnosis"`
	Treatment      *string  `gorm:"type:text" json:"treatment"`
	IsHospitalized bool     `gorm:"default:false" json:"is_hospitalized"`
	IsDone         bool     `gorm:"default:false" json:"is_done"`
	RecordedBy     *uint    `json:"recorded_by"`
	Author         *User    `gorm:"foreignKey:RecordedBy" json:"author,omitempty"`
}
