package sigmaedu

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"sigma-api/internal/core/models"
)

type EduService interface {
	// Dashboard Stats
	GetStats() (map[string]interface{}, error)
	// Subjects
	GetAllSubjects() ([]models.Subject, error)
	CreateSubject(s *models.Subject) error
	UpdateSubject(id uint, s *models.Subject) error
	DeleteSubject(id uint) error
	// Study Hours
	GetAllStudyHours() ([]models.StudyHour, error)
	CreateStudyHour(h *models.StudyHour) error
	UpdateStudyHour(id uint, h *models.StudyHour) error
	DeleteStudyHour(id uint) error
	GetSchedules() ([]models.TeacherSchedule, error)
	GetSchedulesByClass(classID uint) ([]models.TeacherSchedule, error)
	CreateSchedule(sc *models.TeacherSchedule) error
	UpdateSchedule(id uint, sc *models.TeacherSchedule) error
	DeleteSchedule(id uint) error
	// Academic Calendar
	GetAllCalendarEvents() ([]models.AcademicCalendar, error)
	CreateCalendarEvent(e *models.AcademicCalendar) error
	UpdateCalendarEvent(id uint, e *models.AcademicCalendar) error
	DeleteCalendarEvent(id uint) error
	// Attendance
	GetAttendances(classID uint, date string, subjectID *uint) ([]models.Attendance, error)
	GetAttendanceReport(classID uint, subjectID *uint, startDate, endDate string) ([]models.Attendance, error)
	GetAttendancesByStudent(studentID uint) ([]models.Attendance, error)
	RecordAttendance(a models.Attendance) error
	BulkRecordAttendance(attendances []models.Attendance) error
	DeleteAttendance(id uint) error
	// Assessments
	GetAssessments(classID uint, subjectID uint, term, year, assessmentType string) ([]models.Assessment, error)
	BulkRecordAssessments(assessments []models.Assessment) error
	DeleteAssessments(classID uint, subjectID uint, term, year, assessmentType string) error
	GetAssessmentsByStudent(studentID uint) ([]models.Assessment, error)
	// Tahfidz
	GetTahfidzRecords(classID uint, date string) ([]models.QuranMemorization, error)
	RecordTahfidz(record models.QuranMemorization) error
	BulkRecordTahfidz(records []models.QuranMemorization) error
	DeleteTahfidz(classID uint, date string, tahfidzType string) error
	DeleteTahfidzRecord(id uint) error
	GetTahfidzHistoryByStudent(studentID uint) ([]models.QuranMemorization, error)
	// Lesson Memorization
	GetLessonMemorizations(classID uint, date string) ([]models.LessonMemorization, error)
	BulkRecordLessonMemorizations(records []models.LessonMemorization) error
	GetLessonMemorizationHistory(studentID uint) ([]models.LessonMemorization, error)
	DeleteLessonMemorizationRecord(id uint) error
	// Teacher Attendance
	GetTeacherAttendance(date string) ([]models.TeacherAttendance, error)
	BulkRecordTeacherAttendance(records []models.TeacherAttendance) error
	// Teaching Journal
	GetTeachingJournals(date string, classroomID *uint, subjectID *uint) ([]models.TeachingJournal, error)
	CreateTeachingJournal(journal *models.TeachingJournal) error
	GetLastJournal(teacherID uint, subjectID uint) (*models.TeachingJournal, error)
	UpdateTeachingJournal(id uint, journal *models.TeachingJournal) error
	DeleteTeachingJournal(id uint) error
	// Classrooms
	GetAllClassrooms() ([]models.Classroom, error)
	CreateClassroom(c *models.Classroom) error
	UpdateClassroom(id uint, c *models.Classroom) error
	DeleteClassroom(id uint) error
	// Helpers
	GetAllTeachers() ([]models.Teacher, error)
	GetStudentsByClass(classID uint) ([]models.Student, error)
	GetUnassignedStudents() ([]models.Student, error)
	BulkUpdateStudentClass(studentIDs []uint, classID *uint) error
}

type eduService struct {
	db *gorm.DB
}

func NewEduService(db *gorm.DB) EduService {
	return &eduService{db}
}

func (s *eduService) GetStats() (map[string]interface{}, error) {
	var totalSubjects int64
	var totalSchedules int64
	var totalAttendanceToday int64

	s.db.Model(&models.Subject{}).Count(&totalSubjects)
	s.db.Model(&models.TeacherSchedule{}).Count(&totalSchedules)
	
	today := time.Now().Truncate(24 * time.Hour)
	s.db.Model(&models.Attendance{}).Where("date >= ?", today).Count(&totalAttendanceToday)

	return map[string]interface{}{
		"total_subjects":          totalSubjects,
		"total_schedules":         totalSchedules,
		"total_attendance_today": totalAttendanceToday,
	}, nil
}

