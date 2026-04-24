package sigmabase

import (
	"time"

	"gorm.io/gorm"
	"sigma-api/internal/core/models"
)

type StudentService interface {
	GetStudentByID(id uint) (*models.Student, error)
	RefreshViolationPoints(studentID uint) error
	IsPermitBlocked(studentID uint) (bool, map[string]interface{}, error)
	// Students
	CreateStudent(student *models.Student) error
	BulkCreateStudents(students []models.Student) error
	UpdateStudent(id uint, data *models.Student) error
	DeleteStudent(id uint) error
	// Teachers
	GetAllTeachers() ([]models.Teacher, error)
	GetTeacherByID(id uint) (*models.Teacher, error)
	CreateTeacher(teacher *models.Teacher) error
	BulkCreateTeachers(teachers []models.Teacher) error
	UpdateTeacher(id uint, data map[string]interface{}) error
	DeleteTeacher(id uint) error
	// Alumni
	GetAllAlumni() ([]models.Alumni, error)
	GetAlumniByID(id uint) (*models.Alumni, error)
	CreateAlumni(alumni *models.Alumni) error
	UpdateAlumni(id uint, data map[string]interface{}) error
	DeleteAlumni(id uint) error
	GraduateStudent(studentID uint, alumniData map[string]interface{}) (uint, error)
	// Units
	GetAllUnits() ([]models.InstitutionalUnit, error)
	GetUnitByID(id uint) (*models.InstitutionalUnit, error)
	CreateUnit(unit *models.InstitutionalUnit) error
	UpdateUnit(id uint, data map[string]interface{}) error
	DeleteUnit(id uint) error
	SeedUnits() error
	// Classrooms
	GetAllClassrooms() ([]models.Classroom, error)
	// Dashboard Stats
	GetStats() (map[string]interface{}, error)
}

type studentService struct {
	db *gorm.DB
}

func NewStudentService(db *gorm.DB) StudentService {
	return &studentService{db}
}

func (s *studentService) GetStudentByID(id uint) (*models.Student, error) {
	var student models.Student
	if err := s.db.Preload("Classroom").First(&student, id).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (s *studentService) GetTeacherByID(id uint) (*models.Teacher, error) {
	var teacher models.Teacher
	if err := s.db.First(&teacher, id).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (s *studentService) CreateTeacher(teacher *models.Teacher) error {
	return s.db.Create(teacher).Error
}

func (s *studentService) BulkCreateTeachers(teachers []models.Teacher) error {
	return s.db.CreateInBatches(teachers, 100).Error
}

func (s *studentService) UpdateTeacher(id uint, data map[string]interface{}) error {
	// Map-based updates process only the provided keys, preserving other fields.
	// This is the most robust way to handle partial updates and boolean false values.
	return s.db.Model(&models.Teacher{}).Where("id = ?", id).Updates(data).Error
}

func (s *studentService) DeleteTeacher(id uint) error {
	return s.db.Delete(&models.Teacher{}, id).Error
}

func (s *studentService) GetAllTeachers() ([]models.Teacher, error) {
	var teachers []models.Teacher
	err := s.db.Find(&teachers).Error
	return teachers, err
}

func (s *studentService) GetAllAlumni() ([]models.Alumni, error) {
	var alumni []models.Alumni
	err := s.db.Find(&alumni).Error
	return alumni, err
}

func (s *studentService) GetAlumniByID(id uint) (*models.Alumni, error) {
	var alumni models.Alumni
	if err := s.db.First(&alumni, id).Error; err != nil {
		return nil, err
	}
	return &alumni, nil
}

func (s *studentService) CreateAlumni(alumni *models.Alumni) error {
	return s.db.Create(alumni).Error
}

func (s *studentService) UpdateAlumni(id uint, data map[string]interface{}) error {
	return s.db.Model(&models.Alumni{}).Where("id = ?", id).Updates(data).Error
}

func (s *studentService) DeleteAlumni(id uint) error {
	return s.db.Delete(&models.Alumni{}, id).Error
}

func (s *studentService) GraduateStudent(studentID uint, alumniData map[string]interface{}) (uint, error) {
	var newAlumniID uint
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch Student Origin
		var student models.Student
		if err := tx.First(&student, studentID).Error; err != nil {
			return err
		}

		// 2. Create Alumni Record
		alumni := models.Alumni{
			Name:           student.Name,
			Gender:         student.Gender,
			NIK:            student.NIK,
			BirthPlace:     student.BirthPlace,
			BirthDate:      student.BirthDate,
			NIS:            &student.NIS,
			Address:        student.Address,
		}
		
		if val, ok := alumniData["graduation_year"]; ok { alumni.GraduationYear = val.(string) }
		if val, ok := alumniData["batch"]; ok { alumni.Batch = val.(string) }
		if val, ok := alumniData["whatsapp"]; ok { v := val.(string); alumni.WhatsApp = &v }
		if val, ok := alumniData["email"]; ok { v := val.(string); alumni.Email = &v }
		if val, ok := alumniData["nik"]; ok { 
			v := val.(string)
			if v != "" { alumni.NIK = &v } else { alumni.NIK = nil }
		}
		if val, ok := alumniData["photo"]; ok { v := val.(string); alumni.Photo = &v }
		if val, ok := alumniData["service_status"]; ok { alumni.ServiceStatus = val.(string) }

		if err := tx.Create(&alumni).Error; err != nil {
			return err
		}
		newAlumniID = alumni.ID

		// 3. Mark Student as Inactive (Archived)
		if err := tx.Model(&student).Update("is_active", false).Error; err != nil {
			return err
		}

		return nil
	})
	return newAlumniID, err
}

