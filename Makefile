
.PHONY: help build run test deps migrate-up migrate-down migrate-create \
        docker-build docker-up docker-down docker-logs docker-restart \
        clean lint fmt dev install-tools

# Variables
APP_NAME := devx-service-backend
BACKEND_DIR := ./backend-service
MIGRATIONS_DIR := $(BACKEND_DIR)/migrations
DOCKER_COMPOSE := docker compose

.DEFAULT_GOAL := help

help: ## Show available commands
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# ==========================
# Backend
# ==========================

build: ## Build the backend application
	@echo "Building $(APP_NAME)..."
	cd $(BACKEND_DIR) && mkdir -p bin
	cd $(BACKEND_DIR) && go build -o bin/project ./cmd/project
	cd $(BACKEND_DIR) && go build -o bin/migrator ./cmd/migrator

run: ## Run backend locally
	cd $(BACKEND_DIR) && go run ./cmd/project

test: ## Run backend tests
	@echo "Running tests..."
	cd $(BACKEND_DIR) && go test -v ./...

deps: ## Download and tidy backend dependencies
	@echo "Installing dependencies..."
	cd $(BACKEND_DIR) && go mod download && go mod tidy

# ==========================
# Database migrations
# ==========================

migrate-up: ## Apply database migrations locally
	@echo "Running migrations up..."
	cd $(BACKEND_DIR) && go run ./cmd/migrator migrations up

migrate-down: ## Roll back database migration locally
	@echo "Rolling back migration..."
	cd $(BACKEND_DIR) && go run ./cmd/migrator migrations down

migrate-create: ## Create migration (NAME=add_users_table)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=your_migration_name"; \
		exit 1; \
	fi
	@NEXT_NUM=$$(ls -1 $(MIGRATIONS_DIR)/*.sql 2>/dev/null | wc -l | xargs); \
	NEXT_NUM=$$((NEXT_NUM + 1)); \
	FILE_NUM=$$(printf "%05d" $$NEXT_NUM); \
	FILENAME="$(MIGRATIONS_DIR)/$${FILE_NUM}_$(NAME).sql"; \
	echo "Creating migration: $$FILENAME"; \
	echo "-- +goose Up" > "$$FILENAME"; \
	echo "-- +goose StatementBegin" >> "$$FILENAME"; \
	echo "-- Write your UP migration here" >> "$$FILENAME"; \
	echo "-- +goose StatementEnd" >> "$$FILENAME"; \
	echo "" >> "$$FILENAME"; \
	echo "-- +goose Down" >> "$$FILENAME"; \
	echo "-- +goose StatementBegin" >> "$$FILENAME"; \
	echo "-- Write your DOWN migration here" >> "$$FILENAME"; \
	echo "-- +goose StatementEnd" >> "$$FILENAME"; \
	echo "Migration created: $$FILENAME"

# ==========================
# Docker
# ==========================

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

docker-up: ## Start all containers (migrations run automatically)
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d --build
	@echo "Application is ready on http://localhost:8080"

docker-down: ## Stop all containers without deleting volumes
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker container logs
	$(DOCKER_COMPOSE) logs -f

docker-restart: docker-down docker-up ## Restart all containers

# ==========================
# Development
# ==========================

clean: ## Remove backend build artifacts
	@echo "Cleaning..."
	cd $(BACKEND_DIR) && rm -rf bin/
	cd $(BACKEND_DIR) && go clean

lint: ## Run backend linter
	@echo "Running linter..."
	cd $(BACKEND_DIR) && golangci-lint run

fmt: ## Format backend Go code
	@echo "Formatting code..."
	cd $(BACKEND_DIR) && go fmt ./...
	cd $(BACKEND_DIR) && gofmt -s -w ./cmd ./internal

dev: ## Run backend with hot reload (requires Air)
	@echo "Starting development mode..."
	cd $(BACKEND_DIR) && air

install-tools: ## Install development tools
	@echo "Installing development tools..."
	cd $(BACKEND_DIR) && go install github.com/air-verse/air@latest
	cd $(BACKEND_DIR) && go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest


