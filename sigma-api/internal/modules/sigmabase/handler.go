package sigmabase

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/core/models"
	"sigma-api/internal/shared/database"
	"sigma-api/internal/shared/utils"
	"time"
	"fmt"
	"os"
	"path/filepath"
)

type SigmabaseHandler struct {
	studentService StudentService
}

func NewSigmabaseHandler(s StudentService) *SigmabaseHandler {
	return &SigmabaseHandler{s}
}

// GetStudents returns a list of students with optional search
func (h *SigmabaseHandler) GetStudents(c fiber.Ctx) error {
	var students []models.Student
	query := database.DB.Preload("Classroom")

	// Pre-formatting for "Thin UI"
	// We can add search filters here (gender, class, name)
	if gender := c.Query("gender"); gender != "" {
		query = query.Where("gender = ?", gender)
	}

	if err := query.Find(&students).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(students)
}

// GetStudentDetails returns full detail for a specific student
func (h *SigmabaseHandler) GetStudentDetails(c fiber.Ctx) error {
	id := c.Params("id")
	student, err := h.studentService.GetStudentByID(utils.StringToUint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	// Siksa backend: Check permit status during detail retrieval
	isBlocked, permitStats, _ := h.studentService.IsPermitBlocked(student.ID)

	return c.JSON(fiber.Map{
		"student":      student,
		"is_blocked":    isBlocked,
		"permit_stats": permitStats,
	})
}

// GetTeachers returns a list of all teachers
func (h *SigmabaseHandler) GetTeachers(c fiber.Ctx) error {
	teachers, err := h.studentService.GetAllTeachers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(teachers)
}

// GetTeacherDetails returns full detail for a specific teacher
func (h *SigmabaseHandler) GetTeacherDetails(c fiber.Ctx) error {
	id := c.Params("id")
	teacher, err := h.studentService.GetTeacherByID(utils.StringToUint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Teacher not found"})
	}

	return c.JSON(teacher)
}

// GetAlumni returns a list of all alumni
func (h *SigmabaseHandler) GetAlumni(c fiber.Ctx) error {
	alumni, err := h.studentService.GetAllAlumni()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(alumni)
}

// GetClassrooms returns a list of all classrooms
func (h *SigmabaseHandler) GetClassrooms(c fiber.Ctx) error {
	classrooms, err := h.studentService.GetAllClassrooms()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(classrooms)
}
// GetStats returns summary statistics for the dashboard
func (h *SigmabaseHandler) GetStats(c fiber.Ctx) error {
	stats, err := h.studentService.GetStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stats)
}
type StudentDTO struct {
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Class       string `json:"class"`
	BirthPlace  string `json:"birth_place"`
	BirthDate   string `json:"birth_date"`
	NIK         string `json:"nik"`
	Address     string `json:"address"`
	ParentName  string `json:"parent_name"`
	ParentPhone string `json:"parent_phone"`
	NIS         string `json:"nis"`
	NISN        string `json:"nisn"`
	EntryYear   string `json:"entry_year"`
}

type TeacherDTO struct {
	Name           string `json:"name"`
	Gender         string `json:"gender"`
	NIK            string `json:"nik"`
	BirthPlace     string `json:"birth_place"`
	BirthDate      string `json:"birth_date"`
	NIP            string `json:"nip"`
	Position       string `json:"position"`
	Status         string `json:"status"`
	Education      string `json:"education"`
	GraduationYear string `json:"graduation_year"`
	Email          string `json:"email"`
	Address        string `json:"address"`
	WhatsApp       string `json:"whatsapp"`
	Photo          string `json:"photo"`
	EntryYear      string `json:"entry_year"`
	IsActive       *bool  `json:"is_active"`
}

type AlumniDTO struct {
	Name           string `json:"name"`
	Gender         string `json:"gender"`
	NIS            string `json:"nis"`
	BirthPlace     string `json:"birth_place"`
	BirthDate      string `json:"birth_date"`
	GraduationYear string `json:"graduation_year"`
	Batch          string `json:"batch"`
	NIK            string `json:"nik"`
	Email          string `json:"email"`
	Address        string `json:"address"`
	WhatsApp       string `json:"whatsapp"`
	Photo          string `json:"photo"`
	ServiceStatus  string `json:"service_status"`
}

