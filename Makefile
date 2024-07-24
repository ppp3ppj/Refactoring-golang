#!make
include .env

.PHONY:
run:
	@echo "Running the program..."
	go mod tidy
	go run cmd/main.go

test:
	@echo "Running the tests..."
	go test -v ./... | ./colorize

cover:
	@echo "Running the tests with coverage..."
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out

migrateinit:
	@echo "Creating migration grw_db"
	migrate create -ext sql -dir db/migrations -seq int_schema

migrateup:
	@echo "Migrating-up database"
	migrate -source file://db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up 1

migratedown:
	@echo "Migrating-down database"
	migrate -source file://db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down 1

.PHONY:
start-db:
	@echo "Running the database..."
	docker-compose up -d
