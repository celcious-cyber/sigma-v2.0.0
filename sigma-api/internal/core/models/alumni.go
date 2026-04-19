package models

import (
	"time"

	"gorm.io/gorm"
)

// Alumni represents a former student
type Alumni struct {
	gorm.Model
	Name           string     `gorm:"type:string;not null" json:"name"`
	Gender         string     `gorm:"type:string;not null;default:'L'" json:"gender"`
	BirthPlace     *string    `json:"birth_place"`
	BirthDate      *time.Time `json:"birth_date"`
	NIK            *string    `gorm:"type:string;uniqueIndex" json:"nik"`
	NIS            *string    `gorm:"type:string;uniqueIndex" json:"nis"`
	GraduationYear string     `json:"graduation_year"`
	Batch          string     `json:"batch"`
	Address        *string    `gorm:"type:text" json:"address"`
	WhatsApp       *string    `json:"whatsapp"`
	Email          *string    `json:"email"`
	Photo          *string    `json:"photo"`
	ServiceStatus  string     `gorm:"type:string;default:'Tidak Mengabdi'" json:"service_status"`
}
// TableName sets the singular table name for Alumni
func (Alumni) TableName() string {
	return "alumni"
}
