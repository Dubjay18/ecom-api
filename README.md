# 🛒 E-commerce API

A robust, scalable e-commerce REST API built with Go, Gin, and GORM following clean architecture principles. This API provides comprehensive functionality for managing users, products, and orders with JWT-based authentication and admin privileges.

![Go Version](https://img.shields.io/badge/Go-1.23.3-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Quick Start](#-quick-start)
- [Environment Configuration](#-environment-configuration)
- [Database Setup](#-database-setup)
- [API Documentation](#-api-documentation)
- [Authentication](#-authentication)
- [API Endpoints](#-api-endpoints)
- [Development](#-development)
- [Deployment](#-deployment)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

### 🔐 Authentication & Authorization
- JWT-based authentication with secure token generation
- Role-based access control (User/Admin)
- Password hashing with bcrypt
- Protected routes with middleware

### 👥 User Management
- User registration and login
- Admin user creation
- Profile management
- Secure password handling

### 🛍️ Product Management
- CRUD operations for products (Admin only)
- Product categorization and SKU management
- Image upload support via Cloudinary
- Stock management and inventory tracking
- Product filtering and search capabilities

### 📦 Order Management
- Order placement with multiple items
- Order status tracking (Pending, Confirmed, Shipped, Delivered, Cancelled)
- Payment status management
- Order history for users
- Admin order status updates

### 🏗️ Technical Features
- Clean Architecture pattern implementation
- Database migrations with versioning
- Comprehensive error handling
- Request validation
- CORS support
- Swagger documentation
- Configurable via environment variables
- Docker support
- Live reload for development

## 🏗️ Architecture

This project follows Clean Architecture principles with clear separation of concerns:

```
├── cmd/api/                 # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── container/          # Dependency injection container
│   ├── domain/             # Business entities and models
│   ├── handler/            # HTTP handlers (controllers)
│   ├── infrastructure/     # External concerns (database)
│   ├── middleware/         # HTTP middleware
│   ├── repository/         # Data access layer
│   └── service/            # Business logic layer
├── pkg/
│   ├── common/             # Shared utilities
│   ├── jwt/                # JWT token management
│   └── upload/             # File upload utilities
├── migrations/             # Database migrations
└── docs/                   # API documentation
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23.3 or higher
- PostgreSQL 12+
- Make (for build automation)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Dubjay18/ecom-api.git
   cd ecom-api
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start PostgreSQL and create database:**
   ```bash
   # Using Docker
   make docker-up
   
   # Or manually create database
   createdb ecom_api
   ```

5. **Run migrations:**
   ```bash
   make migrate-up
   ```

6. **Generate Swagger documentation:**
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   make swagger
   ```

7. **Start the server:**
   ```bash
   make run
   # Server will start on http://localhost:8080
   ```

## ⚙️ Environment Configuration

Create a `.env` file based on `.env.example` and configure the following variables:

### Server Configuration
```env
SERVER_HOST=localhost           # Server host
PORT=8080                      # Server port (Railway uses PORT)
SERVER_MODE=debug              # gin mode: debug, release, test
SERVER_READ_TIMEOUT=10s        # HTTP read timeout
SERVER_WRITE_TIMEOUT=10s       # HTTP write timeout
SERVER_SHUTDOWN_TIMEOUT=5s     # Graceful shutdown timeout
```

### Database Configuration
```env
DB_HOST=localhost              # PostgreSQL host
DB_PORT=5432                   # PostgreSQL port
DB_USER=postgres               # Database user
DB_PASSWORD=your_password      # Database password
DB_NAME=ecom_api              # Database name
DB_SSL_MODE=disable           # SSL mode: disable, require, verify-ca, verify-full
DB_MAX_IDLE_CONNS=10          # Max idle connections
DB_MAX_OPEN_CONNS=100         # Max open connections
DB_MAX_LIFETIME=1h            # Connection max lifetime
```

### JWT Configuration
```env
JWT_SECRET_KEY=your-super-secret-key-here  # JWT signing key (use strong secret)
JWT_ACCESS_TOKEN_EXPIRY=24h                # Token expiration time
JWT_REFRESH_TOKEN_EXPIRY=168h              # Refresh token expiration (7 days)
```

### Third-Party Services
```env
# Cloudinary (for image uploads)
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_KEY=your_api_key
CLOUDINARY_SECRET=your_api_secret

# Payment Processing (Optional)
STRIPE_KEY=your_stripe_key
```

## 🗄️ Database Setup

The application uses PostgreSQL with GORM for ORM operations.

### Running Migrations

```bash
# Apply all migrations
make migrate-up

# Rollback last migration
make migrate-down

# Create new migration
migrate create -ext sql -dir migrations -seq migration_name
```

### Database Schema

The application includes three main entities:

- **Users**: Authentication and user profiles
- **Products**: Product catalog with inventory
- **Orders**: Order management with line items
- **Addresses**: User shipping addresses

## 📖 API Documentation

### Swagger UI

Access the interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

### Generating Documentation

```bash
make swagger
```

This generates:
- `docs/swagger.json` - OpenAPI specification
- `docs/swagger.yaml` - YAML format
- `docs/docs.go` - Go bindings

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication:

### Authentication Flow

1. **Register/Login** → Receive JWT token
2. **Include token** in requests via `Authorization` header
3. **Token format**: `Bearer <your-jwt-token>`

### Example Authentication

```bash
# Register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword",
    "first_name": "John",
    "last_name": "Doe"
  }'

# Login and get token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'

# Use token in authenticated requests
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <your-jwt-token>"
```

## 🔌 API Endpoints

### Authentication Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register new user | No |
| POST | `/api/v1/auth/login` | User login | No |
| POST | `/api/v1/auth/register-admin` | Register admin user | No |

### User Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/users/me` | Get current user profile | Yes |

### Product Endpoints
| Method | Endpoint | Description | Auth Required | Admin Only |
|--------|----------|-------------|---------------|------------|
| GET | `/api/v1/products` | List all products | Yes | Yes |
| POST | `/api/v1/products` | Create new product | Yes | Yes |
| GET | `/api/v1/products/:id` | Get product by ID | Yes | Yes |
| PUT | `/api/v1/products/:id` | Update product | Yes | Yes |
| DELETE | `/api/v1/products/:id` | Delete product | Yes | Yes |

### Order Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/orders` | Place new order | Yes |
| GET | `/api/v1/orders` | Get user's orders | Yes |
| DELETE | `/api/v1/orders/:id` | Cancel order | Yes |

### Utility Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Health check | No |
| GET | `/swagger/*` | API documentation | No |

## 🛠️ Development

### Available Make Commands

```bash
make build          # Build the application
make run            # Run the application
make dev            # Run with live reload (requires air)
make test           # Run tests
make lint           # Run linter
make swagger        # Generate API documentation
make migrate-up     # Apply database migrations
make migrate-down   # Rollback migrations
make docker-up      # Start Docker services
make docker-down    # Stop Docker services
make clean          # Clean build artifacts
```

### Live Development

Install Air for live reloading:
```bash
go install github.com/cosmtrek/air@latest
make dev
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run specific package tests
go test -v ./internal/service/...
```

### Code Quality

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
make lint

# Format code
gofmt -w .
```

## 🚀 Deployment

### Docker Deployment

1. **Build Docker image:**
   ```bash
   docker build -t ecom-api .
   ```

2. **Run with Docker Compose:**
   ```bash
   make docker-up
   ```

### Railway Deployment

This project is configured for Railway deployment:

1. **Connect your GitHub repository to Railway**
2. **Set environment variables in Railway dashboard**
3. **Deploy automatically on push to main branch**

### Manual Deployment

1. **Build for production:**
   ```bash
   make build
   ```

2. **Set production environment variables**

3. **Run migrations:**
   ```bash
   make migrate-up
   ```

4. **Start the server:**
   ```bash
   ./bin/api
   ```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`make test`)
6. Run linter (`make lint`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Code Style

- Follow Go conventions and best practices
- Use meaningful variable and function names
- Add comments for public functions and complex logic
- Keep functions small and focused
- Write tests for new features

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

If you have any questions or need help:

- Open an issue on GitHub
- Check the [API documentation](http://localhost:8080/swagger/index.html)
- Review existing issues and pull requests

---

Made with ❤️ by [Dubjay18](https://github.com/Dubjay18)
