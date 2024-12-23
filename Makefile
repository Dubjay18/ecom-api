.PHONY: all build run test clean swagger lint dev

all: clean swagger build

build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test -v -cover ./...

clean:
	rm -rf bin/
	rm -rf docs/

swagger:
	swag init -g cmd/api/main.go -o ./docs

lint:
	golangci-lint run

dev:
	air -c .air.toml

migrate-up:
	migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down