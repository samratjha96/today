# Today Dashboard Makefile

.PHONY: help dev prod prod-build prod-up clean test network

# Default target
help:
	@echo "Today Dashboard - Available commands:"
	@echo "  make dev      - Run all components in development mode"
	@echo "  make prod     - Run all components in production mode"
	@echo "  make network  - Create Docker shared network"
	@echo "  make clean    - Clean all build artifacts and containers"
	@echo "  make test     - Run all tests"

# Network setup
network:
	@echo "Creating shared Docker network..."
	@docker network create shared-web || echo "Network already exists"

# Development
dev:
	@echo "Starting development environment..."
	@echo "Starting backend on port 3001 and frontend on port 8019"
	@cd backend/go-backend && go run main.go & npm run dev

# Production
prod: network prod-build prod-up

prod-build:
	@echo "Building production containers..."
	@docker-compose build
	@cd caddy && docker-compose build

prod-up:
	@echo "Starting production services..."
	@cd caddy && docker-compose up -d
	@docker-compose up -d
	@echo "Services started. Running at http://localhost with API at http://localhost/api"
	@echo "Use 'docker-compose logs -f' to view logs."

# Cleanup
clean:
	@echo "Stopping and removing containers..."
	@docker-compose down
	@cd caddy && docker-compose down
	@echo "Removing build artifacts..."
	@rm -rf dist node_modules
	@rm -rf backend/go-backend/app backend/go-backend/go-backend
	@echo "Cleanup complete."

# Tests
test:
	@echo "Running tests..."
	@cd backend/go-backend && go test ./...
	@npm run test || echo "No frontend tests configured."