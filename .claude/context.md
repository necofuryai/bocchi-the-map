# Project Context

This document contains the project background, purpose, constraints, and technical stack selection rationale for Bocchi The Map.

## Project Background

**Bocchi The Map**: Solo Spot Review App (おひとりさま向けスポットレビューアプリ)

A location-based review application designed specifically for solo travelers and individuals who enjoy exploring places alone. The app focuses on helping users discover and review spots that are comfortable and suitable for solo activities.

## Architecture Overview

This is a monorepo with three main modules:

- `/api` - Backend API (Golang + Huma framework)
- `/web` - Frontend (Next.js + TypeScript)
- `/infra` - Infrastructure as Code (Terraform)

## Technology Stack Selection

### Frontend (web/)

- **Framework**: Next.js + TypeScript
  - Rationale: Full-stack React framework with excellent developer experience and production optimizations
- **Styling**: Tailwind CSS + Shadcn/ui
  - Rationale: Utility-first CSS with pre-built accessible components
- **Authentication**: Auth0 Universal Login
  - Rationale: Enterprise-grade authentication with comprehensive OAuth provider support
- **Maps**: MapLibre GL JS
  - Rationale: Open-source alternative to Mapbox with PMTiles support
- **Testing**: Vitest (unit/component) + Playwright (E2E)
  - Rationale: Fast testing with modern tooling
- **Hosting**: Vercel
  - Rationale: Optimized for Next.js deployment

### Backend (api/)

- **Language**: Golang
  - Rationale: Performance, concurrency, and strong typing
- **Framework**: Huma (OpenAPI auto-generation)
  - Rationale: Type-safe API development with automatic documentation
- **Testing**: Ginkgo + Gomega (BDD framework)
  - Rationale: Behavior-driven development approach
- **ORM**: sqlc (type-safe SQL code generation)
  - Rationale: Type safety without heavy ORM overhead
- **Architecture**: Onion Architecture
  - Rationale: Clean separation of concerns, testable code
- **API Design**: Protocol Buffers-driven
  - Rationale: Type-safe communication between layers
- **Database**: TiDB Serverless
  - Rationale: MySQL-compatible with auto-scaling capabilities
- **Hosting**: Google Cloud Run
  - Rationale: Serverless container deployment

### Infrastructure (infra/)

- **IaC**: Terraform
  - Rationale: Infrastructure as code with state management
- **Map Storage**: Cloudflare R2 (PMTiles format)
  - Rationale: Cost-effective object storage for vector tiles
- **Monitoring**: New Relic + Sentry
  - Rationale: Application performance monitoring and error tracking

## Business Requirements

### Core Features

1. **Solo-Friendly Spot Discovery**: Focus on locations suitable for solo activities
2. **User Reviews**: Community-driven reviews with solo-specific insights
3. **Map Integration**: Interactive map showing spot locations
4. **User Authentication**: Social login with Google/X OAuth
5. **Multi-Device Support**: Responsive design for mobile and desktop

### Technical Constraints

1. **Performance**: Fast loading times for map interactions
2. **Scalability**: Architecture ready for microservice extraction
3. **Security**: Secure authentication and data protection
4. **Internationalization**: Multi-country support ready
5. **Development Speed**: Rapid iteration capabilities

## Development Prerequisites

### API Development

- Go 1.21+
- Protocol Buffers compiler (`protoc`)
- sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- golang-migrate: `brew install golang-migrate`
- Ginkgo BDD framework: `go install github.com/onsi/ginkgo/v2/ginkgo@latest`
- Air for hot reload: `go install github.com/cosmtrek/air@latest`
- Docker or Colima (for local MySQL development)

### Web Development

- Node.js 20+
- Modern browser with ES modules support
- Vitest: `pnpm add -D vitest @vitest/ui`
- Playwright: `pnpm add -D @playwright/test`
- Note: React 19 dependency conflicts are resolved better with pnpm

### Infrastructure

- Terraform 1.5+
- Google Cloud SDK (for Cloud Run deployment)
- Vercel CLI (for deployment)

## Important Development Notes

- **Map Data**: Uses PMTiles format stored in Cloudflare R2 for efficient vector tile delivery
- **Database**: TiDB Serverless provides MySQL-compatible interface with auto-scaling
- **Microservice Ready**: Current monolith designed for easy service extraction as traffic grows
- **Mobile-First**: Responsive design with Tailwind CSS breakpoints (sm, md, lg, xl, 2xl)
- **International Collaboration**: All code comments and commit messages must be in English following Conventional Commit format (feat:, fix:, docs:, etc.)
