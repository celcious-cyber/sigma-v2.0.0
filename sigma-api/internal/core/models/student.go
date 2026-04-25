package models

import (
	"time"

	"gorm.io/gorm"
)

// Student represents the santri/siswa entity
type Student struct {
	gorm.Model
	ClassroomID *uint      `json:"classroom_id"`
	Classroom   *Classroom `gorm:"foreignKey:ClassroomID" json:"classroom,omitempty"`
	NIS         string     `gorm:"type:string;uniqueIndex;not null" json:"nis"`
	NISN        *string    `gorm:"type:string;uniqueIndex" json:"nisn"`
	NIK         *string    `gorm:"type:string;uniqueIndex" json:"nik"`
	Name        string     `gorm:"type:string;not null" json:"name"`
	BirthPlace  *string    `json:"birth_place"`
	BirthDate   *time.Time `json:"birth_date"`
	Gender      string     `gorm:"type:string;default:'L'" json:"gender"` // L = Laki-laki, P = Perempuan
	EntryYear   string     `json:"entry_year"` // The batch/enrollment year (e.g., "2024")
	Class       string     `json:"class"` // The class string representation (e.g., "7A")
	ParentName  *string    `json:"parent_name"`
	ParentPhone *string    `json:"parent_phone"`
	Phone       *string    `json:"phone"`
	Address                *string    `gorm:"type:text" json:"address"`
	TotalViolationPoints   int        `gorm:"default:0" json:"total_violation_points"` // Added for discipline summary
	Password               *string    `gorm:"type:string" json:"-"` // Added for Sigmagate login
	IsActive               bool       `gorm:"default:true" json:"is_active"` // Added for archiving (graduation)
	AvatarURL              *string    `json:"avatar_url,omitempty"`
}