func (s *studentService) GetAllClassrooms() ([]models.Classroom, error) {
	var classrooms []models.Classroom
	err := s.db.Find(&classrooms).Error
	return classrooms, err
}

// RefreshViolationPoints recalculates total points for a student
func (s *studentService) RefreshViolationPoints(studentID uint) error {
	var total int64
	s.db.Model(&models.Violation{}).Where("student_id = ?", studentID).Select("sum(points)").Row().Scan(&total)
	
	// Assuming we add a field TotalViolationPoints to Student model or just return it
	// In Laravel it was a field, let's update Student model if needed or just handle it here.
	return s.db.Model(&models.Student{}).Where("id = ?", studentID).Update("total_violation_points", total).Error
}

// IsPermitBlocked implements the quota logic (2x per semester Pulang, 3x per month Keluar)
func (s *studentService) IsPermitBlocked(studentID uint) (bool, map[string]interface{}, error) {
	now := time.Now()
	
	// Month check (Keluar - 3x limit)
	var monthlyCount int64
	s.db.Model(&models.Permit{}).
		Where("student_id = ? AND type = 'Keluar' AND is_emergency = ?", studentID, false).
		Where("strftime('%Y-%m', out_time) = ?", now.Format("2006-01")).
		Count(&monthlyCount)

	// Semester check (Pulang - 2x limit)
	// Semester 1: July-Dec (7-12), Semester 2: Jan-June (1-6)
	var semesterStart time.Time
	if now.Month() >= 7 {
		semesterStart = time.Date(now.Year(), 7, 1, 0, 0, 0, 0, time.Local)
	} else {
		semesterStart = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	}

	var semesterCount int64
	s.db.Model(&models.Permit{}).
		Where("student_id = ? AND type = 'Pulang' AND is_emergency = ?", studentID, false).
		Where("out_time >= ?", semesterStart).
		Count(&semesterCount)

	isBlocked := semesterCount >= 2 || monthlyCount >= 3
	
	stats := map[string]interface{}{
		"monthly_keluar":  monthlyCount,
		"semester_pulang": semesterCount,
		"limits": map[string]int{
			"monthly_keluar":  3,
			"semester_pulang": 2,
		},
	}

	return isBlocked, stats, nil
}
func (s *studentService) GetStats() (map[string]interface{}, error) {
	var totalSantri int64
	var totalGuru int64
	var totalAlumni int64

	s.db.Model(&models.Student{}).Count(&totalSantri)
	s.db.Model(&models.Teacher{}).Count(&totalGuru)
	s.db.Model(&models.Alumni{}).Count(&totalAlumni)
	
	// Gender distribution (Students)
	var santriLaki, santriPerempuan int64
	s.db.Model(&models.Student{}).Where("gender = ?", "L").Count(&santriLaki)
	s.db.Model(&models.Student{}).Where("gender = ?", "P").Count(&santriPerempuan)

	// Teacher Status Distribution
	type StatusCount struct {
		Status string
		Count  int64
	}
	var teacherStatus []StatusCount
	s.db.Model(&models.Teacher{}).Select("status, count(*) as count").Group("status").Scan(&teacherStatus)

	// Alumni Service Status
	var alumniMengabdi, alumniTidak int64
	s.db.Model(&models.Alumni{}).Where("service_status = ?", "Mengabdi").Count(&alumniMengabdi)
	s.db.Model(&models.Alumni{}).Where("service_status != ?", "Mengabdi").Count(&alumniTidak)

	// Trend Pendaftaran Students (Last 6 Months)
	type MonthlyCount struct {
		Month string
		Count int
	}
	var trend []MonthlyCount
	s.db.Raw(`SELECT strftime('%Y-%m', created_at) as month, count(*) as count 
			  FROM students 
			  GROUP BY month 
			  ORDER BY month DESC 
			  LIMIT 6`).Scan(&trend)

	trendLabels := []string{}
	trendData := []int{}
	for i := len(trend) - 1; i >= 0; i-- {
		trendLabels = append(trendLabels, trend[i].Month)
		trendData = append(trendData, trend[i].Count)
	}

	return map[string]interface{}{
		"total_santri": totalSantri,
		"total_guru":   totalGuru,
		"total_alumni": totalAlumni,
		"chart_gender": map[string]int64{
			"L": santriLaki,
			"P": santriPerempuan,
		},
		"chart_teacher_status": teacherStatus,
		"chart_alumni_service": map[string]int64{
			"Mengabdi": alumniMengabdi,
			"Lainnya":  alumniTidak,
		},
		"alumni_regions": s.getTopRegions(),
		"chart_trend": map[string]interface{}{
			"labels": trendLabels,
			"data":   trendData,
		},
	}, nil
}

