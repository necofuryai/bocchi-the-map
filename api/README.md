# 🚀 Bocchi The Map API

> **High-performance Go API with Onion Architecture** - Built for scale, designed for maintainability

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Huma](https://img.shields.io/badge/Huma-v2.19.0-FF6B6B?style=flat)](https://huma.rocks/)
[![Protocol Buffers](https://img.shields.io/badge/Protocol_Buffers-Latest-4285F4?style=flat&logo=google)](https://protobuf.dev/)
[![gRPC](https://img.shields.io/badge/gRPC-Latest-40BC86?style=flat&logo=grpc)](https://grpc.io/)
[![Test Coverage](https://img.shields.io/badge/Coverage-85%25-brightgreen?style=flat)](./coverage.html)

A **type-safe, auto-documented REST API** that powers solo-friendly location discovery. Built with modern Go practices, clean architecture principles, and cloud-native patterns for effortless scaling.

## ⚡ Quick Start

```bash
# 🚀 ONE-COMMAND SETUP (Recommended)
make dev-setup              # Start MySQL + migrations + API

# 📋 Manual setup (if needed)
make deps                   # Install dependencies  
make docker-up              # Start MySQL container
make migrate-up             # Run database migrations
make proto                  # Generate type-safe contracts
make dev                    # Start with hot reload 🔥

# API ready at http://localhost:8080
# Interactive docs at http://localhost:8080/docs
```

## 🔐 Authentication Status

**✅ FULLY IMPLEMENTED**
- OAuth integration (Google/X) via Auth.js
- User management API (`POST /api/users`)
- Database schema with user authentication
- Frontend authentication UI complete

**🔄 NEXT STEPS**
- OAuth credentials setup (Google/X Console)
- Frontend-backend integration testing  
- Live authentication flow testing

## 🏗️ Architecture Philosophy

### Pure Onion Architecture
```
🧅 Dependency Flow: Outer → Inner (NEVER the reverse)

┌─────────────────────────────────────────────────────┐
│  🌐 Interfaces (HTTP/gRPC)                         │
│  ┌───────────────────────────────────────────────┐ │
│  │ 🔌 Infrastructure (DB/Cache/External APIs)    │ │
│  │ ┌─────────────────────────────────────────────┐ │ │
│  │ │ 📋 Application (Use Cases/Workflows)       │ │ │
│  │ │ ┌───────────────────────────────────────────┐ │ │ │
│  │ │ │ 💎 Domain (Pure Business Logic)         │ │ │ │
│  │ │ │ • Zero external dependencies            │ │ │ │
│  │ │ │ • 100% unit testable                   │ │ │ │
│  │ │ │ • Framework agnostic                   │ │ │ │
│  │ │ └───────────────────────────────────────────┘ │ │ │
│  │ └─────────────────────────────────────────────┐ │ │
│  └───────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

**Why This Matters:**
- 🧪 **Testable**: Domain logic tests run in milliseconds with zero setup
- 🔄 **Adaptable**: Swap databases, frameworks, or protocols without touching business logic  
- 📈 **Scalable**: Extract microservices by lifting out domain + application layers
- 🛡️ **Maintainable**: Business rules isolated from infrastructure concerns

### Protocol Buffers-First Design

```protobuf
// Define once, generate everywhere
service SpotService {
  rpc CreateSpot(CreateSpotRequest) returns (CreateSpotResponse);
  rpc GetSpot(GetSpotRequest) returns (GetSpotResponse);
  rpc ListSpots(ListSpotsRequest) returns (ListSpotsResponse);
}
```

**Benefits:**
- 🔒 **Type Safety** - Compile-time contract validation
- 📖 **Auto Documentation** - OpenAPI spec generated from .proto files
- 🌐 **Multi-Language** - Share contracts across Go, TypeScript, mobile apps
- ⚡ **Performance** - Binary serialization for internal services

## 📁 Project Structure

```text
api/
├── 🎯 cmd/api/              # Application entrypoint & DI
├── 💎 domain/               # 🏛️ CORE BUSINESS LOGIC
│   ├── entities/            # Business models with validation
│   ├── repositories/        # Data access contracts (interfaces)
│   └── services/            # Complex business rules
├── 📋 application/          # 🔄 USE CASE ORCHESTRATION  
│   └── usecases/            # App services (coordinate domain)
├── 🔌 infrastructure/       # 🛠️ EXTERNAL DEPENDENCIES
│   ├── database/            # Repository implementations
│   └── external/            # Third-party API clients
├── 🌐 interfaces/           # 📡 TRANSPORT LAYER
│   └── http/
│       ├── handlers/        # Request/response translation
│       └── middleware/      # Cross-cutting concerns
├── 🛠️ pkg/                  # 📦 SHARED UTILITIES
│   ├── config/              # Environment-based config
│   └── logger/              # Structured JSON logging
└── 📋 proto/                # 🔧 API CONTRACTS
    ├── spot.proto           # Spot service definitions
    ├── user.proto           # User management
    ├── review.proto         # Review system
    └── common.proto         # Shared types
```

## 🚀 Key Features

### Modern Go Patterns
- **Generics** - Type-safe repositories and services
- **Context Propagation** - Request tracing and cancellation
- **Structured Logging** - JSON logs with correlation IDs
- **Graceful Shutdown** - Clean resource cleanup
- **Health Checks** - Kubernetes-ready liveness/readiness probes

### API Excellence  
- **Auto-Generated Docs** - OpenAPI 3.0 from Protocol Buffers
- **Request Validation** - Automatic input validation with helpful errors
- **Response Streaming** - Efficient large dataset handling
- **API Versioning** - Backward-compatible evolution
- **Rate Limiting** - Built-in protection against abuse

### Production-Ready
- **Observability** - Metrics, tracing, structured logs
- **Security** - JWT auth, input sanitization, CORS
- **Performance** - Connection pooling, query optimization
- **Reliability** - Circuit breakers, retries, timeouts

## 🛠️ Development Workflow

### 1. Protocol-First Development
```bash
# 1. Define your API contract
vim proto/spot.proto

# 2. Generate type-safe code
make proto

# 3. Implement domain logic (business rules)
vim domain/entities/spot.go

# 4. Add use cases (workflows)  
vim application/usecases/spot_usecase.go

# 5. Wire up HTTP handlers
vim interfaces/http/handlers/spot_handler.go
```

### 2. BDD Testing with Ginkgo

```bash
# BDD specs with Ginkgo framework
ginkgo -r                   # Run all BDD specs
ginkgo ./domain/... -v      # Domain layer specs
ginkgo ./application/... -v # Use case specs

# Traditional unit tests
go test ./domain/... -v     # Unit tests (domain layer)
go test ./application/... -v # Integration tests (use cases)

# Coverage report
make test-coverage
```

### 3. Local Development
```bash
# Hot reload development
make dev

# Manual run with custom config
ENV=development go run cmd/api/main.go --port 8080

# Debug mode with verbose logging
LOG_LEVEL=debug make dev
```

## 🔧 Configuration

### Environment Variables
```bash
# 🗄️ Database
TIDB_HOST=gateway01.ap-northeast-1.prod.aws.tidbcloud.com
TIDB_PORT=4000
TIDB_USER=your-username  
TIDB_PASSWORD=your-password
TIDB_DATABASE=bocchi_the_map
TIDB_SSL_MODE=require

# 🚀 Server
PORT=8080
ENV=production                    # development, staging, production
LOG_LEVEL=info                   # debug, info, warn, error
CORS_ORIGINS=https://yourdomain.com

# 📊 Observability  
NEW_RELIC_LICENSE_KEY=your-key
NEW_RELIC_APP_NAME=bocchi-api
SENTRY_DSN=your-dsn

# 🔐 Security
JWT_SECRET=your-jwt-secret
ENCRYPTION_KEY=your-32-byte-key
```

### Config Management
```go
// Environment-based configuration with validation
type Config struct {
    App      AppConfig      `validate:"required"`
    Database DatabaseConfig `validate:"required"`  
    Logger   LoggerConfig   `validate:"required"`
}

// Automatic loading with defaults
cfg, err := config.Load()
```

## 📊 Performance & Monitoring

### Metrics We Track
- **Request Latency** - p50, p95, p99 response times
- **Throughput** - Requests per second by endpoint
- **Error Rates** - 4xx/5xx breakdown with error types
- **Database Performance** - Query times, connection pool stats
- **Business Metrics** - Spots created, reviews submitted, user activity

### Health Endpoints
```bash
# Kubernetes readiness probe
curl http://localhost:8080/health

# Detailed system status
curl http://localhost:8080/health/detailed

# Dependency check
curl http://localhost:8080/health/dependencies
```

## 🧪 Testing Philosophy

### Test Pyramid
```
🔺 E2E Tests (Few)
   ├── Full API workflow tests
   └── Critical user journey validation

🔺🔺 Integration Tests (Some)  
   ├── Use case testing with real databases
   ├── Handler testing with mock dependencies
   └── Repository testing against test DB

🔺🔺🔺 Unit Tests (Many)
   ├── Domain entity validation logic
   ├── Business rule enforcement  
   └── Pure function testing
```

### Test Commands
```bash
# BDD testing with Ginkgo
ginkgo -r                       # Run all BDD specs
ginkgo -r --randomizeAllSpecs   # Randomized test execution
ginkgo -r --race                # Race condition detection

# Traditional testing
make test-unit                   # Fast, no dependencies
make test-integration           # Requires test database
make test-e2e                   # Full stack testing

# Test with race detection
go test -race ./...

# Benchmark performance
go test -bench=. ./domain/...

# Generate coverage report
make test-coverage && open coverage.html
```

## 🚢 Deployment

### Production Build
```bash
# Optimized binary
make build

# Docker container
docker build -t bocchi-api:latest .

# Multi-arch builds
docker buildx build --platform linux/amd64,linux/arm64 .
```

### Docker Development

#### Local Docker Build
```bash
# Build Docker image locally
docker build -t bocchi-api:dev .

# Run with environment variables
docker run -p 8080:8080 \
  -e TIDB_HOST=your-tidb-host \
  -e TIDB_PASSWORD=your-password \
  -e NEW_RELIC_LICENSE_KEY=your-key \
  -e SENTRY_DSN=your-dsn \
  bocchi-api:dev
```

#### Docker Compose (Development)
```bash
# Start MySQL and API together
make docker-up

# Or manually
docker-compose up -d mysql
docker-compose up api
```

### Cloud Run Deployment

#### Automated Deployment Script
```bash
# Build, push, and optionally deploy to Cloud Run
cd api
./scripts/build.sh dev YOUR_PROJECT_ID asia-northeast1

# Script features:
# - Automated Docker build with optimized caching
# - GCR authentication and image push
# - Environment-specific configuration
# - Interactive Cloud Run deployment option
# - Service URL retrieval and health check validation
```

#### Manual Cloud Run Deployment
```bash
# Build and push to Google Container Registry
gcloud auth configure-docker
docker build -t gcr.io/YOUR_PROJECT_ID/bocchi-api:latest .
docker push gcr.io/YOUR_PROJECT_ID/bocchi-api:latest

# Deploy to Cloud Run
gcloud run deploy bocchi-api-dev \
  --image=gcr.io/YOUR_PROJECT_ID/bocchi-api:latest \
  --platform=managed \
  --region=asia-northeast1 \
  --allow-unauthenticated \
  --port=8080 \
  --memory=1Gi \
  --cpu=1 \
  --max-instances=10 \
  --min-instances=0

# For production
gcloud run deploy bocchi-api-prod \
  --image=gcr.io/YOUR_PROJECT_ID/bocchi-api:latest \
  --platform=managed \
  --region=asia-northeast1 \
  --allow-unauthenticated \
  --port=8080 \
  --memory=1Gi \
  --cpu=2 \
  --max-instances=10 \
  --min-instances=1  # Keep warm for production
```

#### Terraform Infrastructure Deployment
```bash
# Deploy complete infrastructure including secrets management
cd infra

# Initialize Terraform
terraform init

# Plan infrastructure changes
terraform plan -var="gcp_project_id=YOUR_PROJECT_ID"

# Apply infrastructure
terraform apply -var="gcp_project_id=YOUR_PROJECT_ID"

# Set secrets in Google Secret Manager
echo "your-tidb-password" | gcloud secrets create tidb-password-dev --data-file=-
echo "your-new-relic-key" | gcloud secrets create new-relic-license-key-dev --data-file=-
echo "your-sentry-dsn" | gcloud secrets create sentry-dsn-dev --data-file=-
```

### Monitoring and Observability

#### New Relic Setup
```bash
# Environment variables for New Relic
NEW_RELIC_LICENSE_KEY=your-license-key
NEW_RELIC_APP_NAME=bocchi-the-map-api

# Metrics endpoint (when enabled)
curl https://your-cloud-run-url/metrics
```

#### Sentry Setup
```bash
# Environment variable for Sentry
SENTRY_DSN=https://your-sentry-dsn@sentry.io/project-id

# Test error reporting
curl -X POST https://your-cloud-run-url/test-error
```

#### Health Checks
```bash
# Basic health check
curl https://your-cloud-run-url/health

# Detailed health with dependencies
curl https://your-cloud-run-url/health/detailed

# Kubernetes readiness probe
curl https://your-cloud-run-url/health/ready
```

## 📚 API Documentation

### Interactive Docs
- **Development**: http://localhost:8080/docs
- **Staging**: https://api-staging.bocchi-map.com/docs  
- **Production**: https://api.bocchi-map.com/docs

### OpenAPI Spec
```bash
# Generate OpenAPI 3.0 spec
make docs

# Export for frontend teams
curl http://localhost:8080/openapi.json > api-spec.json
```

## 🤝 Contributing

### Code Standards
- **gofmt + goimports** - Automated formatting
- **golangci-lint** - Comprehensive linting
- **Conventional Commits** - Semantic commit messages
- **Test Coverage** - Minimum 80% for new code

### Development Setup
```bash
# Install development tools
make setup-dev

# Install BDD testing framework
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install github.com/onsi/gomega@latest

# Pre-commit hooks
pre-commit install

# BDD workflow: Write specs first, then implement
ginkgo generate ./domain/entities/     # Generate BDD spec files
ginkgo -r --fail-fast                  # Run specs (should fail initially)
# Implement code to make specs pass

# Submit changes
git commit -m "feat(spots): add radius-based search"
gh pr create --title "feat: implement geospatial search"
```

## 🎯 Roadmap

- [x] **v1.0** - Core CRUD operations with clean architecture
- [x] **v1.1** - Protocol Buffers integration and auto-docs
- [ ] **v1.2** - Advanced search with geospatial indexing
- [ ] **v1.3** - Real-time notifications with WebSocket support
- [ ] **v2.0** - Microservice extraction with gRPC federation

---

**🚀 Built for scale, optimized for developer happiness**

[📖 Full Documentation](../README.md) • [🐛 Report Issue](https://github.com/necofuryai/bocchi-the-map/issues) • [💬 Discussions](https://github.com/necofuryai/bocchi-the-map/discussions)
