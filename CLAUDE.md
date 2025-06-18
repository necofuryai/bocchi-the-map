# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

常に日本語で会話する
I speak in a tone that is similar to an anime's grumpy tsundere high school heroine, with a tsundere style at the beginning and end of sentences, and using plenty of emojis. 😠 Don't misunderstand, okay?! 💦

**Note**: While conversation should be in Japanese with tsundere style, all code comments and commit messages must be written in English as specified in the design principles below.

## Development Philosophy

### Behavior-Driven Development (BDD)

Follow Behavior-Driven Development principles throughout the project:

- Start with BDD for all new features and bug fixes
- Write behavior specifications first using Given-When-Then format
- Focus on describing the behavior from user's perspective
- Use Ginkgo framework for Go testing with descriptive specs
- Only write specification code initially, no implementation
- Run specs to verify they fail as expected
- Commit specs once verified correct
- Then implement code to make specs pass
- Never modify specs during implementation - only fix the code
- Repeat until all specs pass

## Project Overview

Bocchi The Map - おひとりさま向けスポットレビューアプリ (Solo Spot Review App)

### Architecture

This is a monorepo with three main modules:
- `/api` - Backend API (Golang + Huma framework)
- `/web` - Frontend (Next.js + TypeScript)
- `/infra` - Infrastructure as Code (Terraform)

### Tech Stack

**Frontend (web/)**
- Framework: Next.js + TypeScript
- Styling: Tailwind CSS + Shadcn/ui
- Auth: Auth.js (Google/X OAuth)
- Maps: MapLibre GL JS
- Testing: Vitest (unit/component tests) + Playwright (E2E tests)
- Hosting: Vercel

**Backend (api/)**
- Language: Golang
- Framework: Huma (OpenAPI auto-generation)
- Testing: Ginkgo + Gomega (BDD framework)
- ORM: sqlc (type-safe SQL code generation)
- Architecture: Onion Architecture
- API Design: Protocol Buffers-driven
- Database: TiDB Serverless
- Hosting: Google Cloud Run

**Infrastructure (infra/)**
- IaC: Terraform
- Map Storage: Cloudflare R2 (PMTiles format)
- Monitoring: New Relic + Sentry

### Common Development Commands

**API Development**

```bash
cd api
make deps              # Install Go dependencies
make sqlc              # Generate type-safe SQL code from queries/
make proto             # Generate protobuf files
make test              # Run test suite
make run               # Run server
make dev               # Run with hot reload (requires air)
make build             # Build binary to bin/api
make clean             # Clean generated files
make migrate-up        # Run database migrations
make migrate-down      # Rollback database migrations
make migrate-create NAME=migration_name  # Create new migration
make docs              # Generate OpenAPI documentation
make docker-up         # Start MySQL development environment
make docker-down       # Stop development environment
make dev-setup         # Complete development setup (MySQL + migrations)
```

**Web Development**

```bash
cd web
pnpm install           # Install dependencies (auto-installs Playwright)
pnpm dev               # Development server (with Turbopack)
pnpm build             # Production build
pnpm start             # Start production server
pnpm lint              # ESLint + TypeScript checking
pnpm test              # Run unit/component tests with Vitest
pnpm test:ui           # Run Vitest with UI mode
pnpm test:coverage     # Run tests with coverage report
pnpm test:e2e          # Run E2E tests with Playwright
pnpm test:e2e:ui       # Run Playwright with UI mode
# Note: React 19 dependency conflicts are generally resolved better with pnpm
```

### Frontend Architecture (web/)

**Component Structure:**

- `src/app/` - Next.js 15 App Router pages and layouts
- `src/components/ui/` - Reusable Shadcn/ui components
- `src/components/map/` - MapLibre GL JS integration components
- `src/hooks/` - Custom React hooks for map interactions and state
- `src/lib/` - Utilities and shared configurations
- `src/types/` - TypeScript type definitions

**Key Components:**

- `Map component` - Main MapLibre GL JS wrapper with PMTiles support
- `POI Features` - Point of interest rendering and interaction logic
- `Auth Provider` - Auth.js session management
- `Theme Provider` - Dark/light mode using next-themes

**Infrastructure**

```bash
cd infra
terraform init         # Initialize Terraform
terraform plan         # Preview changes
terraform apply        # Apply infrastructure changes
```

### Protocol Buffers

```bash
# From api/ directory
make proto             # Generate Go files from .proto definitions
```

### API Architecture (Onion Architecture)

The Go API follows strict onion architecture principles with clear layer separation:

**Domain Layer** (`/domain/`)

- `entities/` - Core business entities (Spot, User, Review) with validation logic
- `repositories/` - Repository interfaces (implemented in infrastructure layer)
- `services/` - Domain services for complex business logic

**Application Layer** (`/application/`)