// CreateStudent handles student registration
func (h *SigmabaseHandler) CreateStudent(c fiber.Ctx) error {
	var input StudentDTO
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid format: " + err.Error()})
	}

	student := models.Student{
		Name:      input.Name,
		Gender:    input.Gender,
		Class:     input.Class,
		EntryYear: input.EntryYear,
		NIS:       input.NIS,
	}
	if input.NISN != "" { student.NISN = &input.NISN }
	if input.NIK != "" { student.NIK = &input.NIK }
	if input.BirthPlace != "" { student.BirthPlace = &input.BirthPlace }
	if input.Address != "" { student.Address = &input.Address }
	if input.ParentName != "" { student.ParentName = &input.ParentName }
	if input.ParentPhone != "" { student.ParentPhone = &input.ParentPhone }
	if input.BirthDate != "" {
		t, err := time.Parse("2006-01-02", input.BirthDate)
		if err == nil { student.BirthDate = &t }
	}

	if err := h.studentService.CreateStudent(&student); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create student: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(student)
}

// BulkCreateStudents handles batch registration
func (h *SigmabaseHandler) BulkCreateStudents(c fiber.Ctx) error {
	var inputs []StudentDTO
	if err := c.Bind().JSON(&inputs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid format: " + err.Error()})
	}

	var students []models.Student
	for _, input := range inputs {
		student := models.Student{
			Name:      input.Name,
			Gender:    input.Gender,
			Class:     input.Class,
			EntryYear: input.EntryYear,
			NIS:       input.NIS,
		}
		if input.NISN != "" { n := input.NISN; student.NISN = &n }
		if input.NIK != "" { nik := input.NIK; student.NIK = &nik }
		if input.BirthPlace != "" { bp := input.BirthPlace; student.BirthPlace = &bp }
		if input.Address != "" { addr := input.Address; student.Address = &addr }
		if input.ParentName != "" { pn := input.ParentName; student.ParentName = &pn }
		if input.ParentPhone != "" { pp := input.ParentPhone; student.ParentPhone = &pp }
		if input.BirthDate != "" {
			t, err := time.Parse("2006-01-02", input.BirthDate)
			if err == nil { student.BirthDate = &t }
		}
		students = append(students, student)
	}

	if err := h.studentService.BulkCreateStudents(students); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to bulk create: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Bulk import successful", "count": len(students)})
}


// UpdateStudent handles student information updates
func (h *SigmabaseHandler) UpdateStudent(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	var input StudentDTO
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid format: " + err.Error()})
	}

	student := models.Student{
		Name:      input.Name,
		Gender:    input.Gender,
		Class:     input.Class,
		EntryYear: input.EntryYear,
		NIS:       input.NIS,
	}
	// Use pointers where GORM expects them
	if input.NISN != "" { student.NISN = &input.NISN }
	if input.NIK != "" { student.NIK = &input.NIK }
	if input.BirthPlace != "" { student.BirthPlace = &input.BirthPlace }
	if input.Address != "" { student.Address = &input.Address }
	if input.ParentName != "" { student.ParentName = &input.ParentName }
	if input.ParentPhone != "" { student.ParentPhone = &input.ParentPhone }
	if input.BirthDate != "" {
		t, err := time.Parse("2006-01-02", input.BirthDate)
		if err == nil { student.BirthDate = &t }
	}

	if err := h.studentService.UpdateStudent(id, &student); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update student: " + err.Error()})
	}

	return c.JSON(student)
}

