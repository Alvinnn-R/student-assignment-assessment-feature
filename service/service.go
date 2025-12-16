package service

import "session-17/repository"

type Service struct {
	AssignmentService AssignmentService
	SubmissionService SubmissionService
	UserService       UserService
	AuthService       AuthService
}

func NewService(repo repository.Repository) Service {
	return Service{
		AssignmentService: NewAssignmentService(repo),
		SubmissionService: NewSubmissionService(repo),
		UserService:       NewUserService(repo),
		AuthService:       NewAuthService(repo),
	}
}
