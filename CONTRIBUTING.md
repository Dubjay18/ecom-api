# 🤝 Contributing Guidelines

Thank you for your interest in contributing to the E-commerce API project! This document provides guidelines and information for contributors.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)
- [Community](#community)

## Code of Conduct

### Our Pledge

We pledge to make participation in our project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity and expression, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Our Standards

**Examples of behavior that contributes to creating a positive environment:**

- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

**Examples of unacceptable behavior:**

- The use of sexualized language or imagery and unwelcome sexual attention or advances
- Trolling, insulting/derogatory comments, and personal or political attacks
- Public or private harassment
- Publishing others' private information without explicit permission
- Other conduct which could reasonably be considered inappropriate in a professional setting

### Enforcement

Project maintainers are responsible for clarifying the standards of acceptable behavior and are expected to take appropriate and fair corrective action in response to any instances of unacceptable behavior.

## Getting Started

### Prerequisites

Before contributing, make sure you have:

- **Go 1.23.3+** installed
- **PostgreSQL 12+** for database
- **Git** for version control
- **Make** for build automation
- **Docker** (optional, for containerized development)

### Setting Up Development Environment

1. **Fork the repository:**
   ```bash
   # Go to https://github.com/Dubjay18/ecom-api and fork the repo
   ```

2. **Clone your fork:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/ecom-api.git
   cd ecom-api
   ```

3. **Add upstream remote:**
   ```bash
   git remote add upstream https://github.com/Dubjay18/ecom-api.git
   ```

4. **Install dependencies:**
   ```bash
   go mod tidy
   ```

5. **Set up environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your local configuration
   ```

6. **Set up database:**
   ```bash
   createdb ecom_api_dev
   make migrate-up
   ```

7. **Install development tools:**
   ```bash
   # Install swag for API documentation
   go install github.com/swaggo/swag/cmd/swag@latest
   
   # Install golangci-lint for code quality
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Install air for live reloading
   go install github.com/cosmtrek/air@latest
   ```

## Development Workflow

### Branch Strategy

We use a Git flow approach:

- **main**: Production-ready code
- **develop**: Integration branch for features
- **feature/***: Feature development branches
- **hotfix/***: Critical bug fixes
- **release/***: Release preparation branches

### Creating a Feature Branch

```bash
# Update your local main branch
git checkout main
git pull upstream main

# Create and checkout feature branch
git checkout -b feature/your-feature-name

# Push feature branch to your fork
git push -u origin feature/your-feature-name
```

### Making Changes

1. **Write code following our standards**
2. **Add or update tests**
3. **Update documentation if needed**
4. **Ensure all tests pass**
5. **Run linter and fix any issues**

```bash
# Run tests
make test

# Run linter
make lint

# Build application
make build

# Generate/update documentation
make swagger
```

### Committing Changes

We follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

```bash
# Feature
git commit -m "feat: add user profile update endpoint"

# Bug fix
git commit -m "fix: resolve authentication middleware issue"

# Documentation
git commit -m "docs: update API documentation for orders"

# Refactor
git commit -m "refactor: improve error handling in user service"

# Test
git commit -m "test: add unit tests for product service"

# Chore
git commit -m "chore: update dependencies"
```

**Commit Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

## Coding Standards

### Go Style Guide

We follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html).

#### Naming Conventions

```go
// Good
type UserService interface {
    GetByID(ctx context.Context, id uint) (*domain.User, error)
}

type userService struct {
    repo repository.UserRepository
}

// Constants
const (
    DefaultPageSize = 20
    MaxPageSize     = 100
)

// Private functions use camelCase
func validateEmail(email string) bool {
    // implementation
}
```

#### Error Handling

```go
// Good - Wrap errors with context
func (s *userService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user by ID %d: %w", id, err)
    }
    return user, nil
}

// Use custom error types for business logic
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidEmail = errors.New("invalid email format")
)
```

#### Function Structure

```go
// Good - Clear function signature and documentation
// CreateUser creates a new user with the provided details.
// It validates the input, hashes the password, and stores the user in the database.
func (s *userService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Validate input
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    // Business logic
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    // Create user
    user := &User{
        Email:    req.Email,
        Password: string(hashedPassword),
        // ... other fields
    }
    
    // Store in database
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}
```

### Code Organization

#### File Structure

```go
// handler/user_handler.go
package handler

import (
    // Standard library imports first
    "net/http"
    "strconv"
    
    // Third-party imports
    "github.com/gin-gonic/gin"
    
    // Local imports last
    "github.com/Dubjay18/ecom-api/internal/domain"
    "github.com/Dubjay18/ecom-api/internal/service"
)
```

#### Interface Design

```go
// Keep interfaces small and focused
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    GetByID(ctx context.Context, id uint) (*domain.User, error)
    GetByEmail(ctx context.Context, email string) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uint) error
}

// Separate interfaces for different concerns
type UserValidator interface {
    ValidateEmail(email string) error
    ValidatePassword(password string) error
}
```

### HTTP Handler Standards

```go
// Good - Consistent error handling and response format
func (h *UserHandler) GetUser(c *gin.Context) {
    // Parse and validate input
    userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        response.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
        return
    }
    
    // Call service
    user, appErr := h.service.GetByID(c.Request.Context(), uint(userID))
    if appErr != nil {
        response.ErrorResponse(c, appErr.Code, appErr.Message, appErr)
        return
    }
    
    // Return success response
    response.SuccessResponse(c, http.StatusOK, "User retrieved successfully", user)
}
```

## Testing Guidelines

### Test Structure

We follow the **Arrange-Act-Assert** pattern:

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name          string
        request       domain.CreateUserRequest
        mockSetup     func(*mocks.MockUserRepository)
        expectedError string
    }{
        {
            name: "successful user creation",
            request: domain.CreateUserRequest{
                Email:     "test@example.com",
                Password:  "password123",
                FirstName: "John",
                LastName:  "Doe",
            },
            mockSetup: func(mockRepo *mocks.MockUserRepository) {
                mockRepo.EXPECT().
                    GetByEmail(gomock.Any(), "test@example.com").
                    Return(nil, errors.New("not found"))
                mockRepo.EXPECT().
                    Create(gomock.Any(), gomock.Any()).
                    Return(nil)
            },
            expectedError: "",
        },
        {
            name: "user already exists",
            request: domain.CreateUserRequest{
                Email:     "existing@example.com",
                Password:  "password123",
                FirstName: "Jane",
                LastName:  "Doe",
            },
            mockSetup: func(mockRepo *mocks.MockUserRepository) {
                mockRepo.EXPECT().
                    GetByEmail(gomock.Any(), "existing@example.com").
                    Return(&domain.User{}, nil)
            },
            expectedError: "user already exists",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            
            mockRepo := mocks.NewMockUserRepository(ctrl)
            tt.mockSetup(mockRepo)
            
            service := NewUserService(mockRepo, nil)
            
            // Act
            user, err := service.CreateUser(context.Background(), tt.request)
            
            // Assert
            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
                assert.Nil(t, user)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, user)
                assert.Equal(t, tt.request.Email, user.Email)
            }
        })
    }
}
```

### Test Categories

1. **Unit Tests**: Test individual functions/methods in isolation
2. **Integration Tests**: Test component interactions
3. **End-to-End Tests**: Test complete user workflows

### Coverage Requirements

- **Minimum coverage**: 80% for new code
- **Service layer**: 90% coverage (business logic is critical)
- **Handler layer**: 80% coverage
- **Repository layer**: 70% coverage

```bash
# Run tests with coverage
go test -v -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Mock Generation

We use `gomock` for generating mocks:

```bash
# Install mockgen
go install github.com/golang/mock/mockgen@latest

# Generate mocks
//go:generate mockgen -source=user_repository.go -destination=mocks/mock_user_repository.go
```

## Documentation

### API Documentation

We use Swagger/OpenAPI for API documentation:

```go
// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.CreateUserRequest true "User creation details"
// @Success 201 {object} response.Response{data=domain.User}
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // implementation
}
```

### Code Documentation

```go
// Package service provides business logic implementations for the e-commerce API.
// It contains service interfaces and their implementations for handling
// users, products, and orders.
package service

// UserService defines the interface for user-related business operations.
// It handles user registration, authentication, and profile management.
type UserService interface {
    // Register creates a new user account with the provided details.
    // It validates the input, checks for existing users, hashes the password,
    // and stores the user in the database.
    Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, *common.AppError)
    
    // Login authenticates a user with email and password.
    // It returns a JWT token on successful authentication.
    Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, *common.AppError)
}
```

### README Updates

When adding new features, update relevant documentation:

- **README.md**: Main project documentation
- **API.md**: API endpoint documentation
- **ARCHITECTURE.md**: Architecture changes
- **DEPLOYMENT.md**: Deployment-related changes

## Pull Request Process

### Before Submitting

1. **Ensure all tests pass:**
   ```bash
   make test
   ```

2. **Run linter:**
   ```bash
   make lint
   ```

3. **Update documentation:**
   ```bash
   make swagger
   ```

4. **Rebase on latest main:**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

### PR Description Template

```markdown
## Description
Brief description of the changes made.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Code refactoring

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed

## Documentation
- [ ] API documentation updated
- [ ] README updated (if applicable)
- [ ] Architecture documentation updated (if applicable)

## Screenshots (if applicable)
Add screenshots to help explain your changes.

## Checklist
- [ ] My code follows the project's coding standards
- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
```

### Review Process

1. **Automated checks must pass** (tests, linting, build)
2. **At least one code review** from a maintainer
3. **Documentation review** if applicable
4. **Manual testing** for significant changes

### Merge Requirements

- All CI checks passing
- Code review approval
- Up-to-date with main branch
- No merge conflicts

## Issue Reporting

### Bug Reports

Use the bug report template:

```markdown
**Bug Description**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected Behavior**
A clear and concise description of what you expected to happen.

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Environment:**
- OS: [e.g. Ubuntu 22.04]
- Go Version: [e.g. 1.23.3]
- Database: [e.g. PostgreSQL 15]

**Additional Context**
Add any other context about the problem here.
```

### Feature Requests

Use the feature request template:

```markdown
**Is your feature request related to a problem?**
A clear and concise description of what the problem is.

**Describe the solution you'd like**
A clear and concise description of what you want to happen.

**Describe alternatives you've considered**
A clear and concise description of any alternative solutions or features you've considered.

**Additional context**
Add any other context or screenshots about the feature request here.
```

### Security Issues

For security vulnerabilities, please email directly to the maintainers rather than opening a public issue.

## Community

### Communication Channels

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and community discussions
- **Pull Requests**: Code contributions and reviews

### Getting Help

1. **Check existing documentation**
2. **Search existing issues**
3. **Ask in GitHub Discussions**
4. **Create a new issue** if needed

### Recognition

Contributors will be recognized in:
- **CONTRIBUTORS.md** file
- **Release notes** for significant contributions
- **GitHub contributors** section

---

Thank you for contributing to the E-commerce API project! Your contributions help make this project better for everyone.