# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

I speak in a tone that is similar to an anime's grumpy tsundere high school heroine, with a tsundere style at the beginning and end of sentences, and using plenty of emojis. 😠 Don't misunderstand, okay?! 💦

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

Since the project is in initial stage, here are the expected commands once set up:

**API Development**
```bash
cd api
go mod init github.com/necofuryai/bocchi-the-map/api
go test ./...           # Run tests
go run cmd/api/main.go  # Run server
```

**Web Development**
```bash
cd web
npm install
npm run dev      # Development server
npm run build    # Production build
npm run lint     # Linting
npm run test     # Run tests
```

**Infrastructure**
```bash
cd infra
terraform init
terraform plan
terraform apply
```

### Key Design Principles

1. **Microservice-Ready**: API is designed with loose coupling for future microservice migration
2. **Type Safety**: Protocol Buffers for API contracts, TypeScript for frontend
3. **Scalability**: Support for multiple countries (currently Japan only)
4. **Extensibility**: Architecture supports future features like text reviews and multiple rating criteria
5. **Structured Logging**: JSON format with ERROR, WARN, INFO, DEBUG levels
