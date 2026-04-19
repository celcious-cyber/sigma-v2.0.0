package sigmacare

import (
	"gorm.io/gorm"
	"sigma-api/internal/core/models"
)

type CareService interface {
	GetCareRecords(gender string) ([]models.CareRecord, error)
	GetFitnessRecords(gender string) ([]models.FitnessRecord, error)
	RecordCare(record models.CareRecord) error
}

type careService struct {
	db *gorm.DB
}

func NewCareService(db *gorm.DB) CareService {
	return &careService{db}
}

func (s *careService) GetCareRecords(gender string) ([]models.CareRecord, error) {
	var records []models.CareRecord
	query := s.db.Preload("Student").Preload("Disease")
	if gender != "" {
		query = query.Joins("JOIN students ON students.id = care_records.student_id").Where("students.gender = ?", gender)
	}
	err := query.Order("created_at desc").Find(&records).Error
	return records, err
}

func (s *careService) GetFitnessRecords(gender string) ([]models.FitnessRecord, error) {
	var records []models.FitnessRecord
	query := s.db.Preload("Student")
	if gender != "" {
		query = query.Joins("JOIN students ON students.id = fitness_records.student_id").Where("students.gender = ?", gender)
	}
	err := query.Order("date desc").Find(&records).Error
	return records, err
}

func (s *careService) RecordCare(record models.CareRecord) error {
	return s.db.Create(&record).Error
}
