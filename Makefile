# Load environment variables
include .env

# Default target
.PHONY: all build run watch setup setup-db setup-db-down docker-run docker-down \
        gen-di gen-docs create-migration migrate-up migrate-down clean

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
	@swag init -g cmd/api/main.go

## ---------------- Database ----------------

# Setup & start database container
setup-db:
	@echo "🐘 Starting database container..."
	@docker compose --env-file ./.env -f ./docker-compose.db.yml up --build -d

# Shutdown database container
setup-db-down:
	@echo "🛑 Stopping database container..."
	@docker compose -f ./docker-compose.db.yml down --rmi all

## ---------------- Docker ----------------

# Run development container
docker-run:
	@echo "🐳 Running dev container..."
	@docker compose --env-file ./.env -f ./docker-compose.dev.yml up --build -d

# Shutdown development container
docker-down:
	@echo "🛑 Stopping dev container..."
	@docker compose -f ./docker-compose.dev.yml down --rmi all

## ---------------- Migrations ----------------

# Create new migration file
create-migration:
	@echo "📜 Creating migration: $(desc)..."
	@migrate create -ext=sql -dir=db/migrations -seq $(desc)

# Apply database migrations
migrate-up:
	@echo "⬆️ Running migrations..."
	@migrate -path=db/migrations -database "postgresql://${DBUSER}:${DBPASS}@${DBHOST}:${DBPORT}/${DBNAME}?sslmode=disable" -verbose up

# Rollback database migrations
migrate-down:
	@echo "⬇️ Rolling back migrations..."
	@migrate -path=db/migrations -database "postgresql://${DBUSER}:${DBPASS}@${DBHOST}:${DBPORT}/${DBNAME}?sslmode=disable" -verbose down

## ---------------- Utilities ----------------

# Clean up generated files
clean:
	@echo "🧹 Cleaning up..."
	@rm -rf main docs
