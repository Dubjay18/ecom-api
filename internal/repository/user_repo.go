package repository

import (
	"context"
	"github.com/Dubjay18/ecom-api/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	// GetByID returns a user by ID
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	// GetByEmail returns a user by email
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	// Update updates a user
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := r.DB.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	user := &domain.User{}
	err := r.DB.WithContext(ctx).First(user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := r.DB.WithContext(ctx).Model(user).Updates(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
