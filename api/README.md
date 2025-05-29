# Bocchi The Map API

Backend API service for Bocchi The Map, built with Go and following Onion Architecture principles.

## Architecture

This project follows Onion Architecture with the following layers:

```
api/
├── cmd/api/              # Application entry point
├── domain/               # Core business logic
│   ├── entities/         # Domain models
│   ├── repositories/     # Repository interfaces
│   └── services/         # Domain services
├── application/          # Application services
│   └── usecases/         # Use case implementations
├── infrastructure/       # External dependencies
│   ├── database/         # Database implementations
│   └── external/         # Third-party integrations
├── interfaces/           # API layer
│   └── http/
│       ├── handlers/     # HTTP handlers
│       └── middleware/   # HTTP middleware
├── pkg/                  # Shared packages
│   ├── logger/           # Structured logging
│   └── config/           # Configuration
└── proto/                # Protocol Buffer definitions
```

## Setup

### Prerequisites
- Go 1.21+
- Protocol Buffers compiler
- TiDB Serverless account

### Installation

```bash
# Install dependencies
go mod download

# Generate protobuf files
make proto

# Run tests
go test ./...

# Run the server
go run cmd/api/main.go
```

## Environment Variables

```bash
# Database
TIDB_HOST=your-tidb-host
TIDB_PORT=4000
TIDB_USER=your-username
TIDB_PASSWORD=your-password
TIDB_DATABASE=bocchi_the_map

# Server
PORT=8080
ENV=development

# Monitoring
NEW_RELIC_LICENSE_KEY=your-key
SENTRY_DSN=your-dsn
```

## API Documentation

The API documentation is automatically generated using Huma's OpenAPI support. Access it at:
- Development: http://localhost:8080/docs

## Development

### Running locally

```bash
# Start the server with hot reload
air

# Or without hot reload
go run cmd/api/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./domain/services/...
```