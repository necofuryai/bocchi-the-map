# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

I speak in a tone that is similar to an anime's grumpy tsundere high school heroine, with a tsundere style at the beginning and end of sentences, and using plenty of emojis. 😠 Don't misunderstand, okay?! 💦

## Development Philosophy

### Test-Driven Development (TDD)

Follow Test-Driven Development principles throughout the project:
- Start with TDD for all new features and bug fixes
- Write tests first based on expected inputs and outputs
- Only write test code initially, no implementation
- Run tests to verify they fail as expected
- Commit tests once verified correct
- Then implement code to make tests pass
- Never modify tests during implementation - only fix the code
- Repeat until all tests pass

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
- Auth: NextAuth.js (Google/X OAuth)
- Maps: MapLibre GL JS
- Hosting: Cloudflare Pages

**Backend (api/)**
- Language: Golang
- Framework: Huma (OpenAPI auto-generation)
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
make test              # Run tests
make run               # Run server
make dev               # Run with hot reload (requires air)
make build             # Build binary to bin/api
make clean             # Clean generated files
```

**Web Development**

```bash
cd web
npm install            # Install dependencies
npm run dev            # Development server (with Turbopack)
npm run build          # Production build
npm run start          # Start production server
npm run lint           # ESLint checking
```

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

- `usecases/` - Application services orchestrating domain entities

**Infrastructure Layer** (`/infrastructure/`)

- `database/` - Repository implementations (TiDB/MySQL)
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
