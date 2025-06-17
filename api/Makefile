.PHONY: all proto sqlc clean test build run deps dev migrate-up migrate-down docs

# Generate SQL code
sqlc:
	@command -v sqlc >/dev/null 2>&1 || { \
		echo "Error: sqlc is not installed. Install with: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest" ; exit 1 ; } ; \
	echo "Generating SQL code..." && sqlc generate

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	@mkdir -p gen
	@if command -v protoc >/dev/null 2>&1; then \
		protoc --go_out=gen --go_opt=paths=source_relative \
			--go-grpc_out=gen --go-grpc_opt=paths=source_relative \
			proto/*.proto; \
	else \
		echo "Error: protoc is not installed"; \
		exit 1; \
	fi

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@rm -rf gen/

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go

# Run the application
run:
	@echo "Running application..."
	@go run cmd/api/main.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Run with hot reload (requires air)
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Error: air is not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi

# Database migration using golang-migrate
migrate-up:
	@echo "Running database migrations..."
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Error: DATABASE_URL environment variable is required"; \
		echo "Example: export DATABASE_URL='mysql://user:password@host:port/database'"; \
		exit 1; \
	fi
	@migrate -path migrations_new -database "$(DATABASE_URL)" up

migrate-down:
	@echo "Rolling back database migration..."
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Error: DATABASE_URL environment variable is required"; \
		echo "Example: export DATABASE_URL='mysql://user:password@host:port/database'"; \
		exit 1; \
	fi
	@migrate -path migrations_new -database "$(DATABASE_URL)" down

migrate-create:
	@echo "Creating new migration..."
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required"; \
		echo "Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@migrate create -ext sql -dir migrations_new -seq $(NAME)

# Generate OpenAPI documentation
docs:
	@echo "Generating OpenAPI documentation..."
	@go run cmd/docs/main.go

# Docker development environment
docker-up:
	@echo "Starting MySQL development environment..."
	@docker-compose up -d mysql
	@echo "Waiting for MySQL to be ready..."
	@sleep 10

docker-down:
	@echo "Stopping development environment..."
	@docker-compose down

docker-logs:
	@echo "Showing MySQL logs..."
	@docker-compose logs -f mysql

# Development workflow
dev-setup: docker-up migrate-up
	@echo "Development environment is ready!"
	@echo "MySQL is running on localhost:3306"
	@echo "Database: bocchi_the_map"
	@echo "User: bocchi_user"

dev-clean: docker-down
	@echo "Development environment cleaned up."