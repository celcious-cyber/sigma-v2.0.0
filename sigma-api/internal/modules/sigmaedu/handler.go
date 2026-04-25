package sigmaedu

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/core/models"
)


type EduHandler struct {
	service EduService
}

func NewEduHandler(s EduService) *EduHandler {
	return &EduHandler{s}
}

// GetStats returns summary statistics for academic module
func (h *EduHandler) GetStats(c fiber.Ctx) error {
	filter := c.Query("filter", "week")
	stats, err := h.service.GetStats(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stats)
}

// GetSubjects returns a list of all subjects
func (h *EduHandler) GetSubjects(c fiber.Ctx) error {
	subjects, err := h.service.GetAllSubjects()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(subjects)
}


// CreateSubject adds a new subject
func (h *EduHandler) CreateSubject(c fiber.Ctx) error {
	var sub models.Subject
	if err := c.Bind().JSON(&sub); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.CreateSubject(&sub); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Subject created successfully", "data": sub})
}

// UpdateSubject updates an existing subject
func (h *EduHandler) UpdateSubject(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var sub models.Subject
	if err := c.Bind().JSON(&sub); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.UpdateSubject(uint(id), &sub); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Subject updated successfully"})
}

// DeleteSubject removes a subject
func (h *EduHandler) DeleteSubject(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.DeleteSubject(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Subject deleted successfully"})
}

// GetStudyHours returns a list of all lesson time slots
func (h *EduHandler) GetStudyHours(c fiber.Ctx) error {
	hours, err := h.service.GetAllStudyHours()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(hours)
}

// CreateStudyHour adds a new lesson time slot
func (h *EduHandler) CreateStudyHour(c fiber.Ctx) error {
	var hour models.StudyHour
	if err := c.Bind().JSON(&hour); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.CreateStudyHour(&hour); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Study hour created successfully", "data": hour})
}

// UpdateStudyHour updates an existing lesson time slot
func (h *EduHandler) UpdateStudyHour(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var hour models.StudyHour
	if err := c.Bind().JSON(&hour); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.UpdateStudyHour(uint(id), &hour); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Study hour updated successfully"})
}

// DeleteStudyHour removes a lesson time slot
func (h *EduHandler) DeleteStudyHour(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.DeleteStudyHour(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Study hour deleted successfully"})
}

// Academic Calendar
func (h *EduHandler) GetCalendarEvents(c fiber.Ctx) error {
	events, err := h.service.GetAllCalendarEvents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(events)
}

func (h *EduHandler) CreateCalendarEvent(c fiber.Ctx) error {
	var event models.AcademicCalendar
	if err := c.Bind().JSON(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.CreateCalendarEvent(&event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Event created successfully", "data": event})
}

func (h *EduHandler) UpdateCalendarEvent(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var event models.AcademicCalendar
	if err := c.Bind().JSON(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.UpdateCalendarEvent(uint(id), &event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Event updated successfully"})
}

func (h *EduHandler) DeleteCalendarEvent(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.DeleteCalendarEvent(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Event deleted successfully"})
}

// Classrooms
func (h *EduHandler) GetClassrooms(c fiber.Ctx) error {
	classrooms, err := h.service.GetAllClassrooms()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(classrooms)
}

func (h *EduHandler) CreateClassroom(c fiber.Ctx) error {
	var classroom models.Classroom
	if err := c.Bind().JSON(&classroom); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.CreateClassroom(&classroom); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Classroom created successfully", "data": classroom})
}

func (h *EduHandler) UpdateClassroom(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var classroom models.Classroom
	if err := c.Bind().JSON(&classroom); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.UpdateClassroom(uint(id), &classroom); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Classroom updated successfully"})
}

func (h *EduHandler) DeleteClassroom(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.DeleteClassroom(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Classroom deleted successfully"})
}



// GetSchedules returns the complete timetable
func (h *EduHandler) GetSchedules(c fiber.Ctx) error {
	schedules, err := h.service.GetSchedules()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(schedules)
}

func (h *EduHandler) CreateSchedule(c fiber.Ctx) error {
	var schedule models.TeacherSchedule
	if err := c.Bind().JSON(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.CreateSchedule(&schedule); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Schedule created successfully", "data": schedule})
}

func (h *EduHandler) UpdateSchedule(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var schedule models.TeacherSchedule
	if err := c.Bind().JSON(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := h.service.UpdateSchedule(uint(id), &schedule); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Schedule updated successfully"})
}

func (h *EduHandler) DeleteSchedule(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.DeleteSchedule(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Schedule deleted successfully"})
}

// Attendance Handlers
func (h *EduHandler) GetAttendances(c fiber.Ctx) error {
	classID, _ := strconv.Atoi(c.Query("classroom_id"))
	date := c.Query("date") // Format: YYYY-MM-DD
	var subjectID *uint
	if sID := c.Query("subject_id"); sID != "" {
		val, _ := strconv.Atoi(sID)
		uVal := uint(val)
		subjectID = &uVal
	}

	attendances, err := h.service.GetAttendances(uint(classID), date, subjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(attendances)
}

func (h *EduHandler) GetAttendanceReport(c fiber.Ctx) error {
	classID, _ := strconv.Atoi(c.Query("classroom_id"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	
	var subjectID *uint
	if sID := c.Query("subject_id"); sID != "" {
		val, _ := strconv.Atoi(sID)
		uVal := uint(val)
		subjectID = &uVal
	}

	report, err := h.service.GetAttendanceReport(uint(classID), subjectID, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(report)
}

func (h *EduHandler) BulkRecordAttendance(c fiber.Ctx) error {
	var attendances []models.Attendance
	if err := c.Bind().JSON(&attendances); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	// Get User ID from JWT context (provided by middleware)
	userID := c.Locals("user_id").(uint)

	for i := range attendances {
		attendances[i].RecordedBy = userID
	}

	if err := h.service.BulkRecordAttendance(attendances); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Attendance recorded successfully"})
}

func (h *EduHandler) DeleteAttendance(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.DeleteAttendance(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Attendance deleted successfully"})
}

func (h *EduHandler) GetStudentsByClass(c fiber.Ctx) error {
	classID, _ := strconv.Atoi(c.Params("id"))
	students, err := h.service.GetStudentsByClass(uint(classID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(students)
}

func (h *EduHandler) GetUnassignedStudents(c fiber.Ctx) error {
	students, err := h.service.GetUnassignedStudents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(students)
}

func (h *EduHandler) BulkUpdateStudentClass(c fiber.Ctx) error {
	type Request struct {
		StudentIDs  []uint `json:"student_ids"`
		ClassroomID *uint  `json:"classroom_id"`
	}
	var req Request
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.service.BulkUpdateStudentClass(req.StudentIDs, req.ClassroomID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Students updated successfully"})
}

// Assessment Handlers
func (h *EduHandler) GetAssessments(c fiber.Ctx) error {
	classID, _ := strconv.Atoi(c.Query("classroom_id"))
	subjectID, _ := strconv.Atoi(c.Query("subject_id"))
	term := c.Query("term")
	year := c.Query("academic_year")
	aType := c.Query("type")

	assessments, err := h.service.GetAssessments(uint(classID), uint(subjectID), term, year, aType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(assessments)
}

func (h *EduHandler) BulkRecordAssessments(c fiber.Ctx) error {
	var assessments []models.Assessment
	if err := c.Bind().JSON(&assessments); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	userID := c.Locals("user_id").(uint)
	for i := range assessments {
		assessments[i].TeacherID = userID
	}

	if err := h.service.BulkRecordAssessments(assessments); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Assessments recorded successfully"})
}

func (h *EduHandler) DeleteAssessments(c fiber.Ctx) error {
	classID, _ := strconv.Atoi(c.Query("classroom_id"))
	subjectID, _ := strconv.Atoi(c.Query("subject_id"))
	term := c.Query("term")
	year := c.Query("academic_year")
	aType := c.Query("type")

	if err := h.service.DeleteAssessments(uint(classID), uint(subjectID), term, year, aType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Assessments deleted successfully"})
}

// Tahfidz Handlers
func (h *EduHandler) GetTahfidzRecords(c fiber.Ctx) error {
	classID, _ := strconv.Atoi(c.Query("classroom_id"))
	date := c.Query("date")

	records, err := h.service.GetTahfidzRecords(uint(classID), date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(records)
}

func (h *EduHandler) BulkRecordTahfidz(c fiber.Ctx) error {
	var records []models.QuranMemorization
	if err := c.Bind().JSON(&records); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	userID := c.Locals("user_id").(uint)
	for i := range records {
		records[i].TeacherID = userID
	}

	if err := h.service.BulkRecordTahfidz(records); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Tahfidz records saved successfully"})
}

func (h *EduHandler) DeleteTahfidz(c fiber.Ctx) error {
	classID, _ := strconv.Atoi(c.Query("classroom_id"))
	date := c.Query("date")
	tType := c.Query("type")

	if err := h.service.DeleteTahfidz(uint(classID), date, tType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Tahfidz records deleted successfully"})
}

func (h *EduHandler) GetTahfidzHistory(c fiber.Ctx) error {
	studentID, _ := strconv.Atoi(c.Params("student_id"))
	records, err := h.service.GetTahfidzHistoryByStudent(uint(studentID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(records)
}

func (h *EduHandler) DeleteTahfidzRecord(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.DeleteTahfidzRecord(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Record deleted successfully"})
}

// Lesson Memorization Handlers
func (h *EduHandler) GetLessonMemorizations(c fiber.Ctx) error {
	classID, _ := strconv.Atoi(c.Query("classroom_id"))
	date := c.Query("date")

	records, err := h.service.GetLessonMemorizations(uint(classID), date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(records)
}

func (h *EduHandler) BulkRecordLessonMemorizations(c fiber.Ctx) error {
	var records []models.LessonMemorization
	if err := c.Bind().JSON(&records); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	userID := c.Locals("user_id").(uint)
	for i := range records {
		records[i].TeacherID = userID
	}

	if err := h.service.BulkRecordLessonMemorizations(records); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Lesson memorization records saved successfully"})
}

func (h *EduHandler) GetLessonMemorizationHistory(c fiber.Ctx) error {
	studentID, _ := strconv.Atoi(c.Params("student_id"))
	records, err := h.service.GetLessonMemorizationHistory(uint(studentID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(records)
}

func (h *EduHandler) DeleteLessonMemorizationRecord(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.DeleteLessonMemorizationRecord(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Record deleted successfully"})
}

// Teacher Attendance Handlers
func (h *EduHandler) GetTeacherAttendance(c fiber.Ctx) error {
	date := c.Query("date")
	records, err := h.service.GetTeacherAttendance(date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(records)
}

func (h *EduHandler) BulkRecordTeacherAttendance(c fiber.Ctx) error {
	var records []models.TeacherAttendance
	if err := c.Bind().JSON(&records); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.service.BulkRecordTeacherAttendance(records); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Teacher attendance records saved successfully"})
}

// Teaching Journal Handlers
func (h *EduHandler) GetTeachingJournals(c fiber.Ctx) error {
	date := c.Query("date")
	
	var classroomID *uint
	if cID := c.Query("classroom_id"); cID != "" {
		val, _ := strconv.Atoi(cID)
		uVal := uint(val)
		classroomID = &uVal
	}
	
	var subjectID *uint
	if sID := c.Query("subject_id"); sID != "" {
		val, _ := strconv.Atoi(sID)
		uVal := uint(val)
		subjectID = &uVal
	}

	records, err := h.service.GetTeachingJournals(date, classroomID, subjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(records)
}

func (h *EduHandler) CreateTeachingJournal(c fiber.Ctx) error {
	var journal models.TeachingJournal
	if err := c.Bind().JSON(&journal); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	// Set teacher from session if not provided
	if journal.TeacherID == 0 {
		journal.TeacherID = c.Locals("user_id").(uint)
	}

	if err := h.service.CreateTeachingJournal(&journal); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(journal)
}

func (h *EduHandler) UpdateTeachingJournal(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var journal models.TeachingJournal
	if err := c.Bind().JSON(&journal); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.service.UpdateTeachingJournal(uint(id), &journal); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Journal updated successfully"})
}

func (h *EduHandler) DeleteTeachingJournal(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.DeleteTeachingJournal(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Journal deleted successfully"})
}

func (h *EduHandler) GetLastJournal(c fiber.Ctx) error {
	teacherID := c.Locals("user_id").(uint)
	subjectID, _ := strconv.Atoi(c.Query("subject_id"))

	journal, err := h.service.GetLastJournal(teacherID, uint(subjectID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No journal found"})
	}
	return c.JSON(journal)
}

