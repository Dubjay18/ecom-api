package container

import (
	"github.com/Dubjay18/ecom-api/internal/config"
	"github.com/Dubjay18/ecom-api/internal/infrastructure/database"
	"github.com/Dubjay18/ecom-api/internal/repository"
	"github.com/Dubjay18/ecom-api/internal/service"
)

type Container struct {
	Config *config.Config
	DB     *database.Database

	// Repositories
	UserRepository    repository.UserRepository
	ProductRepository *repository.ProductRepository
	OrderRepository   *repository.OrderRepository

	// Services
	UserService    service.UserService
	ProductService *service.ProductService
	OrderService   *service.OrderService
}

func NewContainer(cfg *config.Config) (*Container, error) {
	// Initialize database
	db, err := database.NewPostgresDB(&cfg.DB)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	productRepo := repository.NewProductRepository(db.DB)
	orderRepo := repository.NewOrderRepository(db.DB)

	// Initialize services
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo)

	return &Container{
		Config: cfg,
		DB:     db,

		// Repositories
		UserRepository:    userRepo,
		ProductRepository: productRepo,
		OrderRepository:   orderRepo,

		// Services
		UserService:    userService,
		ProductService: productService,
		OrderService:   orderService,
	}, nil
}

func (c *Container) Close() error {
	sqlDB, err := c.DB.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
