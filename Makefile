# Makefile for WMS Backend

.PHONY: build test clean run migrate swagger lint

# Build the application
build:
	go build -o bin/wms-backend main.go

# Run tests
test:
	go test -v ./tests/...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run the application
run:
	go run main.go

# Run migrations
migrate:
	go run main.go -migrate

# Generate Swagger documentation
swagger:
	swag init -g main.go --output docs

# Run linter
lint:
	golangci-lint run

# Download dependencies
deps:
	go mod download
	go mod tidy

# Install development tools
tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker commands
docker-build:
	docker build -t wms-backend:latest .

docker-run:
	docker run -p 8080:8080 --env-file .env wms-backend:latest

# Database commands
db-up:
	docker-compose up -d postgres

db-down:
	docker-compose down

# Development setup
dev-setup:
	cp .env.example .env
	go mod download
	@echo "Setup complete. Edit .env with your configuration."

# Help
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make run            - Run the application"
	@echo "  make migrate        - Run database migrations"
	@echo "  make swagger        - Generate Swagger docs"
	@echo "  make lint           - Run linter"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make deps           - Download dependencies"
	@echo "  make tools          - Install dev tools"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make dev-setup      - Setup development environment"
