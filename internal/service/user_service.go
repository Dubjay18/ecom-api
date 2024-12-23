package service

import (
	"context"

	"github.com/Dubjay18/ecom-api/internal/domain"
	"github.com/Dubjay18/ecom-api/internal/repository"
)

type UserService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func (s *userService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, error) {
	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	return s.repo.Create(ctx, user)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
