# 🛠️ Development Guide

This guide provides detailed information for developers working on the E-commerce API project.

## 📋 Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing Strategy](#testing-strategy)
- [Debugging](#debugging)
- [Performance Optimization](#performance-optimization)
- [Security Considerations](#security-considerations)
- [Common Tasks](#common-tasks)

## Development Environment Setup

### Prerequisites

Ensure you have the following installed:

- **Go 1.23.3+**
- **PostgreSQL 12+**
- **Git**
- **Make**
- **Docker & Docker Compose** (optional but recommended)

### Quick Setup

```bash
# Clone the repository
git clone https://github.com/Dubjay18/ecom-api.git
cd ecom-api

# Install dependencies
go mod tidy

# Setup environment
cp .env.example .env
# Edit .env with your configuration

# Install development tools
make install-tools

# Setup database
make db-setup

# Run migrations
make migrate-up

# Start development server
make dev
```

### Development Tools Installation

```bash
# Install required tools
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/cosmtrek/air@latest
go install github.com/golang/mock/mockgen@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### IDE Configuration

#### VS Code

Recommended extensions:
- Go (official)
- REST Client
- Docker
- GitLens
- Thunder Client

**VS Code Settings** (`.vscode/settings.json`):
```json
{
    "go.testFlags": ["-v"],
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "package",
    "go.formatTool": "goimports",
    "go.useLanguageServer": true,
    "go.docsTool": "godoc",
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.organizeImports": true
    }
}
```

#### GoLand/IntelliJ

1. Enable Go modules support
2. Configure code style according to `gofmt`
3. Set up run configurations for tests and main application

## Project Structure

### Directory Layout

```
ecom-api/
├── cmd/api/                    # Application entry point
│   └── main.go                # Server initialization
├── internal/                  # Private application code
│   ├── config/               # Configuration
│   ├── container/            # Dependency injection
│   ├── domain/               # Business entities
│   ├── handler/              # HTTP handlers
│   ├── infrastructure/       # External services
│   ├── middleware/           # HTTP middleware
│   ├── repository/           # Data access layer
│   ├── service/              # Business logic
│   └── util/                 # Utilities
├── pkg/                      # Public packages
│   ├── common/               # Shared utilities
│   ├── jwt/                  # JWT management
│   └── upload/               # File uploads
├── migrations/               # Database migrations
├── docs/                     # Documentation
├── tests/                    # Test files
├── scripts/                  # Build and deployment scripts
└── deployments/              # Deployment configurations
```

### Adding New Features

When adding a new feature, follow this structure:

1. **Domain Layer**: Define entities and business rules
2. **Repository Layer**: Define data access interface
3. **Service Layer**: Implement business logic
4. **Handler Layer**: Implement HTTP endpoints
5. **Middleware**: Add any required middleware
6. **Tests**: Add comprehensive tests
7. **Documentation**: Update API documentation

### Example: Adding a Category Feature

```go
// 1. Domain (internal/domain/category_domain.go)
type Category struct {
    Base
    Name        string `json:"name" gorm:"uniqueIndex;size:100;not null"`
    Description string `json:"description" gorm:"type:text"`
    ParentID    *uint  `json:"parent_id" gorm:"index"`
    Parent      *Category `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
    Children    []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// 2. Repository Interface (internal/repository/category_repo.go)
type CategoryRepository interface {
    Create(ctx context.Context, category *domain.Category) error
    GetByID(ctx context.Context, id uint) (*domain.Category, error)
    List(ctx context.Context) ([]domain.Category, error)
    Update(ctx context.Context, category *domain.Category) error
    Delete(ctx context.Context, id uint) error
}

// 3. Service (internal/service/category_service.go)
type CategoryService interface {
    Create(ctx context.Context, req domain.CreateCategoryRequest) (*domain.Category, *common.AppError)
    GetByID(ctx context.Context, id uint) (*domain.Category, *common.AppError)
    List(ctx context.Context) ([]domain.Category, *common.AppError)
}

// 4. Handler (internal/handler/category_handler.go)
type CategoryHandler struct {
    service CategoryService
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
    // Implementation
}
```

## Development Workflow

### Daily Development

1. **Start development environment:**
   ```bash
   make dev-up      # Start database and dependencies
   make dev         # Start API with live reload
   ```

2. **Make changes and test:**
   ```bash
   make test        # Run tests
   make lint        # Check code quality
   make swagger     # Update documentation
   ```

3. **Before committing:**
   ```bash
   make pre-commit  # Run all checks
   ```

### Branch Management

```bash
# Create feature branch
git checkout -b feature/category-management

# Regular commits with conventional commit messages
git commit -m "feat: add category entity and domain logic"
git commit -m "feat: implement category repository"
git commit -m "feat: add category service with business logic"
git commit -m "feat: implement category HTTP handlers"
git commit -m "test: add comprehensive category tests"
git commit -m "docs: update API documentation for categories"

# Before pushing
make pre-commit
git push origin feature/category-management
```

### Code Review Checklist

- [ ] Code follows project conventions
- [ ] All tests pass
- [ ] Code coverage meets requirements
- [ ] Documentation is updated
- [ ] No security vulnerabilities
- [ ] Performance impact considered
- [ ] Error handling is appropriate
- [ ] Logging is adequate

## Testing Strategy

### Test Pyramid

```
        ┌─────────────────┐
        │   E2E Tests     │  ← Few, high-level tests
        │                 │
        ├─────────────────┤
        │ Integration     │  ← Medium number of tests
        │    Tests        │
        ├─────────────────┤
        │   Unit Tests    │  ← Many, fast tests
        │                 │
        └─────────────────┘
```

### Unit Tests

```go
// Example: Service layer unit test
func TestUserService_Register(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)
    mockJWT := mocks.NewMockJWTService(ctrl)
    
    service := NewUserService(mockRepo, mockJWT)

    tests := []struct {
        name          string
        request       domain.RegisterRequest
        setupMocks    func()
        expectedError bool
    }{
        {
            name: "successful registration",
            request: domain.RegisterRequest{
                Email:     "test@example.com",
                Password:  "password123",
                FirstName: "John",
                LastName:  "Doe",
            },
            setupMocks: func() {
                mockRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").
                    Return(nil, errors.New("not found"))
                mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
                    Return(nil)
            },
            expectedError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setupMocks()
            
            user, err := service.Register(context.Background(), tt.request)
            
            if tt.expectedError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, user)
            }
        })
    }
}
```

### Integration Tests

```go
// Example: Repository integration test
func TestUserRepository_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    repo := NewUserRepository(db)

    t.Run("Create and Get User", func(t *testing.T) {
        user := &domain.User{
            Email:     "test@example.com",
            Password:  "hashedpassword",
            FirstName: "John",
            LastName:  "Doe",
        }

        // Create user
        err := repo.Create(context.Background(), user)
        assert.NoError(t, err)
        assert.NotZero(t, user.ID)

        // Get user
        retrieved, err := repo.GetByID(context.Background(), user.ID)
        assert.NoError(t, err)
        assert.Equal(t, user.Email, retrieved.Email)
    })
}
```

### End-to-End Tests

```go
// Example: HTTP endpoint E2E test
func TestUserEndpoints_E2E(t *testing.T) {
    // Setup test server
    router := setupTestRouter(t)
    
    t.Run("User Registration Flow", func(t *testing.T) {
        // Register user
        reqBody := `{
            "email": "test@example.com",
            "password": "password123",
            "first_name": "John",
            "last_name": "Doe"
        }`
        
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusCreated, w.Code)
        
        var response map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &response)
        
        assert.Equal(t, "User registered successfully", response["message"])
    })
}
```

### Test Database Setup

```go
func setupTestDB(t *testing.T) *gorm.DB {
    dsn := "host=localhost user=test password=test dbname=test_db port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    require.NoError(t, err)
    
    // Auto migrate
    err = db.AutoMigrate(&domain.User{}, &domain.Product{}, &domain.Order{})
    require.NoError(t, err)
    
    // Clean up after test
    t.Cleanup(func() {
        db.Exec("TRUNCATE users, products, orders CASCADE")
    })
    
    return db
}
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./internal/service/...

# Run tests with race detection
go test -race ./...

# Run integration tests only
go test -tags=integration ./...

# Run tests in verbose mode
go test -v -run TestUserService ./internal/service/
```

## Debugging

### Debugging Tools

1. **Delve Debugger:**
   ```bash
   # Install delve
   go install github.com/go-delve/delve/cmd/dlv@latest
   
   # Debug application
   dlv debug cmd/api/main.go
   
   # Debug with arguments
   dlv debug cmd/api/main.go -- --config=.env
   
   # Debug tests
   dlv test ./internal/service/
   ```

2. **VS Code Debugging:**
   ```json
   // .vscode/launch.json
   {
       "version": "0.2.0",
       "configurations": [
           {
               "name": "Launch API",
               "type": "go",
               "request": "launch",
               "mode": "debug",
               "program": "${workspaceFolder}/cmd/api/main.go",
               "env": {
                   "GIN_MODE": "debug"
               },
               "args": []
           }
       ]
   }
   ```

### Logging and Monitoring

```go
// Structured logging example
func (s *userService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, *common.AppError) {
    logger := s.logger.WithFields(logrus.Fields{
        "operation": "user_registration",
        "email":     req.Email,
    })
    
    logger.Info("Starting user registration")
    
    // Check if user exists
    existingUser, err := s.repo.GetByEmail(ctx, req.Email)
    if err == nil && existingUser != nil {
        logger.Warn("User registration failed: email already exists")
        return nil, &common.AppError{
            Code:    http.StatusConflict,
            Message: "User with this email already exists",
        }
    }
    
    // Create user
    user, err := s.createUser(ctx, req)
    if err != nil {
        logger.WithError(err).Error("Failed to create user")
        return nil, &common.AppError{
            Code:    http.StatusInternalServerError,
            Message: "Failed to create user",
        }
    }
    
    logger.WithField("user_id", user.ID).Info("User registered successfully")
    return user, nil
}
```

### Performance Profiling

```go
// Add profiling endpoints (development only)
import _ "net/http/pprof"

if gin.Mode() == gin.DebugMode {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```

```bash
# CPU profiling
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory profiling
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine profiling
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

## Performance Optimization

### Database Optimization

1. **Query Optimization:**
   ```go
   // Use proper preloading
   var users []domain.User
   db.Preload("Orders").Find(&users)
   
   // Use selective loading
   db.Select("id", "email", "first_name").Find(&users)
   
   // Use pagination
   db.Offset(offset).Limit(limit).Find(&users)
   ```

2. **Indexing:**
   ```sql
   -- Add indexes for frequently queried fields
   CREATE INDEX idx_users_email ON users(email);
   CREATE INDEX idx_products_sku ON products(sku);
   CREATE INDEX idx_orders_user_id ON orders(user_id);
   ```

### API Performance

1. **Pagination:**
   ```go
   type PaginationRequest struct {
       Page     int `form:"page" binding:"min=1"`
       PageSize int `form:"page_size" binding:"min=1,max=100"`
   }
   
   func (r *productRepository) List(ctx context.Context, pagination PaginationRequest) ([]domain.Product, error) {
       offset := (pagination.Page - 1) * pagination.PageSize
       var products []domain.Product
       
       err := r.db.Offset(offset).Limit(pagination.PageSize).Find(&products).Error
       return products, err
   }
   ```

2. **Caching (Optional):**
   ```go
   // Redis caching example
   func (s *productService) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
       // Try cache first
       cacheKey := fmt.Sprintf("product:%d", id)
       if cachedProduct := s.cache.Get(cacheKey); cachedProduct != nil {
           return cachedProduct.(*domain.Product), nil
       }
       
       // Get from database
       product, err := s.repo.GetByID(ctx, id)
       if err != nil {
           return nil, err
       }
       
       // Store in cache
       s.cache.Set(cacheKey, product, 5*time.Minute)
       return product, nil
   }
   ```

## Security Considerations

### Input Validation

```go
// Use binding tags for validation
type CreateProductRequest struct {
    Name        string  `json:"name" binding:"required,min=1,max=255"`
    Price       float64 `json:"price" binding:"required,gt=0"`
    Description string  `json:"description" binding:"max=1000"`
    Stock       int     `json:"stock" binding:"required,gte=0"`
    SKU         string  `json:"sku" binding:"required,alphanum,min=3,max=50"`
}

// Custom validation
func validateSKU(sku string) error {
    if !regexp.MustCompile(`^[A-Z0-9-]+$`).MatchString(sku) {
        return errors.New("SKU must contain only uppercase letters, numbers, and hyphens")
    }
    return nil
}
```

### Authentication & Authorization

```go
// JWT middleware with proper error handling
func AuthMiddleware(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c.GetHeader("Authorization"))
        if token == "" {
            response.ErrorResponse(c, http.StatusUnauthorized, "Missing authorization token", nil)
            c.Abort()
            return
        }
        
        claims, err := validateJWT(token, secretKey)
        if err != nil {
            response.ErrorResponse(c, http.StatusUnauthorized, "Invalid token", err)
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        c.Set("is_admin", claims.IsAdmin)
        c.Next()
    }
}
```

### SQL Injection Prevention

```go
// Always use parameterized queries (GORM handles this)
func (r *userRepository) GetByEmailAndStatus(ctx context.Context, email string, status string) (*domain.User, error) {
    var user domain.User
    
    // Good - parameterized query
    err := r.db.Where("email = ? AND status = ?", email, status).First(&user).Error
    
    // Never do this - vulnerable to SQL injection
    // query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s' AND status = '%s'", email, status)
    
    return &user, err
}
```

## Common Tasks

### Adding a New Endpoint

1. **Define domain model:**
   ```go
   // internal/domain/category_domain.go
   type Category struct {
       Base
       Name string `json:"name" binding:"required"`
   }
   ```

2. **Add repository interface:**
   ```go
   // internal/repository/category_repo.go
   type CategoryRepository interface {
       Create(ctx context.Context, category *domain.Category) error
   }
   ```

3. **Implement service:**
   ```go
   // internal/service/category_service.go
   func (s *categoryService) Create(ctx context.Context, req domain.CreateCategoryRequest) (*domain.Category, error) {
       // Business logic here
   }
   ```

4. **Add HTTP handler:**
   ```go
   // internal/handler/category_handler.go
   func (h *CategoryHandler) CreateCategory(c *gin.Context) {
       // HTTP handling here
   }
   ```

5. **Register routes:**
   ```go
   // cmd/api/main.go
   categoryHandler := handler.NewCategoryHandler(api, container.CategoryService)
   ```

6. **Add tests:**
   ```go
   // tests/
   func TestCategoryService_Create(t *testing.T) { /* tests */ }
   func TestCategoryHandler_CreateCategory(t *testing.T) { /* tests */ }
   ```

7. **Update documentation:**
   ```bash
   make swagger
   ```

### Database Migration

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq add_categories_table

# Edit migration files
# migrations/000004_add_categories_table.up.sql
# migrations/000004_add_categories_table.down.sql

# Apply migration
make migrate-up

# Rollback if needed
make migrate-down
```

### Adding Middleware

```go
// internal/middleware/rate_limit.go
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Rate limiting logic
        c.Next()
    }
}

// Register in main.go
api.Use(middleware.RateLimitMiddleware(100, time.Minute))
```

### Environment-Specific Configuration

```go
// internal/config/config.go
func LoadConfig(path string) (*Config, error) {
    viper.SetConfigName(".env")
    viper.SetConfigType("env")
    viper.AddConfigPath(path)
    
    // Allow environment variables to override
    viper.AutomaticEnv()
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

---

This development guide should help you get started with contributing to the E-commerce API project. For additional questions, please refer to the other documentation files or open an issue on GitHub.