func (s *eduService) GetAllSubjects() ([]models.Subject, error) {
	var subjects []models.Subject
	err := s.db.Find(&subjects).Error
	return subjects, err
}

func (s *eduService) CreateSubject(sub *models.Subject) error {
	return s.db.Create(sub).Error
}

func (s *eduService) UpdateSubject(id uint, sub *models.Subject) error {
	return s.db.Model(&models.Subject{}).Where("id = ?", id).Updates(sub).Error
}

func (s *eduService) DeleteSubject(id uint) error {
	return s.db.Delete(&models.Subject{}, id).Error
}

func (s *eduService) GetAllStudyHours() ([]models.StudyHour, error) {
	var hours []models.StudyHour
	err := s.db.Find(&hours).Error
	return hours, err
}

func (s *eduService) CreateStudyHour(h *models.StudyHour) error {
	return s.db.Create(h).Error
}

func (s *eduService) UpdateStudyHour(id uint, h *models.StudyHour) error {
	return s.db.Model(&models.StudyHour{}).Where("id = ?", id).Updates(h).Error
}

func (s *eduService) DeleteStudyHour(id uint) error {
	return s.db.Delete(&models.StudyHour{}, id).Error
}

func (s *eduService) GetAllCalendarEvents() ([]models.AcademicCalendar, error) {
	var events []models.AcademicCalendar
	err := s.db.Order("start_date asc").Find(&events).Error
	return events, err
}

func (s *eduService) CreateCalendarEvent(e *models.AcademicCalendar) error {
	return s.db.Create(e).Error
}

func (s *eduService) UpdateCalendarEvent(id uint, e *models.AcademicCalendar) error {
	return s.db.Model(&models.AcademicCalendar{}).Where("id = ?", id).Updates(e).Error
}

func (s *eduService) DeleteCalendarEvent(id uint) error {
	return s.db.Delete(&models.AcademicCalendar{}, id).Error
}

func (s *eduService) GetSchedulesByClass(classID uint) ([]models.TeacherSchedule, error) {
	var schedules []models.TeacherSchedule
	err := s.db.Preload("Teacher").Preload("Subject").Where("classroom_id = ?", classID).Find(&schedules).Error
	return schedules, err
}

func (s *eduService) CreateSchedule(sc *models.TeacherSchedule) error {
	return s.db.Create(sc).Error
}

func (s *eduService) UpdateSchedule(id uint, sc *models.TeacherSchedule) error {
	return s.db.Model(&models.TeacherSchedule{}).Where("id = ?", id).Updates(sc).Error
}

func (s *eduService) DeleteSchedule(id uint) error {
	return s.db.Delete(&models.TeacherSchedule{}, id).Error
}

func (s *eduService) GetAllClassrooms() ([]models.Classroom, error) {
	var classrooms []models.Classroom
	err := s.db.Preload("WaliKelas").Find(&classrooms).Error
	return classrooms, err
}

func (s *eduService) CreateClassroom(c *models.Classroom) error {
	return s.db.Create(c).Error
}

func (s *eduService) UpdateClassroom(id uint, c *models.Classroom) error {
	return s.db.Model(&models.Classroom{}).Where("id = ?", id).Updates(c).Error
}

func (s *eduService) DeleteClassroom(id uint) error {
	return s.db.Delete(&models.Classroom{}, id).Error
}

func (s *eduService) GetAllTeachers() ([]models.Teacher, error) {
	var teachers []models.Teacher
	err := s.db.Find(&teachers).Error
	return teachers, err
}

func (s *eduService) GetStudentsByClass(classID uint) ([]models.Student, error) {
	var students []models.Student
	err := s.db.Preload("Classroom").Where("classroom_id = ?", classID).Find(&students).Error
	return students, err
}

func (s *eduService) GetAssessments(classID uint, subjectID uint, term, year, assessmentType string) ([]models.Assessment, error) {
	var assessments []models.Assessment
	query := s.db.Preload("Student").Preload("Subject").
		Where("classroom_id = ? AND term = ? AND academic_year = ?", classID, term, year)
	
	if subjectID > 0 {
		query = query.Where("subject_id = ?", subjectID)
	} else {
		query = query.Where("subject_id IS NULL OR subject_id = 0")
	}
	
	if assessmentType != "" {
		query = query.Where("type = ?", assessmentType)
	}

	err := query.Find(&assessments).Error
	return assessments, err
}

func (s *eduService) BulkRecordAssessments(assessments []models.Assessment) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, a := range assessments {
			var existing models.Assessment
			query := tx.Where("student_id = ? AND classroom_id = ? AND term = ? AND academic_year = ? AND type = ?", 
				a.StudentID, a.ClassroomID, a.Term, a.AcademicYear, a.Type)
			
			if a.SubjectID != nil && *a.SubjectID > 0 {
				query = query.Where("subject_id = ?", *a.SubjectID)
			} else {
				query = query.Where("subject_id IS NULL OR subject_id = 0")
			}
			
			if err := query.First(&existing).Error; err == nil {
				// Update
				if err := tx.Model(&existing).Updates(map[string]interface{}{
					"score": a.Score,
					"grade": a.Grade,
					"remarks": a.Remarks,
				}).Error; err != nil {
					return err
				}
			} else {
				// Create
				if err := tx.Create(&a).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *eduService) DeleteAssessments(classID uint, subjectID uint, term, year, assessmentType string) error {
	query := s.db.Where("classroom_id = ? AND term = ? AND academic_year = ? AND type = ?", 
		classID, term, year, assessmentType)
	
	if subjectID > 0 {
		query = query.Where("subject_id = ?", subjectID)
	} else {
		query = query.Where("subject_id IS NULL OR subject_id = 0")
	}
	
	return query.Delete(&models.Assessment{}).Error
}

func (s *eduService) GetAssessmentsByStudent(studentID uint) ([]models.Assessment, error) {
	var assessments []models.Assessment
	err := s.db.Preload("Subject").Where("student_id = ?", studentID).Order("date desc").Find(&assessments).Error
	return assessments, err
}

func (s *eduService) DeleteAttendance(id uint) error {
	return s.db.Delete(&models.Attendance{}, id).Error
}

func (s *eduService) GetTahfidzRecords(classID uint, date string) ([]models.QuranMemorization, error) {
	var records []models.QuranMemorization
	parsedDate, _ := time.Parse("2006-01-02", date)
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 0, parsedDate.Location())

	err := s.db.Preload("Student").
		Joins("JOIN students ON students.id = quran_memorizations.student_id").
		Where("students.classroom_id = ? AND quran_memorizations.date BETWEEN ? AND ?", classID, startOfDay, endOfDay).
		Find(&records).Error
	return records, err
}

func (s *eduService) RecordTahfidz(record models.QuranMemorization) error {
	return s.db.Create(&record).Error
}

func (s *eduService) BulkRecordTahfidz(records []models.QuranMemorization) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, r := range records {
			var existing models.QuranMemorization
			// Match by student, date, and type
			startOfDay := time.Date(r.Date.Year(), r.Date.Month(), r.Date.Day(), 0, 0, 0, 0, r.Date.Location())
			endOfDay := time.Date(r.Date.Year(), r.Date.Month(), r.Date.Day(), 23, 59, 59, 0, r.Date.Location())
			
			query := tx.Where("student_id = ? AND type = ? AND date BETWEEN ? AND ?", r.StudentID, r.Type, startOfDay, endOfDay)
			
			if err := query.First(&existing).Error; err == nil {
				// Update
				if err := tx.Model(&existing).Updates(map[string]interface{}{
					"surah_name":  r.SurahName,
					"verse_start": r.VerseStart,
					"verse_end":   r.VerseEnd,
					"grade":       r.Grade,
					"remarks":     r.Remarks,
				}).Error; err != nil {
					return err
				}
			} else {
				// Create
				if err := tx.Create(&r).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *eduService) DeleteTahfidz(classID uint, date string, tahfidzType string) error {
	parsedDate, _ := time.Parse("2006-01-02", date)
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 0, parsedDate.Location())

	// Join with students to filter by classroom_id
	return s.db.Where("type = ? AND date BETWEEN ? AND ? AND student_id IN (SELECT id FROM students WHERE classroom_id = ?)", 
		tahfidzType, startOfDay, endOfDay, classID).Delete(&models.QuranMemorization{}).Error
}

func (s *eduService) GetTahfidzHistoryByStudent(studentID uint) ([]models.QuranMemorization, error) {
	var records []models.QuranMemorization
	err := s.db.Where("student_id = ?", studentID).Order("date desc").Find(&records).Error
	return records, err
}

