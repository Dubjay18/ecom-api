package service

import (
	"context"

	"github.com/Dubjay18/ecom-api/internal/domain"
	"github.com/Dubjay18/ecom-api/internal/repository"
	"github.com/Dubjay18/ecom-api/pkg/common"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// Create creates a new product
func (s *ProductService) Create(ctx context.Context, product *domain.Product) *common.AppError {

	err := s.repo.Create(ctx, product)
	if err != nil {
		return common.NewAppError(err, "Failed to create product", common.ErrInternalServer.Code)
	}
	return nil
}

// GetByID returns a product by ID
func (s *ProductService) GetByID(ctx context.Context, id uint) (*domain.Product, *common.AppError) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, common.NewAppError(err, "Failed to get product", common.ErrInternalServer.Code)
	}
	return product, nil
}

// Update updates a product
func (s *ProductService) Update(ctx context.Context, product *domain.Product) *common.AppError {
	err := s.repo.Update(ctx, product)
	if err != nil {
		return common.NewAppError(err, "Failed to update product", common.ErrInternalServer.Code)
	}
	return nil
}

// Delete deletes a product
func (s *ProductService) Delete(ctx context.Context, id uint) *common.AppError {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return common.NewAppError(err, "Failed to delete product", common.ErrInternalServer.Code)
	}
	return nil
}

// List returns list of products
func (s *ProductService) List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, *common.AppError) {
	products, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, common.NewAppError(err, "Failed to list products", common.ErrInternalServer.Code)
	}
	return products, nil
}
