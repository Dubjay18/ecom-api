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

// Place an order for one or more products (authenticated users)
func (s *OrderService) PlaceOrder(ctx context.Context, userID uint, req *domain.CreateOrderRequest) (*domain.Order, *common.AppError) {
	tx := s.orderRepo.BeginTx(ctx)
	defer tx.Rollback()

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

	order := &domain.Order{
		UserID:      userID,
		Status:      domain.StatusPending,
		TotalAmount: total,
		Items:       make([]domain.OrderItem, len(req.Items)),
	}

	for i, item := range req.Items {
		p, _ := s.productRepo.GetByID(ctx, item.ProductID)
		order.Items[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     p.Price,
		}
		p.Stock -= item.Quantity
		if err := s.productRepo.Update(ctx, p); err != nil {
			return nil, common.NewAppError(err, "Failed to update product", common.ErrInternalServer.Code)
		}
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, common.NewAppError(err, "Failed to create order", common.ErrInternalServer.Code)
	}
	if err := tx.Commit(); err != nil {
		return nil, common.NewAppError(err.Error, "Failed to commit transaction", common.ErrInternalServer.Code)
	}
	return order, nil
}

// List all orders for a user (authenticated)
func (s *OrderService) ListUserOrders(ctx context.Context, userID uint) ([]domain.Order, *common.AppError) {
	orders, err := s.orderRepo.List(ctx, userID)
	if err != nil {
		return nil, common.NewAppError(err, "Failed to list orders", common.ErrInternalServer.Code)
	}
	return orders, nil
}

// Cancel an order if still Pending (authenticated)
func (s *OrderService) CancelOrder(ctx context.Context, id uint) *common.AppError {
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
	for _, item := range order.Items {
		p, _ := s.productRepo.GetByID(ctx, item.ProductID)
		p.Stock += item.Quantity
		_ = s.productRepo.Update(ctx, p)
	}
	return nil
}

// Update the status of an order (admin privilege)
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id uint, newStatus domain.OrderStatus) *common.AppError {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return common.NewAppError(err, "Failed to get order", common.ErrInternalServer.Code)
	}
	order.Status = newStatus
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return common.NewAppError(err, "Failed to update order", common.ErrInternalServer.Code)
	}
	return nil
}
