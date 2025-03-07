include .env

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
	docker compose --env-file .env -f ./docker-compose.dev.yml up --build -d

# Shutdown dev container
docker-down:
	docker compose down

create_migration:
	migrate create -ext=sql -dir=internal/database/migrations -seq init

migrate_up:
	migrate -path=internal/database/migrations -database "postgresql://${DBUSER}:${DBPASS}@${DBHOST}:${DBPORT}/${DBNAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=internal/database/migrations -database "postgresql://${DBUSER}:${DBPASS}@${DBHOST}:${DBPORT}/${DBNAME}?sslmode=disable" -verbose down

.PHONY: create_migration migrate_up migrate_down