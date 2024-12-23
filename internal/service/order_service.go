package service

import "github.com/Dubjay18/ecom-api/internal/repository"

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}
