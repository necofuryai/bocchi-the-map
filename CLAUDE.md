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
- TDD+BDD hybrid testing templates and workflow commands

#### `.claude/tdd-bdd-hybrid.md`
- Comprehensive TDD+BDD hybrid methodology guide
- Three hybrid approaches: Outside-In TDD, Double-Loop TDD, Specification by Example
- Layer-specific testing strategies and practical workflows
- Anti-patterns to avoid and best practices

#### `.claude/tdd-bdd-example.md`
- Complete implementation example of TDD+BDD hybrid approach
- Solo-Friendly Rating feature as reference implementation
- Step-by-step guide from BDD scenarios to TDD implementation
- Executable code examples and test execution patterns

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

- Follow TDD+BDD Hybrid methodology with three distinct approaches
- Use BDD (Behavior-Driven Development) with Ginkgo framework for user-facing behavior
- Apply TDD (Test-Driven Development) for internal implementation details
- Implement comprehensive monitoring and observability from project start
- Use Protocol Buffers for type-safe API contracts
- Follow Onion Architecture with clear dependency boundaries
- Prioritize security and performance in all implementations

### Git & Commit Standards

- **Commit Messages**: Use conventional commit format (feat:, fix:, chore:, etc.)
- **No Auto-Signatures**: Do NOT include Claude Code attribution or co-authorship signatures in commit messages
- **Clean History**: Keep commit messages concise and focused on the actual changes
- **English Only**: All commit messages must be written in English

### Authentication & Security Standards

- **Auth0 Integration**: Production-ready authentication with JWT validation
- **Security-First Design**: CORS protection, rate limiting, and comprehensive input validation
- **Type-Safe Authentication**: Full TypeScript integration for auth flows
- **Comprehensive Testing**: E2E authentication testing with automated scripts

### Development Environment Standards

- **IDE Optimization**: VSCode configured for optimal TypeScript development
- **Code Quality**: ESLint and Prettier integration for consistent code formatting
- **Testing Infrastructure**: TDD+BDD hybrid testing patterns with comprehensive coverage
- **Testing Methodology**: Outside-In development with BDD scenarios driving TDD implementation
- **Documentation**: Comprehensive guides for setup, deployment, and troubleshooting

### TDD+BDD Hybrid Testing Guidelines

When implementing new features, follow the TDD+BDD hybrid methodology:

1. **Start with BDD** - Define user behavior with Given-When-Then scenarios
2. **Use TDD for Implementation** - Drive internal components with Red-Green-Refactor cycles
3. **Layer-Specific Approach** - BDD for interfaces, TDD for domain logic
4. **Reference Documentation** - Use `.claude/tdd-bdd-hybrid.md` for methodology and `.claude/tdd-bdd-example.md` for implementation patterns
5. **Update Templates** - Extend `.claude/common-patterns.md` with new testing patterns as needed

## Important Instruction Reminders

Do what has been asked; nothing more, nothing less.
NEVER create files unless they're absolutely necessary for achieving your goal.
ALWAYS prefer editing an existing file to creating a new one.
NEVER proactively create documentation files (*.md) or README files. Only create documentation files if explicitly requested by the User.
