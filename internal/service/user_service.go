package service

import (
	"context"
	"github.com/Dubjay18/ecom-api/internal/domain"
	"github.com/Dubjay18/ecom-api/internal/repository"
	"github.com/Dubjay18/ecom-api/internal/util"
	"github.com/Dubjay18/ecom-api/pkg/common"
	"github.com/Dubjay18/ecom-api/pkg/jwt"
	"log"
	"net/http"
)

var (
	ErrUserNotFound = &common.AppError{
		Code:    http.StatusNotFound,
		Message: "User not found",
	}
)

type UserService interface {
	// Register creates a new user
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, *common.AppError)
	// Login logs in a user
	Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, *common.AppError)
	// GetByID returns a user by ID
	GetByID(ctx context.Context, id uint) (*domain.User, *common.AppError)
	// Update updates a user
	Update(ctx context.Context, user *domain.User) (*domain.User, *common.AppError)
}

type userService struct {
	repo repository.UserRepository
	jwt  *jwt.JWTService
}

func (s *userService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, *common.AppError) {
	// Check if user already exists
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, &common.ErrEmailExists
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return nil, &common.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to hash password",
		}
	}
	user := &domain.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      domain.RoleUser,
	}

	resp, err := s.repo.Create(ctx, user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil, &common.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
		}
	}

	return resp, nil
}

func (s *userService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, *common.AppError) {
	//
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil, ErrUserNotFound
	}

	if !util.CheckPassword(req.Password, user.Password) {
		return nil, &common.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.jwt.GenerateToken(user)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return nil, &common.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate token",
		}
	}

	user.Orders = nil
	user.Addresses = nil
	return &domain.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*domain.User, *common.AppError) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, user *domain.User) (*domain.User, *common.AppError) {
	updated, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, &common.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update user",
		}
	}
	return updated, nil
}

func NewUserService(repo repository.UserRepository, jwt *jwt.JWTService) UserService {
	return &userService{repo: repo,
		jwt: jwt,
	}
}
