# Ecom-API

An e-commerce REST API built with Go, Gin, and GORM. Comprehensive API documentation is available via Swagger.

## Features

- User registration, login, and authentication
- Product management for admin users
- Order creation with payment and shipping statuses
- Configurable via environment variables

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/ecom-api.git
   ```
2. **Navigate to the project directory:**
   ```bash
   cd ecom-api
   ```
3. **Install dependencies:**
   ```bash
   go mod tidy
   ```
4. **Create a `.env` file based on `.env.example` and configure environment variables.**
5. **Build the API:**
   ```bash
   make build
   ```
6. **Start the server:**
   ```bash
   make run
   ```

## Migrations

Manage database migrations using Makefile targets:

- **Apply migrations:**
  ```bash
  make migrate-up
  ```
- **Roll back migrations:**
  ```bash
  make migrate-down
  ```

## API Documentation

The API is documented using Swagger. To generate and access the documentation:

- **Generate Swagger docs:**
  ```bash
  make swagger
  ```
- **Access Swagger UI at:** `http://<host>:<port>/swagger/index.html`
- **API definitions are located in the `docs/` directory.**

## Development

Use the following Makefile targets for development tasks:

- **Run live reload (requires `air`):**
  ```bash
  make dev
  ```
- **Run static analysis:**
  ```bash
  make lint
  ```
- **Run tests:**
  ```bash
  make test
  ```

## Contributing

Contributions are welcome! Please fork the repository and create a pull request with your changes.

## License

This project is licensed under the MIT License.
