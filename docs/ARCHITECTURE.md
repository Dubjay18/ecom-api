# 🏗️ Architecture Documentation

This document describes the architecture and design patterns used in the E-commerce API project.

## 📋 Table of Contents

- [Overview](#overview)
- [Clean Architecture](#clean-architecture)
- [Project Structure](#project-structure)
- [Layers Explanation](#layers-explanation)
- [Data Flow](#data-flow)
- [Design Patterns](#design-patterns)
- [Dependencies](#dependencies)
- [Security Architecture](#security-architecture)

## Overview

The E-commerce API follows **Clean Architecture** principles, ensuring separation of concerns, testability, and maintainability. The architecture is designed to be:

- **Independent of frameworks**: Business logic doesn't depend on external frameworks
- **Testable**: Business rules can be tested without UI, database, or external services
- **Independent of UI**: The UI can be changed without affecting business rules
- **Independent of Database**: Business rules are not bound to the database
- **Independent of external services**: Business rules don't know about external services

## Clean Architecture

```
┌─────────────────────────────────────────────┐
│                   cmd/api                   │  ← Application Entry Point
├─────────────────────────────────────────────┤
│                 Handlers                    │  ← Interface Adapters
│            (Controllers/HTTP)               │
├─────────────────────────────────────────────┤
│               Middleware                    │  ← Cross-cutting Concerns
├─────────────────────────────────────────────┤
│                Services                     │  ← Business Logic
│            (Use Cases)                      │
├─────────────────────────────────────────────┤
│                Repositories                 │  ← Data Access Interface
├─────────────────────────────────────────────┤
│                 Domain                      │  ← Entities & Business Rules
├─────────────────────────────────────────────┤
│              Infrastructure                 │  ← External Concerns
│           (Database, External APIs)         │
└─────────────────────────────────────────────┘
```

## Project Structure

```
ecom-api/
├── cmd/api/                     # Application entry point
│   └── main.go                  # Server initialization and routing
│
├── internal/                    # Private application code
│   ├── config/                  # Configuration management
│   │   ├── config.go           # Configuration structures and loading
│   │   └── logger.go           # Logging configuration
│   │
│   ├── container/              # Dependency injection
│   │   └── container.go        # DI container setup
│   │
│   ├── domain/                 # Business entities and rules
│   │   ├── base.go             # Base entity with common fields
│   │   ├── user_domain.go      # User entity and business rules
│   │   ├── product_domain.go   # Product entity and business rules
│   │   └── order_domain.go     # Order entity and business rules
│   │
│   ├── handler/                # HTTP handlers (controllers)
│   │   ├── user_handler.go     # User-related endpoints
│   │   ├── product_handler.go  # Product-related endpoints
│   │   └── order_handler.go    # Order-related endpoints
│   │
│   ├── infrastructure/         # External concerns
│   │   └── database/          # Database implementation
│   │       └── db.go          # Database connection and setup
│   │
│   ├── middleware/             # HTTP middleware
│   │   ├── auth.go            # Authentication middleware
│   │   ├── admin.go           # Admin authorization middleware
│   │   └── logger.go          # Request logging middleware
│   │
│   ├── repository/             # Data access layer
│   │   ├── user_repo.go       # User data access
│   │   ├── product_repo.go    # Product data access
│   │   └── order_repo.go      # Order data access
│   │
│   ├── service/                # Business logic layer
│   │   ├── user_service.go    # User business logic
│   │   ├── product_service.go # Product business logic
│   │   └── order_service.go   # Order business logic
│   │
│   └── util/                   # Utility functions
│       └── password.go        # Password hashing utilities
│
├── pkg/                        # Public packages
│   ├── common/                 # Shared utilities
│   │   ├── errors.go          # Error handling
│   │   └── response/          # HTTP response utilities
│   │       └── response.go
│   │
│   ├── jwt/                    # JWT token management
│   │   └── jwt.go
│   │
│   └── upload/                 # File upload utilities
│       └── cloudinary.go
│
├── migrations/                 # Database migrations
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_products.up.sql
│   ├── 000002_create_products.down.sql
│   ├── 000003_create_orders.up.sql
│   └── 000003_create_orders.down.sql
│
└── docs/                       # Documentation
    ├── swagger.go              # Swagger definitions
    ├── swagger.json            # Generated OpenAPI spec
    └── swagger.yaml            # Generated OpenAPI spec
```

## Layers Explanation

### 1. Domain Layer (`internal/domain/`)

**Purpose**: Contains business entities and core business rules.

**Responsibilities**:
- Define business entities (User, Product, Order)
- Business rules and validations
- Domain-specific types and constants
- Request/Response models

**Key Files**:
- `base.go`: Common entity fields (ID, timestamps)
- `user_domain.go`: User entity and authentication models
- `product_domain.go`: Product entity and filtering models
- `order_domain.go`: Order entity and status management

**Example Entity**:
```go
type User struct {
    Base
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Password  string    `json:"-" gorm:"not null"`
    FirstName string    `json:"first_name" gorm:"size:100"`
    LastName  string    `json:"last_name" gorm:"size:100"`
    Role      UserRole  `json:"role" gorm:"type:varchar(20);default:'user'"`
    Orders    []Order   `json:"orders,omitempty" gorm:"foreignKey:UserID"`
    Addresses []Address `json:"addresses,omitempty" gorm:"foreignKey:UserID"`
}
```

### 2. Repository Layer (`internal/repository/`)

**Purpose**: Defines interfaces for data access and provides implementations.

**Responsibilities**:
- Data persistence operations (CRUD)
- Database query implementations
- Data mapping between domain entities and database models

**Interface Example**:
```go
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    GetByID(ctx context.Context, id uint) (*domain.User, error)
    GetByEmail(ctx context.Context, email string) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uint) error
}
```

### 3. Service Layer (`internal/service/`)

**Purpose**: Contains business logic and use cases.

**Responsibilities**:
- Implement business use cases
- Coordinate between repositories
- Business rule enforcement
- Transaction management

**Example Service**:
```go
type UserService interface {
    Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, *common.AppError)
    Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, *common.AppError)
    GetByID(ctx context.Context, id uint) (*domain.User, *common.AppError)
}
```

### 4. Handler Layer (`internal/handler/`)

**Purpose**: HTTP request/response handling and routing.

**Responsibilities**:
- HTTP request parsing
- Input validation
- Response formatting
- Route definition
- Middleware application

**Example Handler**:
```go
func (h *UserHandler) Register(c *gin.Context) {
    var req domain.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
        return
    }
    
    user, appErr := h.s.Register(c.Request.Context(), req)
    if appErr != nil {
        response.ErrorResponse(c, appErr.Code, appErr.Message, appErr)
        return
    }
    
    response.SuccessResponse(c, http.StatusCreated, "User registered successfully", user)
}
```

### 5. Infrastructure Layer (`internal/infrastructure/`)

**Purpose**: External services and frameworks integration.

**Responsibilities**:
- Database connections
- External API integrations
- File storage services
- Message queues (if applicable)

### 6. Middleware Layer (`internal/middleware/`)

**Purpose**: Cross-cutting concerns and request/response processing.

**Responsibilities**:
- Authentication and authorization
- Request logging
- CORS handling
- Rate limiting (if implemented)
- Error handling

## Data Flow

### Request Flow (User Registration)

```
1. HTTP Request
   ↓
2. Gin Router
   ↓
3. Middleware Chain
   ├── CORS
   ├── Logging
   └── Recovery
   ↓
4. Handler (UserHandler.Register)
   ├── Parse Request
   ├── Validate Input
   └── Call Service
   ↓
5. Service (UserService.Register)
   ├── Business Logic
   ├── Password Hashing
   └── Call Repository
   ↓
6. Repository (UserRepository.Create)
   ├── Database Query
   └── Return Result
   ↓
7. Response Chain (reverse order)
```

### Authenticated Request Flow

```
1. HTTP Request + JWT Token
   ↓
2. Authentication Middleware
   ├── Extract Token
   ├── Validate Token
   ├── Extract User Info
   └── Set Context
   ↓
3. Authorization Middleware (if admin required)
   ├── Check User Role
   └── Allow/Deny Access
   ↓
4. Handler → Service → Repository
   ↓
5. Response
```

## Design Patterns

### 1. Dependency Injection

The application uses a DI container to manage dependencies:

```go
type Container struct {
    Config *config.Config
    DB     *database.Database
    
    // Repositories
    UserRepository    repository.UserRepository
    ProductRepository repository.ProductRepository
    OrderRepository   repository.OrderRepository
    
    // Services
    UserService    service.UserService
    ProductService *service.ProductService
    OrderService   *service.OrderService
}
```

### 2. Repository Pattern

Abstracts data access logic:

```go
type ProductRepository interface {
    Create(ctx context.Context, product *domain.Product) error
    GetByID(ctx context.Context, id uint) (*domain.Product, error)
    List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error)
    Update(ctx context.Context, product *domain.Product) error
    Delete(ctx context.Context, id uint) error
}
```

### 3. Service Pattern

Encapsulates business logic:

```go
type OrderService struct {
    orderRepo   repository.OrderRepository
    productRepo repository.ProductRepository
}

func (s *OrderService) PlaceOrder(ctx context.Context, userID uint, req *domain.CreateOrderRequest) (*domain.Order, *common.AppError) {
    // Business logic for order placement
    // Validate products, calculate totals, manage inventory
}
```

### 4. Middleware Pattern

Cross-cutting concerns:

```go
func AuthMiddleware(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract and validate JWT token
        // Set user context
        c.Next()
    }
}
```

## Dependencies

### Core Dependencies

- **Gin**: HTTP web framework
- **GORM**: ORM for database operations
- **JWT**: Authentication token management
- **Viper**: Configuration management
- **Logrus**: Structured logging

### Database

- **PostgreSQL**: Primary database
- **Migrations**: Database versioning with golang-migrate

### External Services

- **Cloudinary**: Image storage and processing
- **Stripe**: Payment processing (configurable)

## Security Architecture

### Authentication Flow

```
1. User Registration/Login
   ↓
2. Password Hashing (bcrypt)
   ↓
3. JWT Token Generation
   ├── User ID
   ├── Email
   ├── Role (user/admin)
   └── Expiration
   ↓
4. Token Storage (client-side)
   ↓
5. Subsequent Requests
   ├── Token in Authorization header
   ├── Token validation
   └── User context extraction
```

### Authorization Levels

1. **Public**: No authentication required
   - User registration
   - User login
   - Health check

2. **Authenticated**: Valid JWT token required
   - User profile
   - Order placement
   - Order listing

3. **Admin**: Admin role required
   - Product management
   - Order status updates
   - User management

### Security Measures

- **Password Hashing**: bcrypt with salt
- **JWT Tokens**: Stateless authentication
- **CORS**: Configurable cross-origin requests
- **Input Validation**: Request validation with gin binding
- **SQL Injection Prevention**: GORM ORM usage
- **Rate Limiting**: Can be implemented with middleware

## Testing Strategy

### Unit Testing

- **Service Layer**: Business logic testing
- **Repository Layer**: Data access testing
- **Handler Layer**: HTTP endpoint testing

### Integration Testing

- **Database**: Integration with PostgreSQL
- **API**: End-to-end API testing
- **External Services**: Cloudinary integration testing

### Test Structure

```
tests/
├── unit/
│   ├── service/
│   ├── repository/
│   └── handler/
├── integration/
│   ├── api/
│   └── database/
└── fixtures/
    └── test_data.sql
```

## Performance Considerations

### Database

- **Connection Pooling**: Configured max connections
- **Indexes**: Proper indexing on frequently queried fields
- **Migrations**: Versioned schema changes

### API

- **Pagination**: For large data sets
- **Caching**: Can be implemented with Redis
- **Compression**: Gzip middleware available

### Monitoring

- **Logging**: Structured logging with Logrus
- **Health Checks**: `/health` endpoint
- **Metrics**: Can be added with Prometheus

---

This architecture ensures the application is maintainable, testable, and scalable while following industry best practices and clean code principles.