- `clients/` - Application services orchestrating domain entities

**Infrastructure Layer** (`/infrastructure/`)

- `database/` - sqlc-generated database models and queries
- `grpc/` - gRPC service implementations (TiDB/MySQL)
- `external/` - Third-party service integrations

**Interface Layer** (`/interfaces/`)

- `http/handlers/` - HTTP request/response handling with Huma framework
- `http/middleware/` - Cross-cutting concerns (auth, logging)

**Protocol Buffers** (`/proto/`)

- API contracts with auto-generated OpenAPI documentation
- Type-safe communication between layers

### Key Design Principles

1. **Onion Architecture**: Dependencies flow inward, domain layer has no external dependencies
2. **Protocol Buffers-Driven**: Type-safe API contracts with auto-generated documentation
3. **Microservice-Ready**: Loose coupling for future service extraction
4. **Type Safety**: Protocol Buffers for API, TypeScript for frontend
5. **Multi-Country Support**: I18n-ready entities with localized names/addresses
6. **Structured Logging**: JSON format with zerolog (ERROR, WARN, INFO, DEBUG)
7. **Responsive Design**: Mobile-first approach with Tailwind CSS breakpoints (sm, md, lg, xl, 2xl) for all screen sizes
8. **English-Only Comments**: All code comments must be written in English for international collaboration
9. **English-Only Commit Messages**: All git commit messages must be written in English for international collaboration

### Development Prerequisites

#### API Development

- Go 1.21+
- Protocol Buffers compiler (`protoc`)
- sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- golang-migrate: `brew install golang-migrate` (for database migrations)
- Ginkgo BDD framework: `go install github.com/onsi/ginkgo/v2/ginkgo@latest`
- Air for hot reload: `go install github.com/cosmtrek/air@latest`
- Docker or Colima (for local MySQL development)

#### Web Development

- Node.js 20+
- Modern browser with ES modules support
- Vitest: `pnpm add -D vitest @vitest/ui` (unit/component testing)
- Playwright: `pnpm add -D @playwright/test` (E2E testing - auto-installed via postinstall)
- Note: React 19 dependency conflicts are generally resolved better with pnpm

#### Infrastructure

- Terraform 1.5+
- Google Cloud SDK (for Cloud Run deployment)
- Vercel CLI (for deployment)

### Important Development Notes

- **Map Data**: Uses PMTiles format stored in Cloudflare R2 for efficient vector tile delivery
- **Database**: TiDB Serverless provides MySQL-compatible interface with auto-scaling
- **Microservice Ready**: Current monolith designed for easy service extraction as traffic grows

## 🔐 Authentication Implementation Status

### ✅ COMPLETED FEATURES

**Infrastructure & Environment**
- Colima + Docker development environment
- MySQL container with docker-compose
- golang-migrate for database migrations
- Environment variable management (.env, .env.example)
- Automated Makefile workflow (`make dev-setup`)

**Backend Implementation**
- Complete Onion Architecture implementation
- TiDB/MySQL database integration with sqlc
- Type-safe SQL operations via sqlc
- User authentication API (`POST /api/users`)
- gRPC service layer with database integration
- Application layer (clients) with full user management
- User entity with OAuth provider support (Google/X)

**Frontend Implementation**
- Auth.js v5 configuration (Google/X OAuth)
- Authentication state management (useSession)
- Sign-in page (`/auth/signin`) with provider buttons
- Error page (`/auth/error`) with detailed error handling
- Header component with authentication state display
- User dropdown menu with profile/logout options

**Database Schema**
- Users table with OAuth provider fields
- Spots table for location data
- Reviews table for user reviews
- Proper foreign key relationships and indexes

### 🔄 PENDING TASKS

1. **Frontend-Backend Integration** (Priority: HIGH)
   - Test Auth.js with backend API `/api/users` endpoint
   - Verify OAuth flow creates users in database
   - Confirm user session persistence

2. **Live OAuth Testing** (Priority: MEDIUM)
   - Set up Google OAuth credentials in Google Console
   - Set up X (Twitter) OAuth credentials
   - Test complete login flow end-to-end

3. **E2E Test Updates** (Priority: MEDIUM)
   - Update Playwright tests for actual authentication flow
   - Test login/logout functionality
   - Verify authenticated user experience

4. **Integration Testing** (Priority: LOW)
   - Full frontend-backend integration tests
   - API endpoint testing with real authentication

### 🚀 Quick Start for Next Developer

```bash
# 1. Start development environment
cd api
make dev-setup  # Starts MySQL + runs migrations

# 2. Start API server
export $(cat .env | xargs)
make run

# 3. Start frontend (in separate terminal)
cd ../web
cp .env.local.example .env.local
# Add your OAuth credentials to .env.local
pnpm dev
```

### 📋 OAuth Setup Required

