package service

import "github.com/Dubjay18/ecom-api/internal/repository"

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}