// DeleteStudent removes a student record
func (h *SigmabaseHandler) DeleteStudent(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	if err := h.studentService.DeleteStudent(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete student"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// CreateTeacher handles teacher registration
func (h *SigmabaseHandler) CreateTeacher(c fiber.Ctx) error {
	var input TeacherDTO
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid format: " + err.Error()})
	}

	teacher := models.Teacher{
		Name:           input.Name,
		Gender:         input.Gender,
		NIK:            input.NIK,
		IsActive:       true,
	}
	if input.IsActive != nil { teacher.IsActive = *input.IsActive }
	if input.NIP != "" { teacher.NIP = &input.NIP }
	if input.BirthPlace != "" { teacher.BirthPlace = &input.BirthPlace }
	if input.Position != "" { teacher.Position = &input.Position }
	if input.Status != "" { teacher.Status = &input.Status }
	if input.Education != "" { teacher.Education = &input.Education }
	if input.GraduationYear != "" { teacher.GraduationYear = &input.GraduationYear }
	if input.Email != "" { teacher.Email = &input.Email }
	if input.Address != "" { teacher.Address = &input.Address }
	if input.WhatsApp != "" { teacher.WhatsApp = &input.WhatsApp }
	if input.Photo != "" { teacher.Photo = &input.Photo }
	if input.EntryYear != "" { teacher.EntryYear = &input.EntryYear }
	
	if input.BirthDate != "" {
		t, err := time.Parse("2006-01-02", input.BirthDate)
		if err == nil { teacher.BirthDate = &t }
	}

	if err := h.studentService.CreateTeacher(&teacher); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create teacher: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(teacher)
}

// BulkCreateTeachers handles batch teacher registration
func (h *SigmabaseHandler) BulkCreateTeachers(c fiber.Ctx) error {
	var inputs []TeacherDTO
	if err := c.Bind().JSON(&inputs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid format: " + err.Error()})
	}

	var teachers []models.Teacher
	for _, input := range inputs {
		teacher := models.Teacher{
			Name:           input.Name,
			Gender:         input.Gender,
			NIK:            input.NIK,
			IsActive:       true,
		}
		if input.IsActive != nil { teacher.IsActive = *input.IsActive }
		if input.NIP != "" { n := input.NIP; teacher.NIP = &n }
		if input.BirthPlace != "" { bp := input.BirthPlace; teacher.BirthPlace = &bp }
		if input.Position != "" { p := input.Position; teacher.Position = &p }
		if input.Status != "" { s := input.Status; teacher.Status = &s }
		if input.Education != "" { ed := input.Education; teacher.Education = &ed }
		if input.GraduationYear != "" { gy := input.GraduationYear; teacher.GraduationYear = &gy }
		if input.Email != "" { em := input.Email; teacher.Email = &em }
		if input.Address != "" { addr := input.Address; teacher.Address = &addr }
		if input.WhatsApp != "" { wa := input.WhatsApp; teacher.WhatsApp = &wa }
		if input.Photo != "" { ph := input.Photo; teacher.Photo = &ph }
		if input.EntryYear != "" { ey := input.EntryYear; teacher.EntryYear = &ey }
		
		if input.BirthDate != "" {
			t, err := time.Parse("2006-01-02", input.BirthDate)
			if err == nil { teacher.BirthDate = &t }
		}
		teachers = append(teachers, teacher)
	}

	if err := h.studentService.BulkCreateTeachers(teachers); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to bulk create: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Bulk import successful", "count": len(teachers)})
}

// UpdateTeacher handles teacher information updates
func (h *SigmabaseHandler) UpdateTeacher(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	var input TeacherDTO
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid format: " + err.Error()})
	}

	// Build surgical update map to prevent data loss and fix column naming
	updateMap := make(map[string]interface{})
	
	if input.IsActive != nil { updateMap["is_active"] = *input.IsActive }
	if input.Name != "" { updateMap["name"] = input.Name }
	if input.Gender != "" { updateMap["gender"] = input.Gender }
	if input.NIK != "" { updateMap["nik"] = input.NIK }
	if input.NIP != "" { updateMap["nip"] = input.NIP }
	if input.BirthPlace != "" { updateMap["birth_place"] = input.BirthPlace }
	if input.Position != "" { updateMap["position"] = input.Position }
	if input.Status != "" { updateMap["status"] = input.Status }
	if input.Education != "" { updateMap["education"] = input.Education }
	if input.GraduationYear != "" { updateMap["graduation_year"] = input.GraduationYear }
	if input.Email != "" { updateMap["email"] = input.Email }
	if input.Address != "" { updateMap["address"] = input.Address }
	if input.WhatsApp != "" { updateMap["whats_app"] = input.WhatsApp } // FIX: Case-sensitive DB column
	if input.Photo != "" { updateMap["photo"] = input.Photo }
	if input.EntryYear != "" { updateMap["entry_year"] = input.EntryYear }
	
	if input.BirthDate != "" {
		t, err := time.Parse("2006-01-02", input.BirthDate)
		if err == nil { updateMap["birth_date"] = t }
	}

	if len(updateMap) > 0 {
		if err := h.studentService.UpdateTeacher(id, updateMap); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update teacher: " + err.Error()})
		}
	}

	// Fetch fresh data to ensure response is accurate
	updatedTeacher, _ := h.studentService.GetTeacherByID(id)

	return c.JSON(updatedTeacher)
}

// CreateAlumni handles manual alumni registration
func (h *SigmabaseHandler) CreateAlumni(c fiber.Ctx) error {
	var input AlumniDTO
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid format: " + err.Error()})
	}

	alumni := models.Alumni{
		Name:           input.Name,
		Gender:         input.Gender,
		GraduationYear: input.GraduationYear,
		Batch:          input.Batch,
		ServiceStatus:  input.ServiceStatus,
	}
	if input.NIS != "" { alumni.NIS = &input.NIS }
	if input.NIK != "" { 
		alumni.NIK = &input.NIK 
	} else {
		alumni.NIK = nil
	}
	if input.BirthPlace != "" { alumni.BirthPlace = &input.BirthPlace }
	if input.Email != "" { alumni.Email = &input.Email }
	if input.Address != "" { alumni.Address = &input.Address }
	if input.WhatsApp != "" { alumni.WhatsApp = &input.WhatsApp }
	if input.Photo != "" { alumni.Photo = &input.Photo }
	
	if input.BirthDate != "" {
		t, err := time.Parse("2006-01-02", input.BirthDate)
		if err == nil { alumni.BirthDate = &t }
	}

	if err := h.studentService.CreateAlumni(&alumni); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create alumni: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(alumni)
}

