# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- **Authentication Migration**: Migrated from Auth0 to Supabase Auth for simplified authentication
- Complete modernization of authentication system with SSR compatibility
- Updated all documentation to reflect Supabase Auth authentication system
- Enhanced developer experience with simplified configuration

### Added
- Server-side rendering compatible authentication with Next.js App Router
- Cookie-based session management with automatic token refresh
- Type-safe authentication context and components
- Email/password authentication with email verification
- Protected routes and authentication guards
- Comprehensive error handling and user feedback

### Removed
- Auth0 dependencies and configuration (@auth0/nextjs-auth0)
- Auth0-specific authentication middleware and components
- Legacy JWT token handling code
- Auth0 environment variables and configuration

### Fixed
- TypeScript import issues with authentication modules
- ESLint warnings in authentication components
- Middleware compatibility with Next.js App Router

## [2025-06-28] - Security Update

### Fixed
- **Critical**: Fixed Huma v2 authentication middleware silent context propagation failure
- Protected API endpoints now properly authenticate users
- Implemented proper `huma.WithValue()` context handling pattern

### Security
- All authentication systems now fully functional and production-ready
- Enhanced security features for API endpoints
