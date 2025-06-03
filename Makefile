.PHONY: all build run test clean migrate-up migrate-down sqlc-generate swagger-generate docker-up docker-down go-mod-tidy seed-data

APP_NAME := external-backend-go
BIN_DIR := bin
DB_URL := postgresql://user:password@localhost:5432/mydatabase?sslmode=disable
MIGRATIONS_DIR := db/migrations
SQLC_CONFIG := sqlc.yaml
SWAGGER_MAIN_FILE := cmd/api/main.go

all: build

build:
	@echo "Building Go application..."
	go build -o $(BIN_DIR)/$(APP_NAME) ./$(SWAGGER_MAIN_FILE)
	@echo "Build complete: $(BIN_DIR)/$(APP_NAME)"

run: build
	@echo "Running Go application locally..."
	@$(BIN_DIR)/$(APP_NAME)

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning up generated files..."
	rm -rf $(BIN_DIR)
	rm -rf docs
	rm -rf db/sqlc
	@echo "Cleanup complete."

migrate-up:
	@echo "Waiting for database to start (10 seconds)..."
	@sleep 10
	@echo "Running database migrations (up)..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up
	@echo "Migrations applied."

migrate-down:
	@echo "Running database migrations (down)..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down
	@echo "Migrations rolled back."

sqlc-generate:
	@echo "Generating SQLC code..."
	sqlc generate -f $(SQLC_CONFIG)
	@echo "SQLC code generated in db/sqlc/."

swagger-generate:
	@echo "Generating Swagger documentation..."
	export PATH="$(PATH):$(shell go env GOPATH)/bin" && swag init -g $(SWAGGER_MAIN_FILE)
	@echo "Swagger documentation generated in docs/."

docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d --build
	@echo "Docker services started."

docker-down:
	@echo "Stopping and removing Docker containers..."
	docker-compose down
	@echo "Docker containers stopped and removed."

go-mod-tidy:
	@echo "Tidying Go modules..."
	go mod tidy
	@echo "Go modules tidied."

seed-data:
	@echo "Inserting sample data into database..."
	go run ./cmd/api/main.go --seed
	@echo "Sample data inserted."

setup: go-mod-tidy sqlc-generate swagger-generate docker-up migrate-up
	@echo "Initial Go setup complete. Application is running in Docker."

