# Load environment variables
include .env

# Default target
.PHONY: all build run watch setup setup-db setup-db-down docker-run docker-down \
        gen-di gen-docs clean 

all: build test

## ---------------- Build & Run ----------------

# Build the application
build:
	@echo "🔨 Building application..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@echo "🚀 Running application..."
	@go run ./cmd/api

# Watch for changes (dev only)
watch:
	@echo "👀 Watching for changes..."
	@air -c .air.toml

## ---------------- Setup ----------------

# Install dependencies & setup environment
setup:
	@echo "📦 Setting up project..."
	@go mod download
	@cp .env.example .env || true

# Generate dependency injection
gen-di:
	@echo "⚙️ Generating DI..."
	@wire gen ./cmd/api

# Generate Swagger API docs
gen-docs:
	@echo "📖 Generating API documentation..."
	@swag init -g cmd/api/main.go -o cmd/docs

## ---------------- Database ----------------

# Setup & start environment container
env-up:
	@echo "🐘 Starting environment container..."

	@docker compose --env-file ./.env -f ./deploy/docker-compose.env.yml up --build -d

# Shutdown environment container
env-down:
	@echo "🛑 Stopping environment container..."
	@docker compose --env-file ./.env -f ./deploy/docker-compose.env.yml down

## ---------------- Docker ----------------

# Run development container
docker-run:
	@echo "🐳 Running dev container..."
	@docker compose --env-file ./.env -f ./deploy/docker-compose.dev.yml up --build -d

# Shutdown development container
docker-down:
	@echo "🛑 Stopping dev container..."
	@docker compose -f ./deploy/docker-compose.dev.yml down --rmi all

## ---------------- Migrations ----------------

# Clean up generated files
clean:
	@echo "🧹 Cleaning up..."
	@rm -rf main docs