func (s *eduService) DeleteTahfidzRecord(id uint) error {
	return s.db.Delete(&models.QuranMemorization{}, id).Error
}

func (s *eduService) GetLessonMemorizations(classID uint, date string) ([]models.LessonMemorization, error) {
	var records []models.LessonMemorization
	query := s.db.Joins("JOIN students ON students.id = lesson_memorizations.student_id")
	
	if classID > 0 {
		query = query.Where("students.classroom_id = ?", classID)
	}
	
	if date != "" {
		start := date + " 00:00:00"
		end := date + " 23:59:59"
		query = query.Where("lesson_memorizations.date BETWEEN ? AND ?", start, end)
	}

	err := query.Find(&records).Error
	return records, err
}

func (s *eduService) BulkRecordLessonMemorizations(records []models.LessonMemorization) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, r := range records {
			var existing models.LessonMemorization
			// Use start and end of day to check for existing record on that date
			startOfDay := time.Date(r.Date.Year(), r.Date.Month(), r.Date.Day(), 0, 0, 0, 0, r.Date.Location())
			endOfDay := time.Date(r.Date.Year(), r.Date.Month(), r.Date.Day(), 23, 59, 59, 0, r.Date.Location())
			
			err := tx.Where("student_id = ? AND subject_name = ? AND date BETWEEN ? AND ?", r.StudentID, r.SubjectName, startOfDay, endOfDay).First(&existing).Error
			if err == nil {
				r.ID = existing.ID
				if err := tx.Save(&r).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Create(&r).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *eduService) GetLessonMemorizationHistory(studentID uint) ([]models.LessonMemorization, error) {
	var records []models.LessonMemorization
	err := s.db.Where("student_id = ?", studentID).Order("date desc").Find(&records).Error
	return records, err
}

func (s *eduService) DeleteLessonMemorizationRecord(id uint) error {
	return s.db.Delete(&models.LessonMemorization{}, id).Error
}

func (s *eduService) GetTeacherAttendance(date string) ([]models.TeacherAttendance, error) {
	var records []models.TeacherAttendance
	parsedDate, _ := time.Parse("2006-01-02", date)
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 0, parsedDate.Location())

	err := s.db.Preload("Teacher").
		Where("date BETWEEN ? AND ?", startOfDay, endOfDay).
		Find(&records).Error
	return records, err
}

func (s *eduService) BulkRecordTeacherAttendance(records []models.TeacherAttendance) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, r := range records {
			var existing models.TeacherAttendance
			startOfDay := time.Date(r.Date.Year(), r.Date.Month(), r.Date.Day(), 0, 0, 0, 0, r.Date.Location())
			endOfDay := time.Date(r.Date.Year(), r.Date.Month(), r.Date.Day(), 23, 59, 59, 0, r.Date.Location())
			
			query := tx.Where("teacher_id = ? AND date BETWEEN ? AND ?", r.TeacherID, startOfDay, endOfDay)
			
			if err := query.First(&existing).Error; err == nil {
				if err := tx.Model(&existing).Updates(map[string]interface{}{
					"status":    r.Status,
					"check_in":  r.CheckIn,
					"check_out": r.CheckOut,
					"remarks":   r.Remarks,
				}).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Create(&r).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *eduService) GetTeachingJournals(date string, classroomID *uint, subjectID *uint) ([]models.TeachingJournal, error) {
	var records []models.TeachingJournal
	parsedDate, _ := time.Parse("2006-01-02", date)
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 0, parsedDate.Location())

	query := s.db.Preload("Teacher").Preload("Classroom").Preload("Subject").Preload("StudyHour").
		Where("date BETWEEN ? AND ?", startOfDay, endOfDay)
	
	if classroomID != nil && *classroomID > 0 {
		query = query.Where("classroom_id = ?", *classroomID)
	}
	if subjectID != nil && *subjectID > 0 {
		query = query.Where("subject_id = ?", *subjectID)
	}

	err := query.Find(&records).Error
	return records, err
}

func (s *eduService) CreateTeachingJournal(journal *models.TeachingJournal) error {
	// Validation: No overlap for same teacher, date, and study hour
	var existing models.TeachingJournal
	startOfDay := time.Date(journal.Date.Year(), journal.Date.Month(), journal.Date.Day(), 0, 0, 0, 0, journal.Date.Location())
	endOfDay := time.Date(journal.Date.Year(), journal.Date.Month(), journal.Date.Day(), 23, 59, 59, 0, journal.Date.Location())
	
	if err := s.db.Where("teacher_id = ? AND date BETWEEN ? AND ? AND study_hour_id = ?", 
		journal.TeacherID, startOfDay, endOfDay, journal.StudyHourID).First(&existing).Error; err == nil {
		return fmt.Errorf("jurnal untuk jam pelajaran ini sudah ada")
	}

	return s.db.Create(journal).Error
}

