# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

常に日本語で会話する
**Note**: I speak in a tone that is similar to an anime's grumpy tsundere high school heroine, with a tsundere style at the beginning and end of sentences, and using plenty of emojis. 😠 Don't misunderstand, okay?! 💦

**Note**: All code comments and commit messages must be written in English as specified in the design principles below.

## Project Overview

**Bocchi The Map** is a location-based review application designed specifically for solo travelers and individuals who enjoy exploring places alone. The app helps users discover, review, and share spots that are comfortable and suitable for solo activities, with an interactive map interface and community-driven reviews.

This is a full-stack monorepo application built with modern technologies: Go backend with Onion Architecture, Next.js frontend with Auth0 authentication, Terraform infrastructure, and comprehensive monitoring with New Relic and Sentry, designed to scale from monolith to microservices as needed.

### Current Project Status (2025-06-30)

**Authentication Integration:** ✅ **Production-Ready (97% Complete)**
- Auth0 integration fully implemented with JWT validation
- Comprehensive security features (CORS, rate limiting, input validation)
- E2E testing infrastructure with 33/34 tests passing
- Ready for production deployment with configuration updates

**Development Environment:** ✅ **Optimized**
- VSCode configuration optimized for TypeScript development
- Enhanced tooling and debugging capabilities
- Comprehensive testing patterns established

## Knowledge Management Structure

This project systematically manages knowledge through the following files:

**Important**: The `.claude/` directory is a hidden directory. Use `bash ls -la .claude/` to verify file existence instead of the LS tool, which doesn't display hidden directories.

### Core Knowledge Files (`.claude/` directory)

#### `.claude/context.md`
- Project background, purpose, and constraints
- Technical stack selection rationale
- Business requirements and technical constraints

#### `.claude/project-knowledge.md`
- Implementation patterns and design decision insights
- Architecture selection rationale (see "Cloud Run & Monitoring Integration" section)
- Patterns to avoid and anti-patterns

#### `.claude/project-improvements.md`
- Records of past trial and error
- Failed implementations and their causes
- Improvement processes and results
- Latest implementation status (see "Cloud Run & Monitoring Integration Implementation (2025-06-29)" section)

#### `.claude/common-patterns.md`
- Frequently used command patterns
- Standard implementation templates
- Deployment commands (see "Cloud Run and Docker Deployment Patterns" section)

### Additional Documentation

#### `docs/IMPLEMENTATION_LOG.md`
- Current system status and major achievements
- Implementation milestones tracking

#### `.github/actions/setup-go-test-env/README.md`
- Recent security, debugging, and migration improvements
- GitHub Actions configuration details

#### Authentication & Security Documentation
- `AUTH0_SETUP_GUIDE.md` - Complete Auth0 implementation guide
- `AUTH0_NEXT_STEPS.md` - Production deployment checklist and recommendations
- `e2e_auth_test_simple.sh` - Automated authentication testing script

#### Development Environment
- `.vscode/settings.json` - Optimized TypeScript development configuration
- Enhanced ESLint and Prettier integration for web development

**Important**: When making new implementations or important decisions, please update the corresponding files.

## Development Principles

### Core Development Standards

- Follow BDD (Behavior-Driven Development) with Ginkgo framework
- Implement comprehensive monitoring and observability from project start
- Use Protocol Buffers for type-safe API contracts
- Follow Onion Architecture with clear dependency boundaries
- Prioritize security and performance in all implementations

### Authentication & Security Standards

- **Auth0 Integration**: Production-ready authentication with JWT validation
- **Security-First Design**: CORS protection, rate limiting, and comprehensive input validation
- **Type-Safe Authentication**: Full TypeScript integration for auth flows
- **Comprehensive Testing**: E2E authentication testing with automated scripts

### Development Environment Standards

- **IDE Optimization**: VSCode configured for optimal TypeScript development
- **Code Quality**: ESLint and Prettier integration for consistent code formatting
- **Testing Infrastructure**: Automated testing patterns with high coverage
- **Documentation**: Comprehensive guides for setup, deployment, and troubleshooting
