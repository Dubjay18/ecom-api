package repository

import (
	"context"

	"github.com/Dubjay18/ecom-api/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	BeginTx(ctx context.Context) *gorm.DB
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id uint) (*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	List(ctx context.Context, userID uint) ([]domain.Order, error)
}

type orderRepository struct {
	DB *gorm.DB
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	return r.DB.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) GetByID(ctx context.Context, id uint) (*domain.Order, error) {
	order := &domain.Order{}
	err := r.DB.WithContext(ctx).First(order, id).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *orderRepository) Update(ctx context.Context, order *domain.Order) error {
	return r.DB.WithContext(ctx).Model(order).Updates(order).Error
}

func (r *orderRepository) List(ctx context.Context, userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{DB: db}
}

func (r *orderRepository) BeginTx(ctx context.Context) *gorm.DB {
	return r.DB.WithContext(ctx).Begin()
}
