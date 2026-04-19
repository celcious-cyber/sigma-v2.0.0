package sigmabase

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/middleware"
)

// RegisterAdminRoutes sets up the Sigmabase endpoints for the admin gateway
func RegisterAdminRoutes(router fiber.Router, service StudentService) {
	h := NewSigmabaseHandler(service)
	
	base := router.Group("/base", middleware.Protected())
	
	// Statistics
	base.Get("/stats", h.GetStats)

	// Students
	base.Get("/students", h.GetStudents)
	base.Get("/students/:id", h.GetStudentDetails)
	base.Post("/students", h.CreateStudent)
	base.Post("/students/bulk", h.BulkCreateStudents)
	base.Put("/students/:id", h.UpdateStudent)
	base.Delete("/students/:id", h.DeleteStudent)

	// Teachers
	base.Get("/teachers", h.GetTeachers)
	base.Get("/teachers/:id", h.GetTeacherDetails)
	base.Post("/teachers", h.CreateTeacher)
	base.Post("/teachers/bulk", h.BulkCreateTeachers)
	base.Put("/teachers/:id", h.UpdateTeacher)
	base.Delete("/teachers/:id", h.DeleteTeacher)
	base.Post("/teachers/:id/photo", h.UploadTeacherPhoto)

	// Alumni
	base.Get("/alumni", h.GetAlumni)
	base.Post("/alumni", h.CreateAlumni)
	base.Put("/alumni/:id", h.UpdateAlumni)
	base.Delete("/alumni/:id", h.DeleteAlumni)
	base.Post("/alumni/:id/photo", h.UploadAlumniPhoto)
	base.Post("/students/:id/graduate", h.GraduateStudent)

	// Classrooms
	base.Get("/classrooms", h.GetClassrooms)
}