func (s *studentService) getTopRegions() []map[string]interface{} {
	type RegionCount struct {
		BirthPlace string `gorm:"column:birth_place"`
		Count      int64  `gorm:"column:count"`
	}
	var results []RegionCount
	s.db.Model(&models.Alumni{}).
		Select("birth_place, count(*) as count").
		Where("birth_place != ''").
		Group("birth_place").
		Order("count DESC").
		Limit(5).
		Scan(&results)

	formatted := []map[string]interface{}{}
	for _, r := range results {
		formatted = append(formatted, map[string]interface{}{
			"name":  r.BirthPlace,
			"value": r.Count,
		})
	}
	return formatted
}

func (s *studentService) CreateStudent(student *models.Student) error {
	return s.db.Create(student).Error
}

func (s *studentService) BulkCreateStudents(students []models.Student) error {
	return s.db.CreateInBatches(students, 100).Error
}

func (s *studentService) UpdateStudent(id uint, data *models.Student) error {
	return s.db.Model(&models.Student{}).Where("id = ?", id).Updates(data).Error
}

func (s *studentService) DeleteStudent(id uint) error {
	return s.db.Delete(&models.Student{}, id).Error
}

func (s *studentService) GetAllUnits() ([]models.InstitutionalUnit, error) {
	var units []models.InstitutionalUnit
	if err := s.db.Find(&units).Error; err != nil {
		return nil, err
	}
	return units, nil
}

func (s *studentService) GetUnitByID(id uint) (*models.InstitutionalUnit, error) {
	var unit models.InstitutionalUnit
	if err := s.db.First(&unit, id).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

func (s *studentService) CreateUnit(unit *models.InstitutionalUnit) error {
	return s.db.Create(unit).Error
}

func (s *studentService) UpdateUnit(id uint, data map[string]interface{}) error {
	return s.db.Model(&models.InstitutionalUnit{}).Where("id = ?", id).Updates(data).Error
}

func (s *studentService) DeleteUnit(id uint) error {
	return s.db.Delete(&models.InstitutionalUnit{}, id).Error
}

func (s *studentService) SeedUnits() error {
	units := []models.InstitutionalUnit{
		{Name: "SMP", Code: "SMP-AH", Description: "Sekolah Menengah Pertama Al-Hikmah"},
		{Name: "MA", Code: "MA-AH", Description: "Madrasah Aliyah Al-Hikmah"},
		{Name: "KMI", Code: "KMI-AH", Description: "Kulliyatul Mu'allimin Al-Islamiyah (Pondok)"},
	}

	for _, unit := range units {
		var existing models.InstitutionalUnit
		if err := s.db.Where("name = ?", unit.Name).First(&existing).Error; err != nil {
			s.db.Create(&unit)
		}
	}
	return nil
}

