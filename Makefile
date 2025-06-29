# Subscription Tracker Bot - Makefile
.PHONY: help build up down restart logs clean migrate migrate-down test docker-test

# Variables
DEPLOYMENT_DIR := deployment/docker
COMPOSE_FILE := $(DEPLOYMENT_DIR)/docker-compose.yml
COMPOSE_DEV := $(DEPLOYMENT_DIR)/docker-compose.override.yml
COMPOSE_PROD := $(DEPLOYMENT_DIR)/docker-compose.prod.yml
ENV_FILE := .env
DOCKERFILE := $(DEPLOYMENT_DIR)/Dockerfile

# Help target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Setup commands
setup: ## Initial setup - copy environment file
	@echo "Setting up environment..."
	@if [ ! -f $(ENV_FILE) ]; then \
		cp configs/examples/env.docker.example $(ENV_FILE); \
		echo "Created .env file from configs/examples/env.docker.example"; \
		echo "Please edit .env file with your settings before running 'make up'"; \
	else \
		echo ".env file already exists"; \
	fi

# Docker commands
build: ## Build Docker images
	@echo "Building Docker images..."
	docker-compose -f $(COMPOSE_FILE) build

up: ## Start all services
	@echo "Starting services..."
	docker-compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d

up-dev: ## Start services in development mode
	@echo "Starting services in development mode..."
	docker-compose -f $(COMPOSE_FILE) -f $(COMPOSE_DEV) up -d

up-prod: ## Start services in production mode
	@echo "Starting services in production mode..."
	docker-compose -f $(COMPOSE_FILE) -f $(COMPOSE_PROD) up -d

down: ## Stop all services
	@echo "Stopping services..."
	docker-compose -f $(COMPOSE_FILE) down

restart: ## Restart all services
	@echo "Restarting services..."
	docker-compose -f $(COMPOSE_FILE) restart

logs: ## Show logs from all services
	docker-compose -f $(COMPOSE_FILE) logs -f

logs-bot: ## Show logs from bot service only
	docker-compose -f $(COMPOSE_FILE) logs -f bot

logs-db: ## Show logs from database service only
	docker-compose -f $(COMPOSE_FILE) logs -f postgres

# Database commands
migrate: ## Run database migrations
	@echo "Running database migrations..."
	docker-compose -f $(COMPOSE_FILE) --profile migration run --rm atlas

migrate-status: ## Show migration status
	@echo "Checking migration status..."
	docker-compose -f $(COMPOSE_FILE) --profile migration run --rm atlas migrate status --env docker

# Development commands
shell-bot: ## Open shell in bot container
	docker-compose -f $(COMPOSE_FILE) exec bot sh

shell-db: ## Open psql shell in database
	docker-compose -f $(COMPOSE_FILE) exec postgres psql -U postgres -d sub_cos_counter

# Monitoring commands
ps: ## Show running containers
	docker-compose -f $(COMPOSE_FILE) ps

stats: ## Show container stats
	docker stats $$(docker-compose -f $(COMPOSE_FILE) ps -q)

health: ## Check service health
	@echo "=== Service Health Status ==="
	@docker-compose -f $(COMPOSE_FILE) exec postgres pg_isready -U postgres -d sub_cos_counter || echo "Database: UNHEALTHY"
	@docker-compose -f $(COMPOSE_FILE) exec bot pgrep bot > /dev/null && echo "Bot: HEALTHY" || echo "Bot: UNHEALTHY"

# Cleanup commands
clean: ## Remove containers and networks
	@echo "Cleaning up containers and networks..."
	docker-compose -f $(COMPOSE_FILE) down --remove-orphans

clean-all: ## Remove everything including volumes
	@echo "Removing all containers, networks, and volumes..."
	docker-compose -f $(COMPOSE_FILE) down --volumes --remove-orphans
	docker system prune -f

# Testing commands
test: ## Run tests
	go test ./...

docker-test: ## Run tests inside Docker container
	docker-compose -f $(COMPOSE_FILE) exec bot go test ./...

# Production commands
deploy: setup build migrate up ## Full deployment pipeline
	@echo "Deployment completed!"
	@echo "Check logs with: make logs"

deploy-prod: setup build migrate up-prod ## Full production deployment
	@echo "Production deployment completed!"
	@echo "Check logs with: make logs"

backup-db: ## Backup database
	@echo "Creating database backup..."
	@mkdir -p backups
	docker-compose -f $(COMPOSE_FILE) exec postgres pg_dump -U postgres sub_cos_counter > backups/backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "Backup created in backups/ directory"

# Development helpers
dev-setup: ## Setup for development
	@echo "Setting up development environment..."
	cp configs/examples/env.docker.example .env
	sed -i.bak 's/APP_ENVIRONMENT=production/APP_ENVIRONMENT=development/' .env
	sed -i.bak 's/APP_DEBUG=false/APP_DEBUG=true/' .env
	sed -i.bak 's/LOGGING_LEVEL=info/LOGGING_LEVEL=debug/' .env
	rm -f .env.bak
	@echo "Development environment configured"

rebuild: ## Rebuild and restart services
	@echo "Rebuilding and restarting..."
	docker-compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) down
	docker-compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) build --no-cache
	docker-compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d

# Quick commands
start: up ## Alias for 'up'
stop: down ## Alias for 'down'