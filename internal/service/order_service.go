package service

import (
	"context"
	"net/http"

	"github.com/Dubjay18/ecom-api/internal/domain"
	"github.com/Dubjay18/ecom-api/internal/repository"
	"github.com/Dubjay18/ecom-api/pkg/common"
)

type OrderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewOrderService(or repository.OrderRepository, pr repository.ProductRepository) *OrderService {
	return &OrderService{orderRepo: or, productRepo: pr}
}

// Create creates a new order
func (s *OrderService) Create(ctx context.Context, userID uint, req *domain.CreateOrderRequest) (*domain.Order, *common.AppError) {
	// Start transaction
	tx := s.orderRepo.BeginTx(ctx)
	defer tx.Rollback()

	// Calculate total and check stock
	var total float64
	for _, item := range req.Items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, common.NewAppError(err, "Failed to get product", common.ErrInternalServer.Code)
		}

		if product.Stock < item.Quantity {
			return nil, common.NewAppError(nil, "Insufficient stock", http.StatusBadRequest)
		}

		total += product.Price * float64(item.Quantity)
	}

	// Create order
	order := &domain.Order{
		UserID:      userID,
		Status:      domain.StatusPending,
		TotalAmount: total,
		Items:       make([]domain.OrderItem, len(req.Items)),
	}

	// Create order items and update stock
	for i, item := range req.Items {
		product, _ := s.productRepo.GetByID(ctx, item.ProductID)
		order.Items[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		product.Stock -= item.Quantity
		if err := s.productRepo.Update(ctx, product); err != nil {
			return nil, common.NewAppError(err, "Failed to update product", common.ErrInternalServer.Code)
		}
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, common.NewAppError(err, "Failed to create order", common.ErrInternalServer.Code)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, common.NewAppError(err.Error, "Failed to commit transaction", common.ErrInternalServer.Code)
	}

	return order, nil
}

// GetByID returns an order by ID
func (s *OrderService) GetByID(ctx context.Context, id uint) (*domain.Order, *common.AppError) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, common.NewAppError(err, "Failed to get order", common.ErrInternalServer.Code)
	}
	return order, nil
}

// Update updates an order
func (s *OrderService) Update(ctx context.Context, order *domain.Order) *common.AppError {
	err := s.orderRepo.Update(ctx, order)
	if err != nil {
		return common.NewAppError(err, "Failed to update order", common.ErrInternalServer.Code)
	}
	return nil
}

// List returns list of orders
func (s *OrderService) List(ctx context.Context, userID uint) ([]domain.Order, *common.AppError) {
	orders, err := s.orderRepo.List(ctx, userID)
	if err != nil {
		return nil, common.NewAppError(err, "Failed to list orders", common.ErrInternalServer.Code)
	}
	return orders, nil
}

// Cancel cancels an order
func (s *OrderService) Cancel(ctx context.Context, id uint) *common.AppError {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return common.NewAppError(err, "Failed to get order", common.ErrInternalServer.Code)
	}

	if order.Status != domain.StatusPending {
		return common.NewAppError(nil, "Order cannot be cancelled", http.StatusBadRequest)
	}

	order.Status = domain.StatusCancelled
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return common.NewAppError(err, "Failed to update order", common.ErrInternalServer.Code)
	}

	// Refund stock
	for _, item := range order.Items {
		product, _ := s.productRepo.GetByID(ctx, item.ProductID)
		product.Stock += item.Quantity
		if err := s.productRepo.Update(ctx, product); err != nil {
			return common.NewAppError(err, "Failed to update product", common.ErrInternalServer.Code)
		}
	}

	return nil
}

// Complete completes an order
func (s *OrderService) Complete(ctx context.Context, id uint) *common.AppError {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return common.NewAppError(err, "Failed to get order", common.ErrInternalServer.Code)
	}

	if order.Status != domain.StatusPending {
		return common.NewAppError(nil, "Order cannot be completed", http.StatusBadRequest)
	}

	order.Status = domain.StatusConfirmed
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return common.NewAppError(err, "Failed to update order", common.ErrInternalServer.Code)
	}

	return nil
}
