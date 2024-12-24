package docs

import (
	"github.com/Dubjay18/ecom-api/internal/domain"
	_ "github.com/swaggo/swag"
)

// @title           E-commerce API
// @version         1.0
// @description     A RESTful API for an e-commerce application

// @contact.name   Your Name
// @contact.email  your.email@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// swagger definitions for User endpoints
// @Summary Register a new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterUserRequest true "User registration details"
// @Success 201 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
func registerUserDoc() {}

// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func loginDoc() {}

// swagger definitions for Product endpoints
// @Summary Create a new product
// @Description Create a new product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param product body domain.Product true "Product details"
// @Success 201 {object} domain.Product
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /products [post]
func createProductDoc() {}

// swagger definitions for Order endpoints
// @Summary Place a new order
// @Description Create a new order for authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security Bearer
// @Param order body CreateOrderRequest true "Order details"
// @Success 201 {object} domain.Order
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /orders [post]
func createOrderDoc() {}

// Request/Response models for Swagger
type RegisterUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateOrderRequest struct {
	Items []OrderItem `json:"items" binding:"required,min=1"`
}

type OrderItem struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}