func (s *eduService) GetLastJournal(teacherID uint, subjectID uint) (*models.TeachingJournal, error) {
	var lastJournal models.TeachingJournal
	err := s.db.Where("teacher_id = ? AND subject_id = ?", teacherID, subjectID).
		Order("date desc").First(&lastJournal).Error
	if err != nil {
		return nil, err
	}
	return &lastJournal, nil
}

func (s *eduService) UpdateTeachingJournal(id uint, journal *models.TeachingJournal) error {
	return s.db.Model(&models.TeachingJournal{}).Where("id = ?", id).Updates(journal).Error
}

func (s *eduService) DeleteTeachingJournal(id uint) error {
	return s.db.Delete(&models.TeachingJournal{}, id).Error
}

func (s *eduService) GetAttendancesByStudent(studentID uint) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := s.db.Preload("Subject").Where("student_id = ?", studentID).Order("date desc").Find(&attendances).Error
	return attendances, err
}

func (s *eduService) GetAttendances(classID uint, date string, subjectID *uint) ([]models.Attendance, error) {
	var attendances []models.Attendance
	parsedDate, _ := time.Parse("2006-01-02", date)
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 0, parsedDate.Location())

	query := s.db.Preload("Student").Preload("Subject").Where("classroom_id = ? AND date BETWEEN ? AND ?", classID, startOfDay, endOfDay)
	if subjectID != nil {
		query = query.Where("subject_id = ?", *subjectID)
	}
	err := query.Find(&attendances).Error
	return attendances, err
}

func (s *eduService) GetAttendanceReport(classID uint, subjectID *uint, startDate, endDate string) ([]models.Attendance, error) {
	var attendances []models.Attendance
	
	query := s.db.Preload("Student").Preload("Subject").Where("classroom_id = ?", classID)
	
	if subjectID != nil {
		query = query.Where("subject_id = ?", *subjectID)
	} else {
		// If subjectID is nil, we might want to show all OR just "Harian" (null subject)
		// Usually report should show all unless filtered.
	}
	
	if startDate != "" {
		parsedStart, _ := time.Parse("2006-01-02", startDate)
		query = query.Where("date >= ?", parsedStart)
	}
	if endDate != "" {
		parsedEnd, _ := time.Parse("2006-01-02", endDate)
		// Set to end of day
		endOfDay := time.Date(parsedEnd.Year(), parsedEnd.Month(), parsedEnd.Day(), 23, 59, 59, 0, parsedEnd.Location())
		query = query.Where("date <= ?", endOfDay)
	}
	
	err := query.Order("date desc").Find(&attendances).Error
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

func (s *eduService) BulkRecordAttendance(attendances []models.Attendance) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, a := range attendances {
			// Upsert based on StudentID, ClassroomID, Date, and SubjectID
			var existing models.Attendance
			startOfDay := time.Date(a.Date.Year(), a.Date.Month(), a.Date.Day(), 0, 0, 0, 0, a.Date.Location())
			endOfDay := time.Date(a.Date.Year(), a.Date.Month(), a.Date.Day(), 23, 59, 59, 0, a.Date.Location())
			
			query := tx.Where("student_id = ? AND classroom_id = ? AND date BETWEEN ? AND ?", a.StudentID, a.ClassroomID, startOfDay, endOfDay)
			if a.SubjectID != nil {
				query = query.Where("subject_id = ?", *a.SubjectID)
			} else {
				query = query.Where("subject_id IS NULL")
			}
			
			if err := query.First(&existing).Error; err == nil {
				// Update
				if err := tx.Model(&existing).Update("status", a.Status).Error; err != nil {
					return err
				}
			} else {
				// Create
				if err := tx.Create(&a).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *eduService) GetUnassignedStudents() ([]models.Student, error) {
	var students []models.Student
	err := s.db.Where("classroom_id IS NULL OR classroom_id = 0").Find(&students).Error
	return students, err
}

func (s *eduService) BulkUpdateStudentClass(studentIDs []uint, classID *uint) error {
	return s.db.Model(&models.Student{}).Where("id IN ?", studentIDs).Update("classroom_id", classID).Error
}

