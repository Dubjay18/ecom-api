.PHONY: all build run test clean swagger lint dev help install-tools db-setup pre-commit test-coverage docker-build docker-up docker-down

# Include environment variables if .env exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Default target
all: clean swagger build

# Help target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build the application
build: ## Build the application binary
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go

# Run the application
run: ## Run the application
	@echo "Starting application..."
	go run cmd/api/main.go

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	go test -v -cover ./...

# Run tests with coverage report
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean: ## Clean build artifacts and generated files
	@echo "Cleaning up..."
	rm -rf bin/
	rm -rf docs/docs.go docs/swagger.json docs/swagger.yaml
	rm -f coverage.out coverage.html

# Generate Swagger documentation
swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g cmd/api/main.go -o ./docs; \
	else \
		echo "swag is not installed. Run 'make install-tools' first."; \
		exit 1; \
	fi

# Run linter
lint: ## Run golangci-lint
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint is not installed. Run 'make install-tools' first."; \
		exit 1; \
	fi

# Run with live reload
dev: ## Run with live reload (requires air)
	@echo "Starting development server with live reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		read -p "Go's 'air' is not installed. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/cosmtrek/air@latest; \
			air; \
		else \
			echo "You chose not to install air. Running without live reload..."; \
			make run; \
		fi; \
	fi

# Install development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/golang/mock/mockgen@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Tools installed successfully!"

# Database setup
db-setup: ## Setup database (create database)
	@echo "Setting up database..."
	@if command -v createdb >/dev/null 2>&1; then \
		createdb ecom_api_dev 2>/dev/null || echo "Database might already exist"; \
		createdb ecom_api_test 2>/dev/null || echo "Test database might already exist"; \
	else \
		echo "PostgreSQL tools not found. Please install PostgreSQL client tools."; \
	fi

# Run database migrations up
migrate-up: ## Apply database migrations
	@echo "Applying database migrations..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up; \
	else \
		echo "migrate tool is not installed. Run 'make install-tools' first."; \
		exit 1; \
	fi

# Run database migrations down
migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down 1; \
	else \
		echo "migrate tool is not installed. Run 'make install-tools' first."; \
		exit 1; \
	fi

# Create new migration
migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	migrate create -ext sql -dir migrations -seq $(NAME)

# Pre-commit checks
pre-commit: ## Run pre-commit checks (lint, test, build)
	@echo "Running pre-commit checks..."
	make lint
	make test
	make swagger
	make build
	@echo "All pre-commit checks passed!"

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	gofmt -w .
	goimports -w .

# Tidy dependencies
tidy: ## Tidy Go modules
	@echo "Tidying Go modules..."
	go mod tidy

# Security check
security: ## Run security checks
	@echo "Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec is not installed. Installing..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t ecom-api:latest .

docker-up: ## Start services with Docker Compose
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

docker-down: ## Stop services with Docker Compose
	@echo "Stopping services with Docker Compose..."
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

docker-dev: ## Start development environment with Docker
	@echo "Starting development environment..."
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# Environment setup
env-setup: ## Setup environment file
	@if [ ! -f .env ]; then \
		echo "Creating .env file from .env.example..."; \
		cp .env.example .env; \
		echo "Please edit .env file with your configuration"; \
	else \
		echo ".env file already exists"; \
	fi

# Full setup for new developers
setup: env-setup install-tools db-setup migrate-up ## Complete setup for new developers
	@echo "Setup complete! You can now run 'make dev' to start development."

# Production build
build-prod: ## Build for production
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o bin/api cmd/api/main.go

# Benchmark tests
bench: ## Run benchmark tests
	@echo "Running benchmark tests..."
	go test -bench=. -benchmem ./...
	