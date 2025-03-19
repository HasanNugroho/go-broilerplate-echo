include .env

# Build the application
all: build test

# build the application
build:
	@echo "Building..."
	
	@go build -o main cmd/main.go

# Run the application
run:
	@go run cmd/api/main.go

# watch the application (dev only)
watch:
	@air -c .air.toml

setup:
	go mod download
	cp .env.example .env

# setup database
setup-db:
	docker compose --env-file ./.env -f ./docker-compose.db.yml config
	docker compose --env-file ./.env -f ./docker-compose.db.yml up --build -d

# Shutdown db container
setup-db-down:
	@docker compose -f ./docker-compose.db.yml down --rmi all

# build dev container
docker-run:
	docker compose --env-file ./.env -f ./docker-compose.dev.yml config
	docker compose --env-file ./.env -f ./docker-compose.dev.yml up --build -d

# Shutdown dev container
docker-down:
	@docker compose -f ./docker-compose.dev.yml down --rmi all

create_migration:
	migrate create -ext=sql -dir=migrations -seq $(desc)

migrate_up:
	migrate -path=migrations -database "postgresql://${DBUSER}:${DBPASS}@${DBHOST}:${DBPORT}/${DBNAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=migrations -database "postgresql://${DBUSER}:${DBPASS}@${DBHOST}:${DBPORT}/${DBNAME}?sslmode=disable" -verbose down

.PHONY: create_migration migrate_up migrate_down