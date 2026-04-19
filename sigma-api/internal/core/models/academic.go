package models

import (
	"time"

	"gorm.io/gorm"
)

// Subject represents a course/mata pelajaran
type Subject struct {
	gorm.Model
	Name string `gorm:"not null" json:"name"`
	Code string `gorm:"uniqueIndex;not null" json:"code"`
}

// TeacherSchedule represents a timetable for a teacher
type TeacherSchedule struct {
	gorm.Model
	UserID      uint       `gorm:"not null" json:"user_id"` // The Teacher
	User        *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ClassroomID uint       `gorm:"not null" json:"classroom_id"`
	Classroom   *Classroom `gorm:"foreignKey:ClassroomID" json:"classroom,omitempty"`
	SubjectID   uint       `gorm:"not null" json:"subject_id"`
	Subject     *Subject   `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	DayOfWeek   int        `json:"day_of_week"` // 1-7
	StartTime   string     `json:"start_time"`   // e.g., "07:00"
	EndTime     string     `json:"end_time"`     // e.g., "08:30"
}

// Attendance represents daily or subject-based student presence
type Attendance struct {
	gorm.Model
	StudentID   uint       `gorm:"not null" json:"student_id"`
	Student     *Student   `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	ClassroomID uint       `gorm:"not null" json:"classroom_id"`
	Classroom   *Classroom `gorm:"foreignKey:ClassroomID" json:"classroom,omitempty"`
	Date        time.Time  `json:"date"`
	Status      string     `json:"status"` // Hadir, Izin, Sakit, Alpa
	RecordedBy  uint       `json:"recorded_by"`
	Author      *User      `gorm:"foreignKey:RecordedBy" json:"author,omitempty"`
}

// Assessment represents grades/evaluations for a student
type Assessment struct {
	gorm.Model
	StudentID    uint     `gorm:"not null" json:"student_id"`
	Student      *Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	ClassroomID  uint     `gorm:"not null" json:"classroom_id"`
	Classroom    *Classroom `gorm:"foreignKey:ClassroomID" json:"classroom,omitempty"`
	TeacherID    uint     `gorm:"not null" json:"teacher_id"`
	Author       *User    `gorm:"foreignKey:TeacherID" json:"author,omitempty"`
	SubjectID    *uint    `json:"subject_id"`
	Subject      *Subject `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Type         string   `json:"type"` // academic, attitude, etc.
	Score        *float64 `gorm:"type:decimal(5,2)" json:"score"`
	Grade        *string  `json:"grade"`
	Remarks      *string  `gorm:"type:text" json:"remarks"`
	Term         string   `json:"term"`          // Ganjil, Genap
	AcademicYear string   `json:"academic_year"` // 2024/2025
}

// StudyPeriod represents an academic semester or specific period
type StudyPeriod struct {
	gorm.Model
	Name      string `json:"name"` // e.g., "Semester 1"
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// AcademicCalendar represents a school event or schedule item
type AcademicCalendar struct {
	gorm.Model
	Title       string     `gorm:"not null" json:"title"`
	Description *string    `gorm:"type:text" json:"description"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Type        string     `gorm:"default:'Kegiatan'" json:"type"` // Libur, Ujian, Kegiatan
	Color       *string    `json:"color"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
}

// QuranMemorization represents Tahfidz progress
type QuranMemorization struct {
	gorm.Model
	StudentID  uint      `gorm:"not null" json:"student_id"`
	Student    *Student  `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	TeacherID  uint      `gorm:"not null" json:"teacher_id"`
	Author     *User     `gorm:"foreignKey:TeacherID" json:"author,omitempty"`
	SurahName  string    `json:"surah_name"`
	VerseStart *int      `json:"verse_start"`
	VerseEnd   *int      `json:"verse_end"`
	Type       string    `json:"type"` // Sabaq, Sabqi, etc.
	Grade      *string   `json:"grade"`
	Remarks    *string   `gorm:"type:text" json:"remarks"`
	Date       time.Time `json:"date"`
}
