package model

import "time"

type Submission struct {
	Model
	AssignmentID    int
	StudentID       int
	SubmittedAt     time.Time
	StudentName     string
	AssignmentTitle string
	FileURL         string
	Status          string
	Grade           *float64
}
