# Ecom-API

An e-commerce REST API built with Go, Gin, and GORM. Includes Swagger docs for API reference.

## Features

- User registration, login, and authentication
- Product management for admin users
- Order creation with payment and shipping statuses
- Configurable via environment variables

## Getting Started

1. Clone repository
2. Run `go mod tidy`
3. Create a `.env` file based on `.env.example`
4. Run `make build` to build the API
5. Run `make run` to start the server

## Migrations

Use the Makefile targets:

- `make migrate-up` to apply migrations
- `make migrate-down` to roll back

## Swagger Docs

- Generate with `make swagger`
- Access Swagger UI at `<host>:<port>/swagger/index.html`
- Definitions live in `docs/` directory

## Development

- Use `make dev` to run live reload (requires `air`)
- Use `make lint` to run static analysis
- Use `make test` to run tests
