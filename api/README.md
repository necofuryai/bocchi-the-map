# 🚀 Bocchi The Map API

> **High-performance Go API with Onion Architecture** - Built for scale, designed for maintainability

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
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

**✅ PRODUCTION READY - CRITICAL BUG FIXED (2025-06-28)**
- ✅ **Huma v2 Authentication Middleware**: Fixed critical context propagation issue
- ✅ **OAuth Integration**: Complete Google/X OAuth via Supabase Auth
- ✅ **JWT Authentication**: Secure token generation and validation 
- ✅ **User Management API**: Full CRUD operations with authentication
- ✅ **Protected Endpoints**: `/api/v1/users/me` and preferences properly secured
- ✅ **Review Authentication**: User authentication for review creation
- ✅ **Database Integration**: Complete user authentication schema
- ✅ **Frontend Integration**: Authentication UI and state management

**🎯 RECENT CRITICAL FIX (2025-06-28)**
- **Issue**: Huma v2 middleware context was not being propagated to handlers
- **Impact**: Authentication appeared to work but protected endpoints were actually unprotected  
- **Solution**: Implemented proper `huma.WithValue()` context handling
- **Result**: All protected endpoints now correctly authenticate users

**✅ READY FOR PRODUCTION**
- Authentication system fully functional
- All security endpoints properly protected
- Microservice-ready authentication architecture
- **Protocol Buffers migration completed** - All manual structs replaced
- **Type-safe API contracts** - Full Protocol Buffers implementation

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

### Protocol Buffers-First Design ✅ **FULLY IMPLEMENTED**

```protobuf
// Complete Protocol Buffers implementation
service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
}

service SpotService {
  rpc CreateSpot(CreateSpotRequest) returns (CreateSpotResponse);
  rpc GetSpot(GetSpotRequest) returns (GetSpotResponse);
  rpc ListSpots(ListSpotsRequest) returns (ListSpotsResponse);
}

service ReviewService {
  rpc CreateReview(CreateReviewRequest) returns (CreateReviewResponse);
  rpc GetSpotReviews(GetSpotReviewsRequest) returns (GetSpotReviewsResponse);
  rpc GetUserReviews(GetUserReviewsRequest) returns (GetUserReviewsResponse);
}
```

**✅ MIGRATION COMPLETED:**
- 🏗️ **Full Implementation** - All manual struct definitions replaced with Protocol Buffers
- 🔒 **Type Safety** - Compile-time contract validation across all services
- 📖 **Auto Documentation** - OpenAPI spec generated from .proto files
- 🌐 **Multi-Language Ready** - Share contracts across Go, TypeScript, mobile apps
- ⚡ **Performance** - Binary serialization for internal services
- 🎯 **Zero Breaking Changes** - Seamless migration from manual structs

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
├── 📋 proto/                # 🔧 API CONTRACTS (SOURCE)
│   ├── common.proto         # Shared types & pagination
│   ├── user.proto           # User management service
│   ├── spot.proto           # Spot service definitions
│   └── review.proto         # Review system
└── 🤖 gen/                  # 🔧 GENERATED CODE (DO NOT EDIT)
    ├── common/v1/           # Generated common types
    ├── user/v1/             # Generated user service code
    ├── spot/v1/             # Generated spot service code
    └── review/v1/           # Generated review service code
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

### 1. Protocol-First Development ✅ **ACTIVE WORKFLOW**
```bash
# 1. Define your API contract
vim proto/spot.proto

# 2. Generate type-safe code (generates Go structs & gRPC clients)
make proto

# 3. Implement domain logic using generated types
vim domain/entities/spot.go

# 4. Add use cases with Protocol Buffers integration
vim application/usecases/spot_usecase.go

# 5. Wire up HTTP handlers (auto-converts to/from protobuf)
vim interfaces/http/handlers/spot_handler.go

# 6. Update gRPC services with generated code
vim infrastructure/grpc/spot_service.go
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
# 🗄️ Database Configuration
# Local Development (MySQL 8.0):
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=password
MYSQL_DATABASE=bocchi_the_map

# Production Environment (TiDB Serverless):
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

# Run with environment variables (using MySQL for local development)
docker run -p 8080:8080 \
  -e MYSQL_HOST=host.docker.internal \
  -e MYSQL_PASSWORD=password \
  -e NEW_RELIC_LICENSE_KEY=your-key \
  -e SENTRY_DSN=your-dsn \
  bocchi-api:dev
```

