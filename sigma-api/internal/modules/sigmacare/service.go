package sigmacare

import (
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
	"sigma-api/internal/core/models"
	"sigma-api/internal/shared/utils"
)

type CareService interface {
	GetClinicalDashboardStats() (map[string]interface{}, error)
	CreateMedicalVisit(visit models.MedicalVisit) error
	GetMedicalVisits() ([]models.MedicalVisit, error)
	GetStudentMedicalHistory(studentID uint) ([]models.MedicalVisit, error)
	
	// MCU
	RecordAnthropometry(record models.AnthropometryRecord) error
	GetAnthropometryRecords() ([]models.AnthropometryRecord, error)
	GetMCUHistory(studentID uint) ([]models.AnthropometryRecord, error)

	// Pharmacy
	GetMedicines() ([]models.Medicine, error)
	AddMedicine(medicine models.Medicine) error
	RecordMedicineMutation(mutation models.MedicineMutation) error

	// Facility
	GetBeds() ([]models.Bed, error)

	// Additional Actions
	CreateCertificate(cert models.MedicalCertificate) error
	CreateObservation(obs models.ObservationRecord) error
	UpdateObservation(id uint, updates map[string]interface{}) error
	DischargeObservation(id uint) error
	GetActiveObservations() ([]models.ObservationRecord, error)
}

type careService struct {
	db *gorm.DB
}

func NewCareService(db *gorm.DB) CareService {
	return &careService{db}
}

func (s *careService) GetClinicalDashboardStats() (map[string]interface{}, error) {
	var totalVisits, todayVisits, inpatientCount int64
	s.db.Model(&models.MedicalVisit{}).Count(&totalVisits)
	s.db.Model(&models.MedicalVisit{}).Where("DATE(visit_date) = ?", time.Now().Format("2006-01-02")).Count(&todayVisits)
	s.db.Model(&models.ObservationRecord{}).Where("status = ?", "Active").Count(&inpatientCount)

	// Disease Stats (Top 5)
	type DiseaseCount struct {
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}
	var topDiseases []DiseaseCount
	s.db.Model(&models.MedicalVisit{}).Select("diagnosis as name, count(*) as count").Group("diagnosis").Order("count desc").Limit(5).Scan(&topDiseases)

	// Pharmacy Stock Alerts
	var lowStockItems []models.Medicine
	s.db.Where("stock < ?", 10).Find(&lowStockItems)

	// Recovery Rate
	var treatedCount int64
	s.db.Model(&models.MedicalVisit{}).Where("status = ?", "Treated").Count(&treatedCount)
	recoveryRate := 0.0
	if totalVisits > 0 {
		recoveryRate = float64(treatedCount) / float64(totalVisits) * 100
	}

	return map[string]interface{}{
		"today_visits":    todayVisits,
		"total_visits":    totalVisits,
		"inpatients":      inpatientCount,
		"recovery_rate":   recoveryRate,
		"top_diseases":    topDiseases,
		"low_stock_count": len(lowStockItems),
		"bed_capacity":    12, // Placeholder
	}, nil
}

func (s *careService) CreateMedicalVisit(visit models.MedicalVisit) error {
	visit.VisitDate = time.Now()
	return s.db.Create(&visit).Error
}

func (s *careService) GetMedicalVisits() ([]models.MedicalVisit, error) {
	var visits []models.MedicalVisit
	err := s.db.Preload("Student").Preload("Medic").Order("visit_date desc").Find(&visits).Error
	return visits, err
}

func (s *careService) GetStudentMedicalHistory(studentID uint) ([]models.MedicalVisit, error) {
	var visits []models.MedicalVisit
	err := s.db.Where("student_id = ?", studentID).Order("visit_date desc").Find(&visits).Error
	return visits, err
}

