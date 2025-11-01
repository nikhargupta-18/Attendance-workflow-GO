.PHONY: help build run docker-up docker-down docker-logs migrate clean

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application locally"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make docker-logs  - View Docker logs"
	@echo "  make migrate      - Run database migrations"
	@echo "  make clean        - Clean build files"

build:
	@echo "Building application..."
	go build -o attendance-api ./cmd/server

run:
	@echo "Running application..."
	go run ./cmd/server/main.go

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs:
	@echo "Viewing Docker logs..."
	docker-compose logs -f api

migrate:
	@echo "Running migrations..."
	go run ./cmd/server/main.go

clean:
	@echo "Cleaning build files..."
	rm -f attendance-api
	go clean




