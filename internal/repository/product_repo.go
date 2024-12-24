package repository

import (
	"context"
	"github.com/Dubjay18/ecom-api/internal/domain"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uint) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error)
}

type productRepository struct {
	DB *gorm.DB
}

func (p *productRepository) Create(ctx context.Context, product *domain.Product) error {
	return p.DB.WithContext(ctx).Create(product).Error
}

func (p *productRepository) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	product := &domain.Product{}
	err := p.DB.WithContext(ctx).First(product, id).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *productRepository) Update(ctx context.Context, product *domain.Product) error {
	return p.DB.WithContext(ctx).Model(product).Updates(product).Error
}

func (p *productRepository) Delete(ctx context.Context, id uint) error {
	return p.DB.WithContext(ctx).Delete(&domain.Product{}, id).Error
}

func (p *productRepository) List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error) {
	var products []domain.Product
	query := p.DB.WithContext(ctx)

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	err := query.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{DB: db}
}
