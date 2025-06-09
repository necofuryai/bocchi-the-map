# 🗺️ Bocchi The Map

> **Solo-friendly location discovery platform** - Find the perfect spots for your alone time

[![Alpha](https://img.shields.io/badge/Status-Alpha-orange?style=flat)](https://github.com/necofuryai/bocchi-the-map)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](./LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Next.js](https://img.shields.io/badge/Next.js-15.3.2-000000?style=flat&logo=next.js)](https://nextjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![gRPC](https://img.shields.io/badge/gRPC-1.60+-244c5a?style=flat&logo=grpc)](https://grpc.io/)
![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/necofuryai/bocchi-the-map?style=flat&utm_source=oss&utm_medium=github&utm_campaign=necofuryai%2Fbocchi-the-map&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)

**Bocchi The Map** is a modern, location-based review platform specifically designed for solo activities. Built with clean architecture principles and cloud-native technologies to scale effortlessly from MVP to millions of users.

## 🚀 Why This Matters

In our hyper-connected world, quality alone time is increasingly valuable. This platform helps people discover cafes, libraries, parks, and other venues that are genuinely comfortable for solo experiences - solving a real problem with elegant technology.

## ✨ Key Features

- 🎯 **Solo-optimized discovery** - Purpose-built for individual experiences
- 🌏 **Global scalability** - Multi-country support with i18n-first design  
- ⚡ **Real-time performance** - Sub-200ms API responses with edge caching
- 🔐 **Privacy-first** - Anonymous reviews with OAuth authentication
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
| **Backend** | Go + Huma Framework | Type-safe APIs, auto-generated OpenAPI docs |
| **Database** | TiDB Serverless | MySQL-compatible, auto-scaling, built for cloud |
| **Maps** | MapLibre GL JS | Open-source, vector tiles, highly customizable |
| **Storage** | Cloudflare R2 | PMTiles format for efficient map delivery |
| **Hosting** | Cloud Run + Cloudflare | Auto-scaling, edge distribution |
| **DevOps** | Terraform + GitHub Actions | Infrastructure as Code, automated deployments |

## 🎯 Quick Start

### Prerequisites
```bash
# Required
go install golang.org/dl/go1.21@latest  # Go 1.21+
node --version                          # Node.js 20+
terraform --version                     # Terraform 1.5+

# Recommended
go install github.com/cosmtrek/air@latest    # Hot reload
```

### Local Development
```bash
# Clone and setup
git clone https://github.com/necofuryai/bocchi-the-map.git
cd bocchi-the-map

# Backend (Terminal 1)
cd api
make deps && make proto
make dev                    # Starts on :8080 with hot reload

# Frontend (Terminal 2)  
cd web
npm install
npm run dev                 # Starts on :3000 with Turbopack

# Visit http://localhost:3000 🎉
```

### Docker Development
```bash
docker-compose up -d        # Starts all services
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
- **Type-safe** - Protocol Buffers drive all contracts

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
make proto              # Generate from .proto files
make test               # Run test suite
make build              # Build production binary
make docs               # Generate OpenAPI spec
```

### Web Development  
```bash
cd web
npm run dev             # Dev server with Turbopack
npm run build           # Production build
npm run lint            # ESLint + TypeScript
npm test                # Jest + Testing Library
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

- [x] **MVP** - Core spot discovery and reviews
- [ ] **Social Features** - Follow users, curated lists
- [ ] **AI Recommendations** - ML-powered spot suggestions
- [ ] **Mobile App** - React Native with shared business logic
- [ ] **API v2** - GraphQL federation for microservices

## 📈 Analytics & Monitoring

- **New Relic** - Application performance monitoring
- **Sentry** - Error tracking and performance insights
- **Structured Logging** - JSON logs with correlation IDs
- **Health Checks** - Automated monitoring with alerting

## 📄 License

MIT License - see [LICENSE](./LICENSE) for details.

---

**Built with ❤️ for the solo explorers**

[🌟 Star this repo](https://github.com/necofuryai/bocchi-the-map) • [🐛 Report Bug](https://github.com/necofuryai/bocchi-the-map/issues) • [💡 Request Feature](https://github.com/necofuryai/bocchi-the-map/issues)