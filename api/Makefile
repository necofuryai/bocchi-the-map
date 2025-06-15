.PHONY: proto sqlc clean test build run deps dev migrate-up migrate-down docs

# Generate SQL code
sqlc:
	@echo "Generating SQL code..."
	@if command -v sqlc >/dev/null 2>&1; then \
		sqlc generate; \
	else \
		echo "Error: sqlc is not installed. Install with: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"; \
		exit 1; \
	fi

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

# Database migration
migrate-up:
	@echo "Running database migrations..."
	@go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back database migration..."
	@go run cmd/migrate/main.go down

# Generate OpenAPI documentation
docs:
	@echo "Generating OpenAPI documentation..."
	@go run cmd/docs/main.go