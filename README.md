# 🗺️ Bocchi The Map

> **Solo-friendly location discovery platform** - Find the perfect spots for your alone time

[![Alpha](https://img.shields.io/badge/Status-Alpha-orange?style=flat)](https://github.com/necofuryai/bocchi-the-map)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](./LICENSE)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Next.js](https://img.shields.io/badge/Next.js-15.3.2-000000?style=flat&logo=next.js)](https://nextjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![Vercel](https://img.shields.io/badge/Vercel-000000?style=flat&logo=vercel)](https://vercel.com/)
[![gRPC](https://img.shields.io/badge/gRPC-1.60+-244c5a?style=flat&logo=grpc)](https://grpc.io/)
[![Protocol Buffers](https://img.shields.io/badge/Protocol_Buffers-Fully_Implemented-4285F4?style=flat&logo=google)](https://protobuf.dev/)
![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/necofuryai/bocchi-the-map?style=flat&utm_source=oss&utm_medium=github&utm_campaign=necofuryai%2Fbocchi-the-map&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)

**Bocchi The Map** is a modern, location-based review platform specifically designed for solo activities. Built with clean architecture principles and cloud-native technologies to scale effortlessly from MVP to millions of users.

## 🚀 Why This Matters

In our hyper-connected world, quality alone time is increasingly valuable. This platform helps people discover cafes, libraries, parks, and other venues that are genuinely comfortable for solo experiences - solving a real problem with elegant technology.

## ✨ Key Features

- 🎯 **Solo-optimized discovery** - Purpose-built for individual experiences
- 🌏 **Global scalability** - Multi-country support with i18n-first design  
- ⚡ **Real-time performance** - Sub-200ms API responses with edge caching
- 🔐 **Privacy-first** - Secure user authentication with OAuth integration
- 📱 **Progressive Web App** - Native-like experience across all devices
- 🌙 **Accessible design** - Dark mode, screen reader support, WCAG compliance

## 🏗️ Architecture

**Modern, scalable monorepo** designed for microservice evolution:

```text
📦 bocchi-the-map/
├── 🚀 api/          # Go + Huma (Onion Architecture)
├── 🎨 web/          # Next.js 15 + TypeScript
├── ☁️  infra/       # Terraform (Multi-cloud)
└── 📋 docs/         # Architecture decisions & guides
```

### Tech Stack Highlights

| Layer | Technology | Why |
|-------|------------|-----|
| **Frontend** | Next.js 15 + TypeScript | App Router, Turbopack, React Server Components |
| **Backend** | Go + Huma + Protocol Buffers | Type-safe APIs, auto-generated OpenAPI docs, protobuf-driven contracts |
| **Database** | TiDB Serverless | MySQL-compatible, auto-scaling, built for cloud |
| **Maps** | MapLibre GL JS | Open-source, vector tiles, highly customizable |
| **Storage** | Cloudflare R2 | PMTiles format for efficient map delivery |
| **Hosting** | Cloud Run + Vercel | Auto-scaling, edge distribution |
| **Monitoring** | New Relic + Sentry | APM, error tracking, performance insights |
| **DevOps** | Terraform + GitHub Actions | Infrastructure as Code, automated deployments |

## 🎯 Quick Start

### Prerequisites
```bash
# Required
go install golang.org/dl/go1.24@latest  # Go 1.24+
node --version                          # Node.js 20+
terraform --version                     # Terraform 1.5+

# Package managers
npm install -g pnpm                     # pnpm (preferred over npm)

# Recommended
go install github.com/cosmtrek/air@latest    # Hot reload
go install github.com/onsi/ginkgo/v2/ginkgo@latest  # BDD testing
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest # Type-safe SQL
```

### Local Development
```bash
# Clone and setup
git clone https://github.com/necofuryai/bocchi-the-map.git
cd bocchi-the-map

# Backend (Terminal 1)
cd api
make deps && make proto     # Install deps + generate type-safe code from .proto files
make dev                    # Starts on :8080 with hot reload

# Frontend (Terminal 2)  
cd web
pnpm install                # Auto-installs Playwright for E2E tests
pnpm dev                    # Starts on :3000 with Turbopack

# Visit http://localhost:3000 🎉
```

### Database Setup

The project uses environment-specific database initialization files:

- **Development/Testing**: `init-test.sql` - Used automatically with `docker-compose up`
- **Production**: `init-production.sql` - Used with production docker-compose configuration

```bash
# Start database for development/testing
docker-compose up -d

# Start database for production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

The database configuration supports both local development and production environments with appropriate schema initialization for each context.

### Testing

```bash
# API testing (BDD with Ginkgo)
cd api
make test                   # Run all tests
ginkgo -r                   # Run BDD specs

# Web testing
cd web
pnpm test                   # Unit tests
pnpm test:e2e               # E2E tests
```

## 🏛️ Architecture Deep Dive

### Backend: Clean Architecture
```
🧅 Onion Architecture Layers:
├── 💎 Domain     # Pure business logic, zero dependencies
├── 📋 Application # Use cases, orchestrates domain
├── 🔌 Infrastructure # Database, external APIs  
└── 🌐 Interfaces # HTTP handlers, middleware
```

**Key Benefits:**
- **Testable** - Domain logic isolated from infrastructure
- **Scalable** - Ready for microservice extraction
- **Maintainable** - Clear dependency boundaries
- **Type-safe** - ✅ **Protocol Buffers fully implemented** - All manual structs replaced with generated code

### Frontend: Modern React
- **App Router** - File-based routing with layouts
- **Server Components** - Zero-JS components by default
- **Streaming SSR** - Progressive page hydration
- **Edge Runtime** - Deploy anywhere, run everywhere

### Infrastructure: Cloud-Native
- **Multi-cloud ready** - Abstract provider interfaces
- **Auto-scaling** - Zero to millions with no config changes
- **Cost-optimized** - Pay only for actual usage
- **Observability** - Structured logging, metrics, tracing

## 📊 Performance Goals

| Metric | Target | Current |
|--------|--------|---------|
| **API Response Time** | < 200ms p95 | 🎯 |
| **Page Load Time** | < 1.5s | 🎯 |
| **Lighthouse Score** | > 95 | 🎯 |
| **Bundle Size** | < 200kb | 🎯 |

## 🛠️ Development

### API Development
```bash
cd api
make proto              # Generate type-safe Go code from .proto files
make test               # Run test suite (BDD with Ginkgo)
make build              # Build production binary
make docs               # Generate OpenAPI spec from protobuf definitions
```

**Protocol Buffers Development Workflow:**
1. **Define API contracts** - Edit `.proto` files in `proto/` directory
2. **Generate code** - Run `make proto` to generate Go structs and gRPC clients
3. **Implement services** - Use generated types in `infrastructure/grpc/` services
4. **Update handlers** - HTTP handlers automatically convert to/from protobuf types
5. **Test & validate** - Type safety enforced at compile time

### Web Development  
```bash
cd web
pnpm dev                # Dev server with Turbopack
pnpm build              # Production build
pnpm lint               # ESLint + TypeScript checking
pnpm test               # Unit/component tests with Vitest
pnpm test:ui            # Vitest with UI mode
pnpm test:coverage      # Tests with coverage report
pnpm test:e2e           # E2E tests with Playwright
pnpm test:e2e:ui        # Playwright with UI mode
```

### Infrastructure
```bash
cd infra
terraform init          # Initialize providers
terraform plan          # Preview changes
terraform apply         # Deploy infrastructure
```

## 🚢 Deployment

### Automated CI/CD
- **GitHub Actions** - Test, build, deploy pipeline
- **Branch Protection** - PR reviews, status checks required
- **Preview Deployments** - Every PR gets a live environment
- **Blue/Green Deploys** - Zero-downtime production updates

### Manual Deployment
```bash
# Deploy to staging
make deploy-staging

# Deploy to production (requires approval)
make deploy-production
```

### Cloud Run Deployment

#### Prerequisites

```bash
# Install and configure Google Cloud CLI
gcloud auth login

# Get your Google Cloud project ID (if you don't know it)
gcloud config get-value project
# Or list all projects: gcloud projects list

# Set your project ID (replace with your actual project ID)
gcloud config set project YOUR_PROJECT_ID

# Configure Docker for Google Container Registry
gcloud auth configure-docker
```

#### Build and Deploy API

```bash
cd api

# Build and push Docker image
./scripts/build.sh dev YOUR_PROJECT_ID asia-northeast1

# Or manual steps:
docker build -t gcr.io/YOUR_PROJECT_ID/bocchi-api:latest .
docker push gcr.io/YOUR_PROJECT_ID/bocchi-api:latest

# Deploy with Terraform
cd ../infra
terraform init
terraform apply -var="gcp_project_id=YOUR_PROJECT_ID"
```

#### Environment Setup

```bash
# Set required secrets in Google Secret Manager
# Note: You need to provide the actual secret values from your environment
# Example methods to provide secrets:
# 1. From environment variables: echo "$TIDB_PASSWORD" | gcloud secrets create tidb-password-dev --data-file=-
# 2. From file: gcloud secrets create tidb-password-dev --data-file=path/to/secret.txt
# 3. Interactive input (type secret and press Ctrl+D) - Recommended for security:
gcloud secrets create tidb-password-dev --data-file=-
gcloud secrets create new-relic-license-key-dev --data-file=-
gcloud secrets create sentry-dsn-dev --data-file=-

# Deploy Cloud Run service
gcloud run deploy bocchi-api-dev \
  --image=gcr.io/YOUR_PROJECT_ID/bocchi-api:latest \
  --platform=managed \
  --region=asia-northeast1 \
  --allow-unauthenticated
```

<!-- ## 🤝 Contributing

We welcome contributions! This project follows modern open-source practices:

1. **Fork & Clone** - Standard GitHub workflow
2. **Feature Branches** - `feat/your-feature-name`
3. **Conventional Commits** - Semantic commit messages
4. **PR Template** - Guided review process
5. **Automated Testing** - CI runs full test suite 

```bash
# Setup development environment
make setup-dev

# Run full test suite
make test-all

# Submit PR
gh pr create --title "feat: your amazing feature"
``` -->

## 📝 Documentation

- 📖 [API Documentation](./api/README.md) - Onion architecture guide
- 🎨 [Frontend Guide](./web/README.md) - React patterns & components  
- ☁️ [Infrastructure Docs](./infra/README.md) - Terraform modules & deployment
- 🏗️ [Architecture Decisions](./docs/IMPLEMENTATION_LOG.md) - Technical choices & rationale

## 🎯 Roadmap

- [x] **MVP** - Core spot discovery and reviews ✅
- [x] **Authentication System** - Auth0 + JWT with enterprise security ✅ 
- [x] **Production Infrastructure** - Cloud Run + monitoring ✅
- [x] **Huma v2 Integration** - Type-safe APIs with auto-docs ✅
- [x] **Protocol Buffers Migration** - Full protobuf implementation with generated code ✅
- [ ] **Social Features** - Follow users, curated lists
- [ ] **AI Recommendations** - ML-powered spot suggestions
- [ ] **Mobile App** - React Native with shared business logic
- [ ] **API v2** - GraphQL federation for microservices

## 🔐 Latest Updates (2025-06-30)

**✅ Protocol Buffers Migration Completed** 🎉
- **Full Implementation**: All manual struct definitions replaced with generated Protocol Buffers code
- **Type Safety**: 100% compile-time contract validation across all services (User, Spot, Review)
- **Zero Breaking Changes**: Seamless migration preserving existing API behavior
- **Performance**: Binary serialization for internal service communication
- **Multi-Language Ready**: Shared contracts enable future TypeScript client generation

**Previous Database Migration & CI/CD Improvements (2025-06-29)** 🚧
- **Migration Fixes**: Resolved idx_location index conflicts in reviews table migrations
- **GitHub Actions**: Enhanced BDD test security with database URL consistency and debug logging
- **Production Ready**: All migration files synchronized between development and production
- **CI Stability**: Improved test environment setup with better error handling and security measures

**Previous Security Update (2025-06-28)** 🚨
- **Issue**: Huma v2 authentication middleware had silent context propagation failure
- **Impact**: Protected API endpoints were not properly authenticating users
- **Resolution**: Implemented proper `huma.WithValue()` context handling pattern
- **Status**: ✅ **All authentication systems now fully functional and production-ready**

**Enhanced Security Features:**
- ✅ Proper JWT token validation and user context propagation
- ✅ Protected endpoints (`/api/v1/users/me`, preferences, reviews) secured
- ✅ Auth0 Universal Login with comprehensive OAuth provider support
- ✅ Microservice-ready authentication architecture

## 📈 Analytics & Monitoring

### Observability Stack

- **🔍 New Relic** - Application performance monitoring, custom metrics, distributed tracing
- **🚨 Sentry** - Error tracking, performance insights, real-time alerting
- **📊 Structured Logging** - JSON logs with correlation IDs, centralized via Cloud Logging
- **💓 Health Checks** - Kubernetes-ready probes with dependency validation

### Key Metrics Tracked

- **Performance**: API response times (p50, p95, p99), throughput, error rates
- **Business**: Spot discoveries, review submissions, user engagement patterns
- **Infrastructure**: Memory usage, CPU utilization, database connection pools
- **User Experience**: Page load times, frontend errors, conversion funnels

### Monitoring Endpoints

```bash
# Health check
curl https://api.bocchi-map.com/health

# Detailed system status
curl https://api.bocchi-map.com/health/detailed

# Metrics (Prometheus format)
curl https://api.bocchi-map.com/metrics

# Example metrics output:
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",endpoint="/api/spots",status="200"} 1234
http_requests_total{method="POST",endpoint="/api/reviews",status="201"} 567

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",endpoint="/api/spots",le="0.1"} 800
http_request_duration_seconds_bucket{method="GET",endpoint="/api/spots",le="0.5"} 1200
http_request_duration_seconds_sum{method="GET",endpoint="/api/spots"} 145.67
http_request_duration_seconds_count{method="GET",endpoint="/api/spots"} 1234
```

### Alerting & Incident Response

- **Critical Alerts**: > 5% error rate, > 2s p95 latency, dependency failures
- **Escalation**: Slack notifications → PagerDuty → On-call engineer
- **Runbooks**: Automated remediation for common issues

## 📄 License

MIT License - see [LICENSE](./LICENSE) for details.

---

**Built with ❤️ for the solo explorers**

[🌟 Star this repo](https://github.com/necofuryai/bocchi-the-map) • [🐛 Report Bug](https://github.com/necofuryai/bocchi-the-map/issues) • [💡 Request Feature](https://github.com/necofuryai/bocchi-the-map/issues)
