package service

import (
	"errors"
	"session-17/model"
	"session-17/repository"
)

type AuthService interface {
	Login(email, password string) (*model.User, error)
}

type authService struct {
	Repo repository.Repository
}

func NewAuthService(repo repository.Repository) AuthService {
	return &authService{Repo: repo}
}

func (s *authService) Login(email, password string) (*model.User, error) {
	user, err := s.Repo.UserRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Password != password {
		return nil, errors.New("incorrect password")
	}

	return user, nil
}
