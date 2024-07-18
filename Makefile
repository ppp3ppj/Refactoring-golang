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

.PHONY:
start-db:
	@echo "Running the database..."
	docker-compose up -d
