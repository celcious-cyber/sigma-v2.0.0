package models

import (
	"time"

	"gorm.io/gorm"
)

// IncomingLetter represents letters received by the secretariat
type IncomingLetter struct {
	gorm.Model
	Number      *string    `json:"number"`
	Sender      string     `gorm:"not null" json:"sender"`
	Address     *string    `gorm:"type:text" json:"address"`
	Subject     string     `gorm:"not null" json:"subject"`
	Date        time.Time  `json:"date"`
	ReceivedAt  time.Time  `json:"received_at"`
	FilePath    *string    `json:"file_path"`
	Description *string    `gorm:"type:text" json:"description"`
}

// OutgoingLetter represents letters sent by the secretariat
type OutgoingLetter struct {
	gorm.Model
	Number      *string    `json:"number"`
	Recipient   string     `gorm:"not null" json:"recipient"`
	Address     *string    `gorm:"type:text" json:"address"`
	Subject     string     `gorm:"not null" json:"subject"`
	Date        time.Time  `json:"date"`
	FilePath    *string    `json:"file_path"`
	Description *string    `gorm:"type:text" json:"description"`
}

// LetterCode represents classification codes for letters
type LetterCode struct {
	gorm.Model
	Code        string `gorm:"uniqueIndex;not null" json:"code"`
	Description string `gorm:"type:text" json:"description"`
}

// Visitor represents a guest entry in the guestbook
type Visitor struct {
	gorm.Model
	Name      string     `gorm:"not null" json:"name"`
	Phone     *string    `json:"phone"`
	Purpose   string     `gorm:"not null" json:"purpose"`
	StudentID *uint      `json:"student_id"`
	Student   *Student   `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	CheckIn   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"check_in"`
	CheckOut  *time.Time `json:"check_out"`
	Status    string     `gorm:"default:'Di Lokasi'" json:"status"` // Di Lokasi, Selesai
}

// Announcement represents a global system/school notification
type Announcement struct {
	gorm.Model
	Title       string     `gorm:"not null" json:"title"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	Category    string     `json:"category"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedBy   uint       `json:"created_by"`
	Author      *User      `gorm:"foreignKey:CreatedBy" json:"author,omitempty"`
}
