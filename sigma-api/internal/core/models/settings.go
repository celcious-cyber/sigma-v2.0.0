package models

import (
	"time"

	"gorm.io/gorm"
)

// Book represents a library resource
type Book struct {
	gorm.Model
	ISBN      *string `gorm:"uniqueIndex" json:"isbn"`
	Title     string  `gorm:"not null" json:"title"`
	Author    string  `json:"author"`
	Publisher string  `json:"publisher"`
	Year      int     `json:"year"`
	Category  string  `json:"category"`
	Stock     int     `gorm:"default:0" json:"stock"`
}

// Borrowing represents a book loan record
type Borrowing struct {
	gorm.Model
	StudentID  uint       `gorm:"not null" json:"student_id"`
	Student    *Student   `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	BookID     uint       `gorm:"not null" json:"book_id"`
	Book       *Book      `gorm:"foreignKey:BookID" json:"book,omitempty"`
	BorrowDate time.Time  `json:"borrow_date"`
	DueDate    time.Time  `json:"due_date"`
	ReturnDate *time.Time `json:"return_date"`
	Status     string     `gorm:"default:'Borrowed'" json:"status"` // Borrowed, Returned, Overdue
}

// Setting represents global application configuration
type Setting struct {
	gorm.Model
	SiteName        *string `json:"site_name"`
	Tagline         *string `json:"tagline"`
	About           *string `gorm:"type:text" json:"about"`
	Logo            *string `json:"logo"`
	Favicon         *string `json:"favicon"`
	Address         *string `gorm:"type:text" json:"address"`
	Email           *string `json:"email"`
	Phone           *string `json:"phone"`
	WhatsApp        *string `json:"whatsapp"`
	Instagram       *string `json:"instagram"`
	Facebook        *string `json:"facebook"`
	YouTube         *string `json:"youtube"`
	CountStudents   *string `json:"count_students"`
	CountTeachers   *string `json:"count_teachers"`
	CountAlumni     *string `json:"count_alumni"`
	CountFacilities *string `json:"count_facilities"`
}
