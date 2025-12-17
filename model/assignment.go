package model

import "time"

type Assignment struct {
	Model
	CourseID    int
	LecturerID  int
	Title       string
	Description string
	Deadline    time.Time
	Grade       *float64 // Nilai siswa untuk assignment ini (jika sudah ada)
	Status      string   // Status submission (submitted/late)
}