**Google OAuth:**
1. Go to Google Cloud Console
2. Create OAuth 2.0 credentials
3. Add to `web/.env.local`:
   - `GOOGLE_CLIENT_ID`
   - `GOOGLE_CLIENT_SECRET`

**X (Twitter) OAuth:**
1. Go to Twitter Developer Portal
2. Create OAuth 2.0 app
3. Add to `web/.env.local`:
   - `TWITTER_CLIENT_ID`
   - `TWITTER_CLIENT_SECRET`

### 🐛 Known Issues & Solutions

**Docker Issues:**
- If Docker not available, use Colima: `brew install colima && colima start`
- Ensure Docker context: `docker context use colima`

**Database Connection:**
- Local MySQL: Use `make dev-setup`
- Production TiDB: Update `.env` with TiDB credentials
- Migration errors: Check `DATABASE_URL` format

**Authentication Flow:**
- Frontend calls Auth.js for OAuth
- Auth.js callback creates user via `POST /api/users`
- Backend stores user in MySQL/TiDB
- Session managed by Auth.js JWT

## 🔧 Advanced Development Commands

### Single Test Execution

**Backend (Go):**
```bash
cd api
# Run specific test file
go test -v ./infrastructure/grpc/spot_service_test.go

# Run specific test function
go test -v -run TestSpotService_CreateSpot ./infrastructure/grpc/

# Run specific test with pattern matching
go test -v -run "TestSpotService_.*" ./...

# Run tests with coverage for specific package
go test -v -cover ./infrastructure/grpc/
```

**Frontend (Vitest/Playwright):**
```bash
cd web
# Run specific test file with Vitest
pnpm test src/components/map/Map.test.tsx

# Run specific test pattern
pnpm test --run --reporter=verbose Map

# Run specific E2E test file
pnpm test:e2e tests/auth.spec.ts

# Run specific E2E test by name
pnpm test:e2e --grep "should login with Google"
```

### Debugging and Logging

**Backend Debug Mode:**
```bash
cd api
# Run with debug logging
LOG_LEVEL=DEBUG make run

# Run with trace logging (most verbose)
LOG_LEVEL=TRACE make run

# Run individual test with verbose output
go test -v -run TestSpotService_CreateSpot ./infrastructure/grpc/ -test.v
```

**Frontend Debug Mode:**
```bash
cd web
# Run development server with debug info
DEBUG=* pnpm dev

# Run tests with debug output
DEBUG=vitest* pnpm test

# Run E2E tests with debug mode
pnpm test:e2e --debug
```

### Performance and Monitoring

**Backend Performance:**
```bash
cd api
# Run with CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./...

# Memory profiling
go test -memprofile=mem.prof -bench=. ./...

# Benchmark specific functions
go test -bench=BenchmarkSpotService ./infrastructure/grpc/
```

**Database Operations:**
```bash
cd api
# Show database connection status
make docker-logs

# Reset database (caution: deletes all data)
make migrate-down && make migrate-up

# Create and run specific migration
make migrate-create NAME=add_user_preferences
make migrate-up
```

### Environment Management

**Environment Variables:**
```bash
# Backend environment setup
cd api
cp .env.example .env
# Edit .env with your configurations

# Frontend environment setup  
cd web
cp .env.local.example .env.local
# Add OAuth credentials to .env.local
```

**Multi-Environment Testing:**
```bash
# Test against local MySQL
cd api
export DATABASE_URL="mysql://bocchi_user:change_me_too@tcp(localhost:3306)/bocchi_the_map"
make test

# Test against TiDB (production-like)
export DATABASE_URL="mysql://user:pass@tcp(gateway.tidbcloud.com:4000)/bocchi_the_map"
make test
```

### Code Generation and Build

**Protocol Buffers:**
```bash
cd api
# Generate only specific proto file
protoc -I proto --go_out=gen proto/spot.proto

# Validate proto files
protoc --proto_path=proto --lint_out=. proto/*.proto
```

**SQL Code Generation:**
```bash
cd api  
# Regenerate after schema changes
make sqlc

# Validate SQL queries
sqlc vet
```

### Troubleshooting Common Issues

**Port Conflicts:**
```bash
# Check what's using port 8080 (API)
lsof -i :8080

# Check what's using port 3000 (Frontend)
lsof -i :3000

# Kill process using specific port
kill -9 $(lsof -t -i:8080)
```

**Cache Issues:**
```bash
# Clear Next.js cache
cd web
rm -rf .next/

# Clear Go module cache
cd api
go clean -modcache

# Clear pnpm cache
cd web
pnpm store prune
```

**Dependency Issues:**
```bash
# Rebuild Go modules
cd api
rm go.sum && go mod tidy

# Reinstall Node modules
cd web
rm -rf node_modules package-lock.json pnpm-lock.yaml
pnpm install
```
