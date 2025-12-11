# Today Dashboard Makefile

.PHONY: help dev prod restart down logs clean

.DEFAULT_GOAL := help

help: ## Show available commands
	@echo "\033[0;34mToday Dashboard - Available Commands\033[0m"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[0;32m%-12s\033[0m %s\n", $$1, $$2}'

dev: ## Start development environment (frontend + backend)
	@[ -d "node_modules" ] || npm install
	@cd backend/go-backend && go mod download
	@echo "\033[0;34mStarting development environment...\033[0m"
	@echo "Frontend: http://localhost:5173"
	@echo "Backend:  http://localhost:3001"
	@trap 'kill 0' EXIT; \
		(cd backend/go-backend && go run main.go) & \
		npm run dev

prod: ## Start production environment
	@docker network inspect shared-web >/dev/null 2>&1 || docker network create shared-web
	@echo "\033[0;34mBuilding and starting production services...\033[0m"
	@cd caddy && docker compose build && docker compose up -d
	@docker compose build && docker compose up -d --remove-orphans
	@echo "\033[0;32m✓ Production services running\033[0m"
	@echo ""
	@echo "Main: http://localhost"
	@echo "API:  http://localhost/api"

restart: down prod ## Restart production services

down: ## Stop production services
	@echo "\033[0;34mStopping services...\033[0m"
	@docker compose down
	@cd caddy && docker compose down

logs: ## View production logs
	@docker compose logs -f

clean: down ## Stop services and remove build artifacts
	@echo "\033[0;34mCleaning up...\033[0m"
	@rm -rf dist node_modules
	@rm -rf backend/go-backend/app backend/go-backend/go-backend backend/go-backend/data
	@echo "\033[0;32m✓ Cleanup complete\033[0m"
