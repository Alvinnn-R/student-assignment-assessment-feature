package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"session-17/model"
	"session-17/repository"
	"time"
)

type AssignmentService interface {
	GetAllAssignments() ([]model.Assignment, error)
	GetAssignmentsWithGrades(studentID int) ([]model.Assignment, error)
	SubmitAssignment(studentID, assignmentID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	GetGradeFormData() ([]model.User, []model.Assignment, error)
	GetAssignmentByID(id int) (*model.Assignment, error)
}

type assignmentService struct {
	Repo repository.Repository
}

func NewAssignmentService(repo repository.Repository) AssignmentService {
	return &assignmentService{Repo: repo}
}

func (s *assignmentService) GetAllAssignments() ([]model.Assignment, error) {
	return s.Repo.AssignmentRepo.FindAll()
}

func (s *assignmentService) GetAssignmentByID(id int) (*model.Assignment, error) {
	return s.Repo.AssignmentRepo.FindByID(id)
}

func (s *assignmentService) SubmitAssignment(studentID, assignmentID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	assignment, err := s.Repo.AssignmentRepo.FindByID(assignmentID)
	if err != nil {
		return "", err
	}

	count, err := s.Repo.SubmissionRepo.CountByStudentAndAssignment(studentID, assignmentID)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "already submitted", nil
	}

	// save file to disk
	uploadDir := "public/uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	filename := fmt.Sprintf("%d_%d_%s", assignmentID, studentID, fileHeader.Filename)
	filepath := fmt.Sprintf("%s/%s", uploadDir, filename)
	accessURL := fmt.Sprintf("/public/uploads/%s", filename)

	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	status := "submitted"
	if time.Now().After(assignment.Deadline) {
		status = "late"
	}

	sub := &model.Submission{
		AssignmentID: assignmentID,
		StudentID:    studentID,
		SubmittedAt:  time.Now(),
		FileURL:      accessURL,
		Status:       status,
	}

	return status, s.Repo.SubmissionRepo.Create(sub)
}
func (s *assignmentService) GetAssignmentsWithGrades(studentID int) ([]model.Assignment, error) {
	// Get all assignments
	assignments, err := s.Repo.AssignmentRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// For each assignment, check if student has submission and grade
	for i := range assignments {
		submission, err := s.Repo.SubmissionRepo.FindByStudentAndAssignment(studentID, assignments[i].ID)
		if err == nil && submission != nil {
			assignments[i].Grade = submission.Grade
			assignments[i].Status = submission.Status
		}
	}

	return assignments, nil
}
func (s *assignmentService) GetGradeFormData() ([]model.User, []model.Assignment, error) {
	students, err := s.Repo.UserRepo.FindAllStudents()
	if err != nil {
		return nil, nil, err
	}

	assignments, err := s.Repo.AssignmentRepo.FindAll()
	if err != nil {
		return nil, nil, err
	}

	return students, assignments, nil
}
