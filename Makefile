# TheraClosure Infrastructure Makefile
# Enterprise-grade local infrastructure management

# Variables
PROJECT_NAME := theraclosure
COMPOSE_FILE_INFRA := infra/docker/docker-compose.infrastructure.yml
COMPOSE_FILE_SERVICES := infra/docker/services/docker-compose.services.yml
NETWORK_NAME := $(PROJECT_NAME)-network

# Docker and Docker Compose detection
DOCKER := $(shell command -v docker 2> /dev/null)
DOCKER_COMPOSE := docker compose

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

.PHONY: help check-deps network setup infra services auth build clean logs status stop restart down reset migrate test

# Default target
help: ## Show this help message
	@echo "$(BLUE)TheraClosure Infrastructure Management$(NC)"
	@echo "======================================"
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(YELLOW)Usage examples:$(NC)"
	@echo "  make setup     - Full infrastructure setup"
	@echo "  make infra     - Start core infrastructure only"  
	@echo "  make services  - Start application services"
	@echo "  make auth      - Start auth service only"
	@echo "  make logs      - View all service logs"
	@echo "  make clean     - Stop all and clean up"

# Check dependencies
check-deps: ## Check if required dependencies are installed
	@echo "$(BLUE)Checking dependencies...$(NC)"
ifndef DOCKER
	@echo "$(RED)Error: Docker is not installed$(NC)"
	@exit 1
endif
	@docker compose version >/dev/null 2>&1 || (echo "$(RED)Error: Docker Compose is not available$(NC)" && exit 1)
	@echo "$(GREEN)✓ Docker and Docker Compose are available$(NC)"

# Create Docker network if it doesn't exist
network: check-deps ## Create Docker network
	@echo "$(BLUE)Creating Docker network: $(NETWORK_NAME)$(NC)"
	@docker network inspect $(NETWORK_NAME) >/dev/null 2>&1 || \
		docker network create $(NETWORK_NAME)
	@echo "$(GREEN)✓ Network $(NETWORK_NAME) ready$(NC)"

# Setup complete infrastructure
setup: check-deps network ## Setup complete infrastructure (recommended for first run)
	@echo "$(BLUE)Setting up TheraClosure infrastructure...$(NC)"
	@make infra
	@echo "$(YELLOW)Waiting for infrastructure to be ready...$(NC)"
	@sleep 10
	@make migrate
	@echo "$(YELLOW)Infrastructure ready, starting services...$(NC)"
	@make services
	@echo "$(GREEN)✓ Complete infrastructure setup finished$(NC)"
	@make status

# Start core infrastructure (PostgreSQL, Redis, etc.)
infra: network ## Start core infrastructure services
	@echo "$(BLUE)Starting core infrastructure...$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) up -d
	@echo "$(GREEN)✓ Core infrastructure started$(NC)"

# Start application services
services: ## Start application services (auth, users, etc.)
	@echo "$(BLUE)Starting application services...$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) -f $(COMPOSE_FILE_SERVICES) up -d --no-recreate
	@echo "$(GREEN)✓ Application services started$(NC)"

# Start only auth service
auth: ## Start only the auth service
	@echo "$(BLUE)Starting auth service...$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) -f $(COMPOSE_FILE_SERVICES) up -d --no-recreate auth-service
	@echo "$(GREEN)✓ Auth service started$(NC)"

# Start only users service
users: ## Start only the users service
	@echo "$(BLUE)Starting users service...$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) -f $(COMPOSE_FILE_SERVICES) up -d --no-recreate users-service
	@echo "$(GREEN)✓ Users service started$(NC)"

# Build services
build: check-deps ## Build all service images
	@echo "$(BLUE)Building service images...$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) -f $(COMPOSE_FILE_SERVICES) build
	@echo "$(GREEN)✓ Service images built$(NC)"

# Build specific service
build-auth: check-deps ## Build auth service image
	@echo "$(BLUE)Building auth service image...$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) -f $(COMPOSE_FILE_SERVICES) build auth-service
	@echo "$(GREEN)✓ Auth service image built$(NC)"

build-users: check-deps ## Build users service image
	@echo "$(BLUE)Building users service image...$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) -f $(COMPOSE_FILE_SERVICES) build users-service
	@echo "$(GREEN)✓ Users service image built$(NC)"

# Run database migrations
migrate: ## Run database migrations
	@echo "$(BLUE)Running database migrations...$(NC)"
	@sleep 5  # Wait for PostgreSQL to be ready
	@docker exec -i $(PROJECT_NAME)-postgres psql -U $(PROJECT_NAME) -d $(PROJECT_NAME)_auth < infra/migrations/001_users.sql 2>/dev/null || echo "$(YELLOW)Migration 001 already applied$(NC)"
	@docker exec -i $(PROJECT_NAME)-postgres psql -U $(PROJECT_NAME) -d $(PROJECT_NAME)_auth < infra/migrations/002_sessions.sql 2>/dev/null || echo "$(YELLOW)Migration 002 already applied$(NC)"
	@echo "$(GREEN)✓ Database migrations completed$(NC)"

