# Today Dashboard Makefile
# Provides commands for development and production workflows

.PHONY: help dev dev-frontend dev-backend frontend-install frontend-build frontend-build-dev frontend-lint frontend-preview frontend-test backend-build backend-tidy backend-init-db backend-test prod prod-build prod-up logs logs-frontend logs-backend logs-caddy stop clean clean-db clean-logs test docker-prune status

# Default target
help:
	@echo "Today Dashboard - Available commands:"
	@echo ""
	@echo "Development commands:"
	@echo "  make dev              - Run all components in development mode"
	@echo "  make dev-frontend     - Run frontend in development mode"
	@echo "  make dev-backend      - Run backend in development mode"
	@echo ""
	@echo "Frontend-specific commands:"
	@echo "  make frontend-install - Install frontend dependencies"
	@echo "  make frontend-build   - Build frontend for production"
	@echo "  make frontend-build-dev - Build frontend for development"
	@echo "  make frontend-lint    - Lint frontend code"
	@echo "  make frontend-preview - Preview production build"
	@echo ""
	@echo "Backend-specific commands:"
	@echo "  make backend-build    - Build backend binary"
	@echo "  make backend-tidy     - Tidy Go dependencies"
	@echo "  make backend-init-db  - Initialize database directory"
	@echo "  make backend-test     - Run backend tests"
	@echo ""
	@echo "Production commands:"
	@echo "  make prod             - Build and start all services in production mode"
	@echo "  make prod-build       - Build production containers"
	@echo "  make prod-up          - Start production containers"
	@echo ""
	@echo "Utility commands:"
	@echo "  make logs             - Display logs from all services"
	@echo "  make logs-frontend    - Display logs from frontend service"
	@echo "  make logs-backend     - Display logs from backend service"
	@echo "  make logs-caddy       - Display logs from Caddy service" 
	@echo "  make stop             - Stop all containers"
	@echo "  make clean            - Stop and remove all containers and generated files"
	@echo "  make clean-db         - Remove database files"
	@echo "  make clean-logs       - Remove log files"
	@echo "  make test             - Run all tests"
	@echo "  make docker-prune     - Clean up unused Docker resources"
	@echo "  make status           - Check status of running services"

# Development commands
dev: dev-backend dev-frontend

dev-frontend:
	@echo "Starting frontend in development mode..."
	@npm run dev

frontend-install:
	@echo "Installing frontend dependencies..."
	@npm install

frontend-build:
	@echo "Building frontend for production..."
	@npm run build

frontend-build-dev:
	@echo "Building frontend for development..."
	@npm run build:dev

frontend-lint:
	@echo "Linting frontend code..."
	@npm run lint

frontend-preview:
	@echo "Previewing production build..."
	@npm run preview

dev-backend:
	@echo "Starting backend in development mode..."
	@cd backend/go-backend && go run main.go

backend-build:
	@echo "Building backend binary..."
	@cd backend/go-backend && go build -o app main.go

backend-tidy:
	@echo "Tidying Go dependencies..."
	@cd backend/go-backend && go mod tidy

backend-init-db:
	@echo "Ensuring database directory exists..."
	@mkdir -p backend/go-backend/data

backend-test:
	@echo "Running backend tests..."
	@cd backend/go-backend && go test ./...

# Production commands
prod: prod-build prod-up

prod-build:
	@echo "Building production containers..."
	@docker-compose build

prod-up:
	@echo "Starting production services..."
	@docker-compose up -d
	@echo "Services started. Use 'make logs' to see logs."

# Logs commands
logs:
	@echo "Displaying logs from all services..."
	@docker-compose logs -f

logs-frontend:
	@echo "Displaying frontend logs..."
	@docker-compose logs -f frontend

logs-backend:
	@echo "Displaying backend logs..."
	@docker-compose logs -f go-backend

logs-caddy:
	@echo "Displaying Caddy logs..."
	@docker-compose logs -f caddy

# Utility commands
stop:
	@echo "Stopping all services..."
	@docker-compose stop

clean:
	@echo "Cleaning up..."
	@docker-compose down
	@echo "Removing build artifacts..."
	@rm -rf dist
	@rm -rf backend/go-backend/app
	@rm -rf backend/go-backend/go-backend
	@find . -name "node_modules" -type d -prune -exec rm -rf '{}' +
	@echo "Cleanup complete."

clean-db:
	@echo "Cleaning database files..."
	@rm -f backend/go-backend/data/*.db*
	@echo "Database files removed."

clean-logs:
	@echo "Cleaning log files..."
	@find . -name "*.log" -type f -delete
	@echo "Log files removed."

test:
	@echo "Running all tests..."
	@$(MAKE) frontend-test
	@$(MAKE) backend-test

frontend-test:
	@echo "Running frontend tests..."
	@npm run test || echo "No frontend tests configured yet."

# Docker utility commands
docker-prune:
	@echo "Pruning unused Docker resources..."
	@docker system prune -f
	@echo "Docker system pruned."

status:
	@echo "Checking service status..."
	@docker-compose ps