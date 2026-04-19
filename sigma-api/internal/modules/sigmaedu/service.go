package sigmaedu

import (
	"gorm.io/gorm"
	"sigma-api/internal/core/models"
)

type EduService interface {
	GetAssessmentsByStudent(studentID uint) ([]models.Assessment, error)
	GetAttendancesByStudent(studentID uint) ([]models.Attendance, error)
	GetSchedules() ([]models.TeacherSchedule, error)
	RecordAttendance(a models.Attendance) error
}

type eduService struct {
	db *gorm.DB
}

func NewEduService(db *gorm.DB) EduService {
	return &eduService{db}
}

func (s *eduService) GetAssessmentsByStudent(studentID uint) ([]models.Assessment, error) {
	var assessments []models.Assessment
	err := s.db.Preload("Subject").Where("student_id = ?", studentID).Order("date desc").Find(&assessments).Error
	return assessments, err
}

func (s *eduService) GetAttendancesByStudent(studentID uint) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := s.db.Preload("Subject").Where("student_id = ?", studentID).Order("date desc").Find(&attendances).Error
	return attendances, err
}

func (s *eduService) GetSchedules() ([]models.TeacherSchedule, error) {
	var schedules []models.TeacherSchedule
	err := s.db.Preload("Teacher").Preload("Subject").Preload("Classroom").Find(&schedules).Error
	return schedules, err
}

func (s *eduService) RecordAttendance(a models.Attendance) error {
	return s.db.Create(&a).Error
}