#### Docker Compose (Development)

```bash
# Start MySQL (local development database) and API together
# Note: Uses MySQL 8.0 for local development, TiDB is used in production
# Requires docker-compose.yml in the root directory
make docker-up

# Or manually
docker-compose up -d mysql
docker-compose up api
```

### Cloud Run Deployment

#### Automated Deployment Script

```bash
# Prerequisites:
# - gcloud CLI authenticated: gcloud auth login
# - Project configured: gcloud config set project YOUR_PROJECT_ID
# - Docker permissions: Make script executable: chmod +x scripts/build.sh
# - Docker Buildx available for multi-platform builds

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

# Deploy to Cloud Run (development - allows unauthenticated access for testing)
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

# For production (requires authentication and dedicated service account)
gcloud run deploy bocchi-api-prod \
  --image=gcr.io/YOUR_PROJECT_ID/bocchi-api:latest \
  --platform=managed \
  --region=asia-northeast1 \
  --port=8080 \
  --memory=1Gi \
  --cpu=2 \
  --max-instances=10 \
  --min-instances=1 \
  --service-account=bocchi-api-service-account@YOUR_PROJECT_ID.iam.gserviceaccount.com  # Keep warm for production

# ⚠️ IMPORTANT: Production deployment requires authentication
# The --allow-unauthenticated flag is intentionally omitted for security
# You MUST set up IAM policies and service account authentication before accessing the service
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

# Set secrets in Google Secret Manager (idempotent - safe to re-run)
echo "your-tidb-password" | \
  (gcloud secrets versions add tidb-password-dev --data-file=- --quiet 2>/dev/null || \
   gcloud secrets create tidb-password-dev --data-file=- --quiet)

echo "your-new-relic-key" | \
  (gcloud secrets versions add new-relic-license-key-dev --data-file=- --quiet 2>/dev/null || \
   gcloud secrets create new-relic-license-key-dev --data-file=- --quiet)

echo "your-sentry-dsn" | \
  (gcloud secrets versions add sentry-dsn-dev --data-file=- --quiet 2>/dev/null || \
   gcloud secrets create sentry-dsn-dev --data-file=- --quiet)
```

### Monitoring and Observability

#### New Relic Setup

```bash
# Environment variables for New Relic
NEW_RELIC_LICENSE_KEY=your-license-key
NEW_RELIC_APP_NAME=bocchi-the-map-api

# Metrics endpoint (when enabled)
# ⚠️ WARNING: Metrics may contain sensitive information
# Protect with IAM or API key restrictions if publicly exposed
curl https://your-cloud-run-url/metrics
```

#### Sentry Setup

```bash
# Environment variable for Sentry
SENTRY_DSN=https://your-sentry-dsn@sentry.io/project-id

# Test error reporting (⚠️ DEVELOPMENT ONLY - should not be enabled in production)
# This endpoint should be disabled in production environments for security
curl -X POST https://your-cloud-run-url/test-error
```

#### Health Checks

| Endpoint | Purpose | Response Content | Recommended Use |
|----------|---------|------------------|-----------------|
| `/health` | Basic health check | Simple OK/ERROR status | Cloud Run health check path |
| `/health/ready` | Kubernetes readiness | Service readiness status | K8s readiness probe |
| `/health/detailed` | Detailed system status | Full dependency check | Monitoring and debugging |

```bash
# Basic health check (recommended for Cloud Run --health-check-path)
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
