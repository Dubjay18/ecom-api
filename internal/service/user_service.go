package service

import "github.com/Dubjay18/ecom-api/internal/repository"

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}