# View logs
logs: ## View logs from all services
	@echo "$(BLUE)Showing logs from all services...$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) -f $(COMPOSE_FILE_SERVICES) logs -f

# View logs for specific service
logs-infra: ## View logs from infrastructure services
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) logs -f

logs-services: ## View logs from application services  
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_SERVICES) logs -f

logs-auth: ## View logs from auth service
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_SERVICES) logs -f auth-service

logs-postgres: ## View logs from PostgreSQL
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) logs -f postgres

logs-redis: ## View logs from Redis
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) logs -f redis

# Check status
status: ## Show status of all services
	@echo "$(BLUE)Service Status:$(NC)"
	@echo "$(YELLOW)Infrastructure Services:$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) ps
	@echo ""
	@echo "$(YELLOW)Application Services:$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_SERVICES) ps

# Stop services
stop: ## Stop all services
	@echo "$(BLUE)Stopping all services...$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_SERVICES) stop
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) stop
	@echo "$(GREEN)✓ All services stopped$(NC)"

# Restart services
restart: ## Restart all services
	@echo "$(BLUE)Restarting all services...$(NC)"
	@make stop
	@sleep 2
	@make setup
	@echo "$(GREEN)✓ All services restarted$(NC)"

# Stop and remove containers (keeps volumes)
down: ## Stop and remove containers (keeps data)
	@echo "$(BLUE)Stopping and removing containers...$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_SERVICES) down
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) down
	@echo "$(GREEN)✓ Containers removed$(NC)"

# Complete cleanup
clean: ## Stop all services and clean up containers, networks (keeps volumes)
	@echo "$(BLUE)Cleaning up infrastructure...$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) -f $(COMPOSE_FILE_SERVICES) down --remove-orphans || true
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) down --remove-orphans || true
	@echo "$(GREEN)✓ Cleanup completed$(NC)"

# Nuclear reset (removes everything including volumes)
reset: ## Complete reset - removes all containers, networks, and volumes
	@echo "$(RED)WARNING: This will remove all data including databases!$(NC)"
	@read -p "Are you sure? Type 'yes' to confirm: " confirm && [ "$$confirm" = "yes" ]
	@echo "$(BLUE)Performing complete reset...$(NC)"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_SERVICES) down --volumes --remove-orphans
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE_INFRA) down --volumes --remove-orphans
	@docker network rm $(NETWORK_NAME) 2>/dev/null || true
	@echo "$(GREEN)✓ Complete reset finished$(NC)"

# Health checks
health: ## Check health of all services
	@echo "$(BLUE)Checking service health...$(NC)"
	@echo "$(YELLOW)PostgreSQL:$(NC)"
	@docker exec $(PROJECT_NAME)-postgres pg_isready -U $(PROJECT_NAME) || echo "$(RED)PostgreSQL not ready$(NC)"
	@echo "$(YELLOW)Redis:$(NC)"
	@docker exec $(PROJECT_NAME)-redis redis-cli ping || echo "$(RED)Redis not ready$(NC)"
	@echo "$(YELLOW)Auth Service:$(NC)"
	@curl -f http://localhost:3001/api/v1/health 2>/dev/null && echo "$(GREEN)Auth service healthy$(NC)" || echo "$(RED)Auth service not ready$(NC)"

# Development helpers
dev-setup: ## Setup for development (infra only, no services)
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@make infra
	@make migrate
	@echo "$(GREEN)✓ Development environment ready$(NC)"
	@echo "$(YELLOW)You can now run services locally with 'go run cmd/main.go'$(NC)"

# Test infrastructure
test: ## Test complete infrastructure setup
	@echo "$(BLUE)Testing infrastructure...$(NC)"
	@make setup
	@sleep 10
	@make health
	@echo "$(GREEN)✓ Infrastructure test completed$(NC)"

# Database shell
db-shell: ## Open PostgreSQL shell
	@docker exec -it $(PROJECT_NAME)-postgres psql -U $(PROJECT_NAME) -d $(PROJECT_NAME)_auth

# Redis shell
redis-shell: ## Open Redis shell
	@docker exec -it $(PROJECT_NAME)-redis redis-cli

# Update main $(DOCKER_COMPOSE).yml to point to new structure
update-main-compose: ## Database management
pgadmin: ## Open pgAdmin web interface
	@echo "$(BLUE)Opening pgAdmin...$(NC)"
	@echo "$(YELLOW)pgAdmin is available at: http://localhost:5050$(NC)"
	@echo "$(YELLOW)Email: admin@theraclosure.com$(NC)"
	@echo "$(YELLOW)Password: admin123$(NC)"
	@echo "$(YELLOW)Pre-configured databases:$(NC)"
	@echo "  - theraclosure (main)"
	@echo "  - theraclosure_auth"  
	@echo "  - theraclosure_users"
	@echo "  - theraclosure_payments"
	@echo "  - theraclosure_core"