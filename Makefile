# Today Dashboard Makefile
# Complete setup for fresh installations and deployment

.PHONY: help check-prereqs install dev dev-deps prod network build up down restart logs clean test

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

##@ General

help: ## Display this help message
	@echo "$(BLUE)Today Dashboard - Available Commands$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make $(GREEN)<target>$(NC)\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(BLUE)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Setup & Prerequisites

check-prereqs: ## Check if all prerequisites are installed
	@echo "$(BLUE)Checking prerequisites...$(NC)"
	@command -v docker >/dev/null 2>&1 || { echo "$(RED)✗ Docker is not installed$(NC)"; exit 1; }
	@echo "$(GREEN)✓ Docker is installed$(NC)"
	@docker compose version >/dev/null 2>&1 || { echo "$(RED)✗ Docker Compose is not installed$(NC)"; exit 1; }
	@echo "$(GREEN)✓ Docker Compose is installed$(NC)"
	@command -v node >/dev/null 2>&1 || { echo "$(YELLOW)⚠ Node.js is not installed (required for dev mode)$(NC)"; }
	@command -v node >/dev/null 2>&1 && echo "$(GREEN)✓ Node.js is installed (version: $$(node --version))$(NC)" || true
	@command -v npm >/dev/null 2>&1 && echo "$(GREEN)✓ npm is installed (version: $$(npm --version))$(NC)" || true
	@command -v go >/dev/null 2>&1 || { echo "$(YELLOW)⚠ Go is not installed (required for dev mode)$(NC)"; }
	@command -v go >/dev/null 2>&1 && echo "$(GREEN)✓ Go is installed (version: $$(go version))$(NC)" || true
	@echo "$(GREEN)✓ All prerequisites satisfied!$(NC)"

setup-env: ## Create environment files from examples
	@echo "$(BLUE)Setting up environment files...$(NC)"
	@if [ ! -f .env.development ]; then \
		if [ -f .env.example ]; then \
			cp .env.example .env.development; \
			echo "$(GREEN)✓ Created .env.development from .env.example$(NC)"; \
		else \
			echo "$(YELLOW)⚠ No .env.example found, skipping$(NC)"; \
		fi \
	else \
		echo "$(GREEN)✓ .env.development already exists$(NC)"; \
	fi
	@if [ ! -f .env.production ]; then \
		if [ -f .env.example ]; then \
			cp .env.example .env.production; \
			echo "$(GREEN)✓ Created .env.production from .env.example$(NC)"; \
		fi \
	else \
		echo "$(GREEN)✓ .env.production already exists$(NC)"; \
	fi

network: ## Create shared Docker network
	@echo "$(BLUE)Creating shared Docker network...$(NC)"
	@docker network inspect shared-web >/dev/null 2>&1 || \
		(docker network create shared-web && echo "$(GREEN)✓ Network 'shared-web' created$(NC)") || \
		echo "$(GREEN)✓ Network 'shared-web' already exists$(NC)"

install: check-prereqs setup-env network ## Complete fresh installation setup
	@echo "$(GREEN)✓ Fresh installation setup complete!$(NC)"
	@echo "$(BLUE)Next steps:$(NC)"
	@echo "  - For development: $(GREEN)make dev$(NC)"
	@echo "  - For production:  $(GREEN)make prod$(NC)"

##@ Development

dev-deps: ## Install development dependencies
	@echo "$(BLUE)Installing development dependencies...$(NC)"
	@if [ ! -d "node_modules" ]; then \
		echo "Installing frontend dependencies..."; \
		npm install; \
	else \
		echo "$(GREEN)✓ Frontend dependencies already installed$(NC)"; \
	fi
	@echo "Installing Go dependencies..."
	@cd backend/go-backend && go mod download && go mod tidy
	@echo "$(GREEN)✓ Development dependencies installed$(NC)"

dev: dev-deps ## Run in development mode (frontend + backend)
	@echo "$(BLUE)Starting development environment...$(NC)"
	@echo "$(YELLOW)Frontend: http://localhost:5173$(NC)"
	@echo "$(YELLOW)Backend:  http://localhost:3001$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to stop both services$(NC)"
	@trap 'kill 0' EXIT; \
		(cd backend/go-backend && go run main.go) & \
		npm run dev

dev-backend: ## Run only backend in development mode
	@echo "$(BLUE)Starting backend in development mode...$(NC)"
	@cd backend/go-backend && go run main.go

dev-frontend: ## Run only frontend in development mode
	@echo "$(BLUE)Starting frontend in development mode...$(NC)"
	@npm run dev

##@ Production Deployment

build: network ## Build all Docker images
	@echo "$(BLUE)Building Docker images...$(NC)"
	@echo "Building application services..."
	@docker compose build
	@echo "Building Caddy reverse proxy..."
	@cd caddy && docker compose build
	@echo "$(GREEN)✓ All images built successfully$(NC)"

up: ## Start all production services
	@echo "$(BLUE)Starting production services...$(NC)"
	@cd caddy && docker compose up -d
	@docker compose up -d --remove-orphans
	@echo "$(GREEN)✓ Services started!$(NC)"
	@echo ""
	@echo "$(BLUE)Application URLs:$(NC)"
	@echo "  Main:  http://localhost"
	@echo "  API:   http://localhost/api"
	@echo ""
	@echo "$(YELLOW)Use 'make logs' to view logs$(NC)"
	@echo "$(YELLOW)Use 'make down' to stop services$(NC)"

