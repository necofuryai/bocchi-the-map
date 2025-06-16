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
make proto             # Generate protobuf files
make test              # Run test suite
make run               # Run server
make dev               # Run with hot reload (requires air)
make build             # Build binary to bin/api
make clean             # Clean generated files
make migrate-up        # Run database migrations
make migrate-down      # Rollback database migrations
make docs              # Generate OpenAPI documentation
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
- Ginkgo BDD framework: `go install github.com/onsi/ginkgo/v2/ginkgo@latest`
- Air for hot reload: `go install github.com/cosmtrek/air@latest`

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
