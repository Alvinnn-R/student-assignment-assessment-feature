package repository

import (
	"context"
	"database/sql"
	"errors"
	"session-17/database"
	"session-17/model"
)

type AssignmentRepository interface {
	Create(assignment *model.Assignment) error
	FindByID(id int) (*model.Assignment, error)
	FindAll() ([]model.Assignment, error)
	Update(assignment *model.Assignment) error
	Delete(id int) error
}

type assignmentRepository struct {
	db database.PgxIface
}

func NewAssignmentRepository(db database.PgxIface) AssignmentRepository {
	return &assignmentRepository{db: db}
}

func (r *assignmentRepository) Create(a *model.Assignment) error {
	query := `
		INSERT INTO assignments (course_id, lecturer_id, title, description, deadline, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id
	`
	err := r.db.QueryRow(context.Background(), query, a.CourseID, a.LecturerID, a.Title, a.Description, a.Deadline).Scan(&a.ID)
	return err
}

func (r *assignmentRepository) FindByID(id int) (*model.Assignment, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at, course_id, lecturer_id, title, description, deadline
		FROM assignments WHERE id = $1 AND deleted_at IS NULL
	`
	var a model.Assignment
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&a.ID, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
		&a.CourseID, &a.LecturerID, &a.Title, &a.Description, &a.Deadline,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *assignmentRepository) FindAll() ([]model.Assignment, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at, course_id, lecturer_id, title, description, deadline
		FROM assignments WHERE deleted_at IS NULL
		ORDER BY deadline ASC
	`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []model.Assignment
	for rows.Next() {
		var a model.Assignment
		err := rows.Scan(
			&a.ID, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
			&a.CourseID, &a.LecturerID, &a.Title, &a.Description, &a.Deadline,
		)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, nil
}

func (r *assignmentRepository) Update(a *model.Assignment) error {
	query := `
		UPDATE assignments
		SET course_id = $1, lecturer_id = $2, title = $3, description = $4,
		    deadline = $5, updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(context.Background(), query, a.CourseID, a.LecturerID, a.Title, a.Description, a.Deadline, a.ID)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("assignment not found or already deleted")
	}
	return err
}

func (r *assignmentRepository) Delete(id int) error {
	query := `
		UPDATE assignments SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("assignment not found or already deleted")
	}
	return err
}
