package models

import (
	"time"

	"gorm.io/gorm"
)

// Medicine represents UKS/Clinic medicine inventory
type Medicine struct {
	gorm.Model
	Code        string     `gorm:"uniqueIndex;not null" json:"code"`
	Name        string     `gorm:"not null" json:"name"`
	Stock       int        `gorm:"default:0" json:"stock"`
	Unit        string     `gorm:"default:'Pcs'" json:"unit"` // Tablet, Botol, Pcs
	Category    string     `json:"category"` // Obat Bebas, Keras, Alkes
	ExpiredAt   *time.Time `json:"expired_at"`
	Description *string    `gorm:"type:text" json:"description"`
}

// MedicineMutation logs inventory changes
type MedicineMutation struct {
	gorm.Model
	MedicineID uint      `json:"medicine_id"`
	Medicine   *Medicine `json:"medicine,omitempty"`
	Type       string    `json:"type"` // IN, OUT, ADJUST
	Quantity   int       `json:"quantity"`
	Notes      string    `json:"notes"`
	CreatedBy  string    `json:"created_by"`
}

// Disease represents a catalog of diseases/illnesses (ICD-10)
type Disease struct {
	gorm.Model
	Code        string  `gorm:"uniqueIndex" json:"code"` // ICD-10 Code
	Name        string  `gorm:"not null" json:"name"`
	Category    string  `json:"category"`
	Description *string `gorm:"type:text" json:"description"`
}

// MedicalVisit represents a polyclinic visit record
type MedicalVisit struct {
	gorm.Model
	StudentID     uint      `json:"student_id"`
	Student       *Student  `json:"student,omitempty"`
	VisitDate     time.Time `json:"visit_date"`
	Anamnesis     string    `gorm:"type:text" json:"anamnesis"`
	Diagnosis     string    `gorm:"type:text" json:"diagnosis"`
	ICD10Code     *string   `json:"icd10_code"`
	Treatment     string    `gorm:"type:text" json:"treatment"`
	MedicineGiven string    `json:"medicine_given"` // JSON or text summary
	Status        string    `json:"status"`         // Treated, Referred, Resting
	RecordedBy    uint      `json:"recorded_by"`
	Medic         *User     `gorm:"foreignKey:RecordedBy" json:"medic,omitempty"`
}

// ObservationRecord tracks students resting in UKS
type ObservationRecord struct {
	gorm.Model
	VisitID    uint          `json:"visit_id"`
	Visit      *MedicalVisit `json:"visit,omitempty"`
	BedID      uint          `json:"bed_id"`
	Bed        *Bed          `json:"bed,omitempty"`
	CheckIn    time.Time     `json:"check_in"`
	CheckOut   *time.Time    `json:"check_out"`
	Temperature float64      `json:"temperature"`
	BloodPressure string     `json:"blood_pressure"`
	Meds       string        `json:"meds"`
	Notes      string        `gorm:"type:text" json:"notes"`
	Status     string        `json:"status"` // Active, Discharged
}

// MedicalCertificate for sick/healthy letters
type MedicalCertificate struct {
	gorm.Model
	VisitID       uint          `json:"visit_id"`
	Visit         *MedicalVisit `json:"visit,omitempty"`
	Type          string        `json:"type"` // SAKIT, SEHAT
	StartDate     time.Time     `json:"start_date"`
	EndDate       time.Time     `json:"end_date"`
	RefNumber     string        `gorm:"uniqueIndex" json:"ref_number"`
	DigitalSign   string        `json:"digital_sign"`
}

// AnthropometryRecord for MCU (3 monthly)
type AnthropometryRecord struct {
	gorm.Model
	StudentID  uint      `json:"student_id"`
	Student    *Student  `json:"student,omitempty"`
	Date       time.Time `json:"date"`
	Height     float64   `json:"height"` // cm
	Weight     float64   `json:"weight"` // kg
	BMI        float64   `json:"bmi"`
	BMICategory string   `json:"bmi_category"`
	VisionLeft  string    `json:"vision_left"`
	VisionRight string    `json:"vision_right"`
	BloodType   string    `json:"blood_type"`
	Allergies   string    `json:"allergies"`
	RecordedBy  uint      `json:"recorded_by"`
}

// MedicalFacility represents ward/room
type MedicalFacility struct {
	gorm.Model
	Name     string `json:"name"`
	Location string `json:"location"`
	Type     string `json:"type"` // UKS Putra, UKS Putri, Klinik
}

// Bed management
type Bed struct {
	gorm.Model
	FacilityID uint             `json:"facility_id"`
	Facility   *MedicalFacility `json:"facility,omitempty"`
	Number     string           `json:"number"`
	IsOccupied bool             `gorm:"default:false" json:"is_occupied"`
}