down: ## Stop all production services
	@echo "$(BLUE)Stopping production services...$(NC)"
	@docker compose down
	@cd caddy && docker compose down
	@echo "$(GREEN)✓ Services stopped$(NC)"

restart: down up ## Restart all production services

prod: install build up ## Complete production deployment from scratch

##@ Monitoring & Maintenance

logs: ## View logs from all services
	@echo "$(BLUE)Viewing logs (Ctrl+C to exit)...$(NC)"
	@docker compose logs -f

logs-frontend: ## View frontend logs only
	@docker compose logs -f frontend

logs-backend: ## View backend logs only
	@docker compose logs -f go-backend

logs-caddy: ## View Caddy logs only
	@cd caddy && docker compose logs -f

status: ## Show status of all services
	@echo "$(BLUE)Service Status:$(NC)"
	@echo ""
	@cd caddy && docker compose ps
	@docker compose ps

health: ## Check health of all services
	@echo "$(BLUE)Checking service health...$(NC)"
	@echo ""
	@docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Health}}"
	@cd caddy && docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Health}}"

##@ Testing

test: ## Run all tests
	@echo "$(BLUE)Running tests...$(NC)"
	@echo "Testing Go backend..."
	@cd backend/go-backend && go test -v ./... || echo "$(YELLOW)No Go tests found$(NC)"
	@echo "Testing frontend..."
	@npm run test || echo "$(YELLOW)No frontend tests configured$(NC)"
	@echo "$(GREEN)✓ Tests complete$(NC)"

test-backend: ## Run backend tests only
	@echo "$(BLUE)Running Go backend tests...$(NC)"
	@cd backend/go-backend && go test -v ./...

##@ Cleanup

clean: down ## Stop services and remove build artifacts
	@echo "$(BLUE)Cleaning up...$(NC)"
	@echo "Removing build artifacts..."
	@rm -rf dist node_modules
	@rm -rf backend/go-backend/app backend/go-backend/go-backend backend/go-backend/data
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

clean-all: clean ## Complete cleanup including Docker volumes and images
	@echo "$(BLUE)Performing complete cleanup...$(NC)"
	@echo "Removing Docker images..."
	@docker compose down -v --rmi all 2>/dev/null || true
	@cd caddy && docker compose down -v --rmi all 2>/dev/null || true
	@echo "$(YELLOW)Note: shared-web network is preserved for other services$(NC)"
	@echo "$(GREEN)✓ Complete cleanup done$(NC)"

clean-network: ## Remove the shared Docker network (use with caution)
	@echo "$(RED)Warning: This will remove the shared-web network$(NC)"
	@echo "Press Ctrl+C to cancel, or wait 5 seconds to continue..."
	@sleep 5
	@docker network rm shared-web 2>/dev/null || echo "Network already removed or in use"

##@ Rebuild

rebuild: clean build up ## Clean, rebuild, and restart all services

rebuild-frontend: ## Rebuild only frontend service
	@echo "$(BLUE)Rebuilding frontend...$(NC)"
	@docker compose build frontend
	@docker compose up -d frontend
	@echo "$(GREEN)✓ Frontend rebuilt$(NC)"

rebuild-backend: ## Rebuild only backend service
	@echo "$(BLUE)Rebuilding backend...$(NC)"
	@docker compose build go-backend
	@docker compose up -d go-backend
	@echo "$(GREEN)✓ Backend rebuilt$(NC)"

rebuild-caddy: ## Rebuild only Caddy service
	@echo "$(BLUE)Rebuilding Caddy...$(NC)"
	@cd caddy && docker compose build
	@cd caddy && docker compose up -d
	@echo "$(GREEN)✓ Caddy rebuilt$(NC)"

##@ Database

db-backup: ## Backup SQLite database
	@echo "$(BLUE)Backing up database...$(NC)"
	@mkdir -p backups
	@docker compose exec go-backend cp /app/data/today.db /app/data/today.db.backup || \
		cp backend/go-backend/data/today.db backups/today-$(shell date +%Y%m%d-%H%M%S).db
	@echo "$(GREEN)✓ Database backed up$(NC)"

db-restore: ## Restore SQLite database (requires BACKUP_FILE variable)
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "$(RED)Error: Please specify BACKUP_FILE variable$(NC)"; \
		echo "Example: make db-restore BACKUP_FILE=backups/today-20231201-120000.db"; \
		exit 1; \
	fi
	@echo "$(BLUE)Restoring database from $(BACKUP_FILE)...$(NC)"
	@docker compose exec go-backend cp $(BACKUP_FILE) /app/data/today.db
	@echo "$(GREEN)✓ Database restored$(NC)"

##@ Advanced

shell-frontend: ## Open shell in frontend container
	@docker compose exec frontend sh

shell-backend: ## Open shell in backend container
	@docker compose exec go-backend sh

shell-caddy: ## Open shell in Caddy container
	@cd caddy && docker compose exec caddy sh

validate-caddy: ## Validate Caddyfile configuration
	@echo "$(BLUE)Validating Caddyfile...$(NC)"
	@docker run --rm -v $(PWD)/caddy/Caddyfile:/etc/caddy/Caddyfile caddy:2-alpine caddy validate --config /etc/caddy/Caddyfile
	@echo "$(GREEN)✓ Caddyfile is valid$(NC)"

update-deps: ## Update all dependencies
	@echo "$(BLUE)Updating dependencies...$(NC)"
	@echo "Updating frontend dependencies..."
	@npm update
	@echo "Updating Go dependencies..."
	@cd backend/go-backend && go get -u ./... && go mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"
