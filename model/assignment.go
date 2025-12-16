package model

import "time"

type Assignment struct {
	Model
	CourseID    int
	LecturerID  int
	Title       string
	Description string
	Deadline    time.Time
}
