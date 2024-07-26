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
	migrate -source file://db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migratedown:
	@echo "Migrating-down database"
	migrate -source file://db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

migratedirty:
	@echo "Fix migration dirty"
	migrate -source file://db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" force 1

.PHONY:
start-db:
	@echo "Running the database..."
	docker-compose up -d

# For CRUD API Make
# Create a new person
post-person:
	@echo "Post a person...jimin"
	@curl -s -X POST "${APP_HOST}:${APP_PORT}"/persons \
		-H "Content-Type: application/json" \
		-d '{"Key": "jimin", "Name": "Jimin", "Description": "She’s known for her incredible dancing skills and vibrant personality.", "Image": "path/to/jimin.png", "Traits": [{"Personality": "ESFJ", "Like": "dancing, singing", "Zodiac Sign": "Gemini", "Emoji": "💃", "Color": "Red"}], "Tags": ["jimin", "dancer"]}' \
		| jq .

# Update an existing person
put-person:
	@echo "Put a person...jimin"
	@curl -s -X PUT "${APP_HOST}:${APP_PORT}"/person/jimin \
		-H "Content-Type: application/json" \
		-d '{"Key": "jimin", "Name": "Jimin", "Description": "Updated description", "Image": "path/to/jimin_updated.png", "Traits": [{"Personality": "ESFJ", "Like": "dancing, singing, traveling", "Zodiac Sign": "Gemini", "Emoji": "💃", "Color": "Purple"}], "Tags": ["jimin", "dancer", "traveler"]}' \
		| jq .

# Get a persons all
get-persons:
	@echo "Getting a persons..."
	@curl -s "${APP_HOST}:${APP_PORT}"/persons \
		| jq .

# Get a person by key
get-person:
	@echo "Getting a person...jimin"
	@curl -s "${APP_HOST}:${APP_PORT}"/person/jimin \
		| jq .

# Delete a person by key
delete-person:
	@echo "Delete a person...jimin"
	@curl -s -X DELETE "${APP_HOST}:${APP_PORT}"/person/jimin \
		| jq .

# Check if the service is running
check:
	@curl -s -o /dev/null -w "%{http_code}" "${APP_HOST}:${APP_PORT}"/persons | grep -q "200" && echo "Service is running" || echo "Service is not running"
