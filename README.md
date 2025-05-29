# Bocchi The Map 🗺️

おひとりさま向けスポットレビューアプリ / Solo Spot Review App

## Overview

Bocchi The Map is a location-based review application designed for solo activities. Users can discover and review places that are comfortable for enjoying alone time.

## Features

- 🗾 Map-based spot discovery (Japan-focused, with multi-country support planned)
- ⭐ Anonymous 1-5 star rating system
- 🔐 OAuth authentication (Google/X)
- 🌓 Dark mode support
- 🌐 Bilingual support (Japanese/English)

## Tech Stack

### Frontend
- Next.js + TypeScript
- Tailwind CSS + Shadcn/ui
- MapLibre GL JS
- NextAuth.js

### Backend
- Golang + Huma Framework
- Protocol Buffers
- TiDB Serverless
- Google Cloud Run

### Infrastructure
- Terraform
- Cloudflare R2 (Map storage)
- New Relic + Sentry (Monitoring)

## Project Structure

```
bocchi-the-map/
├── api/          # Backend API (Golang)
├── web/          # Frontend (Next.js)
├── infra/        # Infrastructure (Terraform)
└── CLAUDE.md     # AI assistant instructions
```

## Getting Started

### Prerequisites
- Go 1.21+
- Node.js 20+
- Terraform 1.5+

### Development

See individual module READMEs for detailed setup instructions:
- [API Documentation](./api/README.md)
- [Web Documentation](./web/README.md)
- [Infrastructure Documentation](./infra/README.md)

## License

MIT License - see [LICENSE](./LICENSE) file for details.