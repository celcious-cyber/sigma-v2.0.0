package sigmaguard

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"sigma-api/internal/core/models"
	"sigma-api/internal/modules/sigmabase"
)

type GuardService interface {
	GetViolations(gender string) ([]models.Violation, error)
	RecordViolation(v models.Violation) error
	GetPermits(gender string) ([]models.Permit, error)
	CreatePermit(p models.Permit) error
}

type guardService struct {
	db         *gorm.DB
	studentSvc sigmabase.StudentService
}

func NewGuardService(db *gorm.DB, s sigmabase.StudentService) GuardService {
	return &guardService{db, s}
}

// GetViolations returns violations filtered by gender (Putra/Putri split)
func (s *guardService) GetViolations(gender string) ([]models.Violation, error) {
	var violations []models.Violation
	query := s.db.Preload("Student").Preload("ViolationDictionary")
	if gender != "" {
		query = query.Joins("JOIN students ON students.id = violations.student_id").Where("students.gender = ?", gender)
	}
	err := query.Order("created_at desc").Find(&violations).Error
	return violations, err
}

// RecordViolation saves a violation and refreshes the student's total points
func (s *guardService) RecordViolation(v models.Violation) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&v).Error; err != nil {
			return err
		}
		// Siksa backend phase: Automatic point refresh
		// We call the service method or just do it here
		var total int64
		tx.Model(&models.Violation{}).Where("student_id = ?", v.StudentID).Select("sum(points)").Row().Scan(&total)
		return tx.Model(&models.Student{}).Where("id = ?", v.StudentID).Update("total_violation_points", total).Error
	})
	return err
}

// GetPermits returns permits filtered by gender
func (s *guardService) GetPermits(gender string) ([]models.Permit, error) {
	var permits []models.Permit
	query := s.db.Preload("Student")
	if gender != "" {
		query = query.Joins("JOIN students ON students.id = permits.student_id").Where("students.gender = ?", gender)
	}
	err := query.Order("created_at desc").Find(&permits).Error
	return permits, err
}

// CreatePermit handles permit creation with quota validation
func (s *guardService) CreatePermit(p models.Permit) error {
	if !p.IsEmergency {
		isBlocked, _, err := s.studentSvc.IsPermitBlocked(p.StudentID)
		if err != nil {
			return err
		}
		if isBlocked {
			return errors.New("student has exceeded permit quota for this period")
		}
	}
	
	p.CreatedAt = time.Now()
	return s.db.Create(&p).Error
}
