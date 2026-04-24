package sigmaedu

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/middleware"
)

// RegisterAdminRoutes sets up routes for Backoffice Academic management
func RegisterAdminRoutes(router fiber.Router, service EduService) {
	h := NewEduHandler(service)
	
	edu := router.Group("/edu", middleware.Protected())
	edu.Get("/stats", h.GetStats)
	edu.Get("/subjects", h.GetSubjects)
	edu.Post("/subjects", h.CreateSubject)
	edu.Put("/subjects/:id", h.UpdateSubject)
	edu.Delete("/subjects/:id", h.DeleteSubject)
	// Study Hours
	edu.Get("/hours", h.GetStudyHours)
	edu.Post("/hours", h.CreateStudyHour)
	edu.Put("/hours/:id", h.UpdateStudyHour)
	edu.Delete("/hours/:id", h.DeleteStudyHour)
	// Academic Calendar
	edu.Get("/calendar", h.GetCalendarEvents)
	edu.Post("/calendar", h.CreateCalendarEvent)
	edu.Put("/calendar/:id", h.UpdateCalendarEvent)
	edu.Delete("/calendar/:id", h.DeleteCalendarEvent)
	// Classrooms
	edu.Get("/classrooms", h.GetClassrooms)
	edu.Post("/classrooms", h.CreateClassroom)
	edu.Put("/classrooms/:id", h.UpdateClassroom)
	edu.Delete("/classrooms/:id", h.DeleteClassroom)
	edu.Get("/schedules", h.GetSchedules)
	edu.Post("/schedules", h.CreateSchedule)
	edu.Put("/schedules/:id", h.UpdateSchedule)
	edu.Delete("/schedules/:id", h.DeleteSchedule)
	
	// Attendance
	edu.Get("/attendance", h.GetAttendances)
	edu.Get("/attendance/report", h.GetAttendanceReport)
	edu.Post("/attendance/bulk", h.BulkRecordAttendance)
	edu.Delete("/attendance/:id", h.DeleteAttendance)
	edu.Get("/classrooms/:id/students", h.GetStudentsByClass)
	edu.Get("/students/unassigned", h.GetUnassignedStudents)
	edu.Post("/students/assign", h.BulkUpdateStudentClass)
	
	// Assessments
	edu.Get("/assessments", h.GetAssessments)
	edu.Post("/assessments/bulk", h.BulkRecordAssessments)

	// Tahfidz
	edu.Get("/tahfidz", h.GetTahfidzRecords)
	edu.Post("/tahfidz/bulk", h.BulkRecordTahfidz)

	// Lesson Memorization
	edu.Get("/memorization", h.GetLessonMemorizations)
	edu.Post("/memorization/bulk", h.BulkRecordLessonMemorizations)

	// Teacher Attendance
	edu.Get("/teacher-attendance", h.GetTeacherAttendance)
	edu.Post("/teacher-attendance/bulk", h.BulkRecordTeacherAttendance)

	// Teaching Journal
	edu.Get("/journals", h.GetTeachingJournals)
	edu.Post("/journals", h.CreateTeachingJournal)
	edu.Put("/journals/:id", h.UpdateTeachingJournal)
	edu.Delete("/journals/:id", h.DeleteTeachingJournal)
}

// RegisterPortalRoutes sets up routes for the Parent Portal (Sigmagate)
func RegisterPortalRoutes(router fiber.Router, service EduService) {
	// h := NewEduHandler(service)
	
	// portal := router.Group("/edu", middleware.Protected())
	// portal.Get("/assessments", h.GetAssessments)
	// portal.Get("/attendances", h.GetAttendances)
}
