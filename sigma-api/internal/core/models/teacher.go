package models

import (
	"time"

	"gorm.io/gorm"
)

// Teacher represents a teacher/musyrif entity
type Teacher struct {
	gorm.Model
	Name           string     `gorm:"type:string;not null" json:"name"`
	Gender         string     `gorm:"type:string;not null;default:'L'" json:"gender"`
	NIK            string     `gorm:"type:string;uniqueIndex;not null" json:"nik"`
	BirthPlace     *string    `json:"birth_place"`
	BirthDate      *time.Time `json:"birth_date"`
	NIP            *string    `gorm:"type:string;uniqueIndex" json:"nip"`
	Position       *string    `json:"position"`
	Status         *string    `json:"status"`
	Education      *string    `json:"education"`
	GraduationYear *string    `json:"graduation_year"`
	Email          *string    `json:"email"`
	Address        *string    `gorm:"type:text" json:"address"`
	WhatsApp       *string    `json:"whatsapp"`
	Photo          *string    `json:"photo"`
	EntryYear      *string    `json:"entry_year"`
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	UserID         *uint      `json:"user_id"`
	User           *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName sets the singular table name for Teacher
func (Teacher) TableName() string {
	return "teacher"
}
