# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- **Authentication Migration**: Migrated from Supabase Auth to Clerk for enhanced authentication features
- Complete modernization of authentication system with improved SSR compatibility
- Updated all documentation to reflect Clerk authentication system
- Enhanced developer experience with simplified configuration and better UI components

### Added
- Clerk authentication integration with built-in user management UI
- Social login support (Google, GitHub, and more)
- Multi-factor authentication capabilities
- Built-in user profile management
- Enhanced session management with Clerk middleware
- Improved TypeScript support with Clerk SDK

### Removed
- Supabase Auth dependencies (@supabase/supabase-js, @supabase/ssr)
- Supabase-specific authentication middleware and components
- Custom email verification implementation (now handled by Clerk)
- Supabase environment variables and configuration

### Fixed
- Authentication flow edge cases
- Session persistence across page refreshes
- TypeScript type definitions for user objects

## [2025-07] - Clerk Authentication Migration

### Changed
- **Major Update**: Migrated from Supabase Auth to Clerk authentication
- Replaced all Supabase Auth components with Clerk equivalents
- Updated authentication flow to use Clerk's built-in UI components
- Simplified authentication configuration

### Added
- @clerk/nextjs package for Next.js integration
- Clerk authentication middleware
- Built-in sign-in and sign-up pages
- User profile management through UserButton component

### Removed
- All Supabase Auth related code and configuration
- Custom authentication UI components (replaced by Clerk's built-in components)

## [2025-07-13] - Supabase Auth Migration (Now Superseded)

### Changed
- **Historical Note**: This migration was later superseded by the Clerk authentication migration
- Migrated from Auth0 to Supabase Auth (later replaced with Clerk)

## [2025-06-28] - Security Update

### Fixed
- **Critical**: Fixed Huma v2 authentication middleware silent context propagation failure
- Protected API endpoints now properly authenticate users
- Implemented proper `huma.WithValue()` context handling pattern

### Security
- All authentication systems now fully functional and production-ready
- Enhanced security features for API endpoints
