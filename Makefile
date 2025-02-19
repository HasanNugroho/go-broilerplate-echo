# Build the application
all: build test

# build the application
build:
	@echo "Building..."
	
	@go build -o main cmd/main.go

# Run the application
run:
	@go run cmd/main.go

# watch the application (dev only)
watch:
	@air -c .air.toml

setup:
	go mod download
	cp .env.example .env

# build dev container
docker-run:
	docker compose --env-file .env -f build/docker-compose.dev.yml up --build

# Shutdown dev container
docker-down:
	docker compose down