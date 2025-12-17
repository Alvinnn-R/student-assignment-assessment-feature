package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"session-17/model"
	"session-17/service"
	"strconv"
)

type SubmissionHandler struct {
	SubmissionService service.SubmissionService
	Templates         *template.Template
}

func NewSubmissionHandler(templates *template.Template, submissionService service.SubmissionService) SubmissionHandler {
	return SubmissionHandler{
		SubmissionService: submissionService,
		Templates:         templates,
	}
}

// ListSubmissions - Menampilkan semua submission untuk dosen
func (h *SubmissionHandler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
	submissions, err := h.SubmissionService.GetAllSubmissions()
	if err != nil {
		http.Error(w, "Failed to get submissions", http.StatusInternalServerError)
		return
	}

	// Cek query parameter success
	showSuccess := r.URL.Query().Get("success") == "1"

	// Data untuk template
	data := struct {
		Submissions []model.Submission
		ShowSuccess bool
	}{
		Submissions: submissions,
		ShowSuccess: showSuccess,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.Templates.ExecuteTemplate(w, "submissions_list", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GradeFormView - Menampilkan form untuk input nilai
func (h *SubmissionHandler) GradeFormView(w http.ResponseWriter, r *http.Request) {
	studentIDStr := r.URL.Query().Get("student_id")
	assignmentIDStr := r.URL.Query().Get("assignment_id")

	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	assignmentID, err := strconv.Atoi(assignmentIDStr)
	if err != nil {
		http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
		return
	}

	// Ambil detail submission
	submissions, err := h.SubmissionService.GetAllSubmissions()
	if err != nil {
		http.Error(w, "Failed to get submission", http.StatusInternalServerError)
		return
	}

	// Cari submission yang sesuai
	var targetSubmission *struct {
		StudentID       int
		AssignmentID    int
		StudentName     string
		AssignmentTitle string
		Status          string
		FileURL         string
		Grade           *float64
	}

	for _, sub := range submissions {
		if sub.StudentID == studentID && sub.AssignmentID == assignmentID {
			targetSubmission = &struct {
				StudentID       int
				AssignmentID    int
				StudentName     string
				AssignmentTitle string
				Status          string
				FileURL         string
				Grade           *float64
			}{
				StudentID:       sub.StudentID,
				AssignmentID:    sub.AssignmentID,
				StudentName:     sub.StudentName,
				AssignmentTitle: sub.AssignmentTitle,
				Status:          sub.Status,
				FileURL:         sub.FileURL,
				Grade:           sub.Grade,
			}
			break
		}
	}

	if targetSubmission == nil {
		http.Error(w, "Submission not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.Templates.ExecuteTemplate(w, "grade_form", targetSubmission); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SubmitGrade - Proses submit nilai dari form
func (h *SubmissionHandler) SubmitGrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(r.FormValue("student_id"))
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	assignmentID, err := strconv.Atoi(r.FormValue("assignment_id"))
	if err != nil {
		http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
		return
	}

	grade, err := strconv.ParseFloat(r.FormValue("grade"), 64)
	if err != nil {
		http.Error(w, "Invalid grade value", http.StatusBadRequest)
		return
	}

	// Validasi nilai harus 0-100
	if grade < 0 || grade > 100 {
		http.Error(w, "Grade must be between 0 and 100", http.StatusBadRequest)
		return
	}

	// Update nilai di database
	err = h.SubmissionService.GradeSubmission(studentID, assignmentID, grade)
	if err != nil {
		fmt.Println("Error grading:", err)
		http.Error(w, "Failed to submit grade", http.StatusInternalServerError)
		return
	}

	// Redirect ke list submissions dengan pesan sukses
	http.Redirect(w, r, "/user/submissions?success=1", http.StatusSeeOther)
}