// UpdateAlumni handles alumni information updates
func (h *SigmabaseHandler) UpdateAlumni(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	var input AlumniDTO
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid format: " + err.Error()})
	}

	updateMap := make(map[string]interface{})
	if input.Name != "" { updateMap["name"] = input.Name }
	if input.Gender != "" { updateMap["gender"] = input.Gender }
	if input.NIK != "" { 
		updateMap["nik"] = input.NIK 
	} else {
		updateMap["nik"] = nil
	}
	if input.NIS != "" { updateMap["nis"] = input.NIS }
	if input.GraduationYear != "" { updateMap["graduation_year"] = input.GraduationYear }
	if input.Batch != "" { updateMap["batch"] = input.Batch }
	if input.Email != "" { updateMap["email"] = input.Email }
	if input.Address != "" { updateMap["address"] = input.Address }
	if input.WhatsApp != "" { updateMap["whats_app"] = input.WhatsApp }
	if input.ServiceStatus != "" { updateMap["service_status"] = input.ServiceStatus }
	
	if input.BirthDate != "" {
		t, err := time.Parse("2006-01-02", input.BirthDate)
		if err == nil { updateMap["birth_date"] = t }
	}

	if len(updateMap) > 0 {
		if err := h.studentService.UpdateAlumni(id, updateMap); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update alumni: " + err.Error()})
		}
	}

	updatedAlumni, _ := h.studentService.GetAlumniByID(id)
	return c.JSON(updatedAlumni)
}

// DeleteAlumni removes an alumni record
func (h *SigmabaseHandler) DeleteAlumni(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	if err := h.studentService.DeleteAlumni(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete alumni"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GraduateStudent handles the mutation from active student to alumni
func (h *SigmabaseHandler) GraduateStudent(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	var input AlumniDTO
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid graduation data"})
	}

	// Prepare migration data for the service
	migrationData := map[string]interface{}{
		"graduation_year": input.GraduationYear,
		"batch":           input.Batch,
		"whatsapp":        input.WhatsApp,
		"nik":             input.NIK,
		"email":           input.Email,
		"photo":           input.Photo,
		"service_status":  input.ServiceStatus,
	}

	newAlumniID, err := h.studentService.GraduateStudent(id, migrationData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Graduation process failed: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success", 
		"message": "Student has graduated successfully",
		"alumni_id": newAlumniID,
	})
}

// UploadAlumniPhoto handles alumni photo upload
func (h *SigmabaseHandler) UploadAlumniPhoto(c fiber.Ctx) error {
	id := c.Params("id")
	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	uploadDir := "./uploads/alumni"
	os.MkdirAll(uploadDir, 0755)
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", id, time.Now().Unix(), ext)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	dbPath := fmt.Sprintf("/uploads/alumni/%s", filename)
	if err := h.studentService.UpdateAlumni(utils.StringToUint(id), map[string]interface{}{"photo": dbPath}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update db record"})
	}

	return c.JSON(fiber.Map{"status": "success", "path": dbPath})
}

// DeleteTeacher removes a teacher record
func (h *SigmabaseHandler) DeleteTeacher(c fiber.Ctx) error {
	id := utils.StringToUint(c.Params("id"))
	if err := h.studentService.DeleteTeacher(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete teacher"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// UploadTeacherPhoto handles personnel photo upload
func (h *SigmabaseHandler) UploadTeacherPhoto(c fiber.Ctx) error {
	id := c.Params("id")
	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	// 1. Validate File Size (Max 2MB)
	if file.Size > 2*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File too large (max 2MB)"})
	}

	// 2. Validate Extension
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Only JPG/PNG images allowed"})
	}

	// 3. Ensure Directory Exists
	uploadDir := "./uploads/teachers"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create upload directory"})
	}

	// 4. Create Unique Filename
	filename := fmt.Sprintf("%s_%d%s", id, time.Now().Unix(), ext)
	savePath := filepath.Join(uploadDir, filename)

	// 5. Save File
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// 6. Update Database using Map-based Service
	dbPath := fmt.Sprintf("/uploads/teachers/%s", filename)
	if err := h.studentService.UpdateTeacher(utils.StringToUint(id), map[string]interface{}{"photo": dbPath}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update db record"})
	}

	return c.JSON(fiber.Map{"status": "success", "path": dbPath})
}