func (s *careService) RecordAnthropometry(record models.AnthropometryRecord) error {
	record.Date = time.Now()
	// Basic BMI Calculation
	if record.Height > 0 && record.Weight > 0 {
		hMeter := record.Height / 100
		record.BMI = record.Weight / (hMeter * hMeter)
		
		if record.BMI < 18.5 {
			record.BMICategory = "Underweight"
		} else if record.BMI < 25 {
			record.BMICategory = "Normal"
		} else if record.BMI < 30 {
			record.BMICategory = "Overweight"
		} else {
			record.BMICategory = "Obese"
		}
	}
	return s.db.Create(&record).Error
}

func (s *careService) GetAnthropometryRecords() ([]models.AnthropometryRecord, error) {
	var records []models.AnthropometryRecord
	err := s.db.Preload("Student").Order("date desc").Find(&records).Error
	return records, err
}

func (s *careService) GetMCUHistory(studentID uint) ([]models.AnthropometryRecord, error) {
	var records []models.AnthropometryRecord
	err := s.db.Where("student_id = ?", studentID).Order("date desc").Find(&records).Error
	return records, err
}

func (s *careService) GetMedicines() ([]models.Medicine, error) {
	var medicines []models.Medicine
	err := s.db.Find(&medicines).Error
	return medicines, err
}

func (s *careService) AddMedicine(m models.Medicine) error {
	return s.db.Create(&m).Error
}

func (s *careService) RecordMedicineMutation(m models.MedicineMutation) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&m).Error; err != nil {
			return err
		}
		
		change := m.Quantity
		if m.Type == "OUT" {
			change = -m.Quantity
		}
		
		return tx.Model(&models.Medicine{}).Where("id = ?", m.MedicineID).
			UpdateColumn("stock", gorm.Expr("stock + ?", change)).Error
	})
}

func (s *careService) GetBeds() ([]models.Bed, error) {
	var beds []models.Bed
	err := s.db.Preload("Facility").Find(&beds).Error
	return beds, err
}

func (s *careService) CreateCertificate(cert models.MedicalCertificate) error {
	// Generate Reference Number if empty
	if cert.RefNumber == "" {
		cert.RefNumber = "MED/" + time.Now().Format("20060102") + "/" + strings.ToUpper(utils.RandomString(4))
	}
	return s.db.Create(&cert).Error
}

func (s *careService) CreateObservation(obs models.ObservationRecord) error {
	obs.CheckIn = time.Now()
	obs.Status = "Active"
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&obs).Error; err != nil {
			return err
		}
		// Mark bed as occupied
		return tx.Model(&models.Bed{}).Where("id = ?", obs.BedID).Update("is_occupied", true).Error
	})
}

func (s *careService) UpdateObservation(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.ObservationRecord{}).Where("id = ?", id).Updates(updates).Error
}

func (s *careService) DischargeObservation(id uint) error {
	slog.Info("Discharging observation", "id", id)
	return s.db.Transaction(func(tx *gorm.DB) error {
		var obs models.ObservationRecord
		if err := tx.First(&obs, id).Error; err != nil {
			slog.Error("Error finding observation", "id", id, "error", err)
			return err
		}
		
		now := time.Now()
		// Update Observation Record
		if err := tx.Model(&obs).Updates(map[string]interface{}{
			"status":    "Discharged",
			"check_out": &now,
		}).Error; err != nil {
			slog.Error("Error updating observation status", "id", id, "error", err)
			return err
		}

		// Update Medical Visit status to 'Treated'
		if err := tx.Model(&models.MedicalVisit{}).Where("id = ?", obs.VisitID).Update("status", "Treated").Error; err != nil {
			slog.Error("Error updating visit status", "visit_id", obs.VisitID, "error", err)
			return err
		}

		// Free bed
		if err := tx.Model(&models.Bed{}).Where("id = ?", obs.BedID).Update("is_occupied", false).Error; err != nil {
			slog.Error("Error freeing bed", "bed_id", obs.BedID, "error", err)
			return err
		}

		slog.Info("Discharge successful", "id", id)
		return nil
	})
}

func (s *careService) GetActiveObservations() ([]models.ObservationRecord, error) {
	var obs []models.ObservationRecord
	err := s.db.Preload("Visit.Student").Preload("Bed.Facility").Where("status = ?", "Active").Find(&obs).Error
	return obs, err
}
