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

	// Fetch all product details in a single query
	productIDs := make([]uint, len(req.Items))
	for i, item := range req.Items {
		productIDs[i] = item.ProductID
	}

	products, err := s.productRepo.GetByIDs(ctx, productIDs)
	if err != nil {
		return nil, common.NewAppError(err, "Failed to fetch products", common.ErrInternalServer.Code)
	}

	// Map products for quick lookup
	productMap := make(map[uint]*domain.Product)
	for _, product := range products {
		productMap[product.ID] = &product
	}

	// Validate stock and calculate total
	var total float64
	orderItems := make([]domain.OrderItem, len(req.Items))

	for i, item := range req.Items {
		product, exists := productMap[item.ProductID]
		if !exists {
			return nil, common.NewAppError(nil, "Product not found", http.StatusBadRequest)
		}
		if product.Stock < item.Quantity {
			return nil, common.NewAppError(nil, "Insufficient stock for product", http.StatusBadRequest)
		}

		orderItems[i] = domain.OrderItem{
			Product:  *product,
			Quantity: item.Quantity,
			Price:    product.Price * float64(item.Quantity),
		}
		total += orderItems[i].Price
	}

	// Create shipping address
	shippingAddr := &domain.Address{
		UserID:     userID,
		Street:     req.ShippingAddr.Street,
		City:       req.ShippingAddr.City,
		State:      req.ShippingAddr.State,
		Country:    req.ShippingAddr.Country,
		PostalCode: req.ShippingAddr.PostalCode,
	}
	if err := s.orderRepo.CreatAddress(ctx, shippingAddr); err != nil {
		return nil, common.NewAppError(err, "Failed to create shipping address", common.ErrInternalServer.Code)
	}

	// Update product stocks
	for _, item := range req.Items {
		product := productMap[item.ProductID]
		product.Stock -= item.Quantity
		if err := s.productRepo.Update(ctx, product); err != nil {
			return nil, common.NewAppError(err, "Failed to update product stock", common.ErrInternalServer.Code)
		}
	}

	// Create order
	order := &domain.Order{
		UserID:            userID,
		Status:            domain.StatusPending,
		TotalAmount:       total,
		Items:             orderItems,
		ShippingAddressID: shippingAddr.ID,
	}
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, common.NewAppError(err, "Failed to create order", common.ErrInternalServer.Code)
	}

	// Commit transaction
	if err := tx.Commit(); err.Error != nil {
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
	if newStatus != domain.StatusPending && newStatus != domain.StatusShipped && newStatus != domain.StatusDelivered {
		return common.NewAppError(nil, "Invalid status", http.StatusBadRequest)
	}
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
