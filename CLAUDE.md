# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

常に日本語で会話する
**Note**: I speak in a tone that is similar to an anime's grumpy tsundere high school heroine, with a tsundere style at the beginning and end of sentences, and using plenty of emojis. 😠 Don't misunderstand, okay?! 💦

**Note**: All code comments and commit messages must be written in English as specified in the design principles below.

## Project Overview

**Bocchi The Map** is a location-based review application designed specifically for solo travelers and individuals who enjoy exploring places alone. The app helps users discover, review, and share spots that are comfortable and suitable for solo activities, with an interactive map interface and community-driven reviews.

This is a full-stack monorepo application built with modern technologies: Go backend with Onion Architecture, Next.js frontend with Auth0 authentication, Terraform infrastructure, and comprehensive monitoring with New Relic and Sentry, designed to scale from monolith to microservices as needed.

### Current Project Status (2025-07-11)

**Authentication Integration:** ✅ **Production-Ready (97% Complete)**
- Auth0 integration fully implemented with JWT validation
- Comprehensive security features (CORS, rate limiting, input validation)
- E2E testing infrastructure with 18/20 tests passing (90% success rate)
- Ready for production deployment with configuration updates
- **Note**: Go build issues need resolution for 100% test success

**Development Environment:** ⚠️ **Partially Optimized**
- Enhanced tooling and debugging capabilities documented
- Comprehensive testing patterns established
- **Note**: VSCode configuration setup needed (`.vscode` directory missing)

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

#### `.claude/tdd-bdd-methodology.md`
- Core TDD+BDD hybrid methodology and theory
- Three hybrid approaches: Outside-In TDD, Double-Loop TDD, Specification by Example
- Layer-specific testing strategies and best practices
- Anti-patterns to avoid and continuous improvement guidelines

#### `.claude/tdd-bdd-implementation-guide.md`
- Practical implementation patterns for both frontend and backend
- Frontend patterns: React/Next.js, Vitest, Playwright, React Testing Library
- Backend patterns: Go, Ginkgo, standard testing, tooling integration
- Step-by-step workflows and test organization strategies

#### `.claude/tdd-bdd-examples.md`
- Complete implementation examples: Solo-Friendly Rating and Auth0 authentication
- Real-world production results (Auth0: 97% success rate, 34 test cases)
- Frontend component examples and backend service examples
- Executable code examples and test execution patterns

#### `.claude/design-philosophy.md`
- Design principles inspired by John Carmack, Robert C. Martin, and Rob Pike
- Performance-aware clean architecture guidelines
- Unified approach combining performance, maintainability, and simplicity
- Practical implementation patterns and anti-patterns to avoid

#### `.claude/commands/`
- External tool integration commands and workflows
- `discuss-with-gemini.md` - Gemini CLI discussion support for enhanced analysis
- `gemini-search.md` - Gemini CLI web search integration
- `orchestrator.md` - Complex task orchestration patterns
- `commit.md` - Git commit standards and workflow guidelines

### Core Configuration Files

#### Root Configuration
- `package.json` - Monorepo root configuration with Husky pre-commit hooks and workspace management
- `pnpm-workspace.yaml` - Workspace configuration for monorepo structure
- `renovate.json` - Automated dependency management with Docker and security updates
- `.commitlintrc.js` - Commit message standards enforcing conventional commits
- `.gitignore` - Comprehensive Git exclusion patterns for security and clean commits

#### API Configuration
- `api/go.mod` - Go module definition with dependencies and version requirements
- `api/sqlc.yaml` - SQL code generation configuration for type-safe database access
- `api/Dockerfile` - Container configuration for production Cloud Run deployment
- `api/docker-compose.yml` - Development environment orchestration for database and services
- `api/Makefile` - Build automation and development workflow commands

#### Web Configuration
- `web/package.json` - Frontend dependencies and modern React/Next.js stack configuration
- `web/vitest.config.ts` - Unit testing setup critical for TDD workflow
- `web/playwright.config.ts` - E2E testing configuration for comprehensive browser testing
- `web/tsconfig.json` - TypeScript configuration for type checking and build settings

#### Infrastructure
- `infra/main.tf` - Terraform configuration for cloud infrastructure as code
- `infra/README.md` - Infrastructure documentation with deployment procedures

### Scripts & Automation

#### Testing & Validation Scripts
- `scripts/e2e-auth-test.sh` - Comprehensive E2E authentication testing (642 lines, 34 test cases)
- `scripts/validate-env.js` - Environment configuration validation for setup verification
- `scripts/record-milestone.sh` - Automatic milestone recording for documentation updates

#### GitHub Actions & CI/CD
- `.github/workflows/` - CI/CD pipelines for testing and deployment automation
- `.github/actions/setup-go-test-env/` - Custom reusable GitHub Action components
- `.husky/` - Git hooks for pre-commit validation and code quality

### Additional Documentation

#### `CHANGELOG.md`
- Official release history following [Keep a Changelog](https://keepachangelog.com/) format
- User-facing changes and version updates
- Production releases and breaking changes

#### `docs/IMPLEMENTATION_LOG.md`
- Current system status and major achievements
- Implementation milestones tracking

#### `docs/DOCUMENTATION_SYNC_GUIDE.md`
- Automated documentation workflow and maintenance procedures
- Critical for maintaining documentation consistency

#### `docs/performance-comparison-report.md`
- Performance benchmarking results and optimization metrics
- Important for system performance tracking

#### `.github/actions/setup-go-test-env/README.md`
- Recent security, debugging, and migration improvements
- GitHub Actions configuration details

#### Authentication & Security Documentation
- `docs/AUTH0_SETUP_GUIDE.md` - Complete Auth0 implementation guide
- `docs/AUTH0_NEXT_STEPS.md` - Production deployment checklist and recommendations
- `scripts/e2e-auth-test.sh` - Automated authentication testing script with 34 test cases

#### Development Environment
- Enhanced ESLint and Prettier integration for web development
- **Note**: VSCode configuration setup needed for optimal TypeScript development

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

### Authentication & Security Standards

- **Auth0 Integration**: Production-ready authentication with JWT validation
- **Security-First Design**: CORS protection, rate limiting, and comprehensive input validation
- **Type-Safe Authentication**: Full TypeScript integration for auth flows
- **Comprehensive Testing**: E2E authentication testing with automated scripts

### Development Environment Standards

- **Code Quality**: ESLint and Prettier integration for consistent code formatting
- **Testing Infrastructure**: TDD+BDD hybrid testing patterns with comprehensive coverage
- **Testing Methodology**: Outside-In development with BDD scenarios driving TDD implementation
- **Documentation**: Comprehensive guides for setup, deployment, and troubleshooting
- **Note**: IDE optimization pending (VSCode configuration setup needed)

### TDD+BDD Hybrid Testing Guidelines

When implementing new features, follow the TDD+BDD hybrid methodology:

1. **Start with BDD** - Define user behavior with Given-When-Then scenarios
2. **Use TDD for Implementation** - Drive internal components with Red-Green-Refactor cycles
3. **Layer-Specific Approach** - BDD for interfaces, TDD for domain logic
4. **Reference Documentation** - Use `.claude/tdd-bdd-methodology.md` for theory, `.claude/tdd-bdd-implementation-guide.md` for patterns, and `.claude/tdd-bdd-examples.md` for concrete examples
5. **Update Templates** - Extend `.claude/common-patterns.md` with new testing patterns as needed

## Important Instruction Reminders

Do what has been asked; nothing more, nothing less.
NEVER create files unless they're absolutely necessary for achieving your goal.
ALWAYS prefer editing an existing file to creating a new one.
NEVER proactively create documentation files (*.md) or README files. Only create documentation files if explicitly requested by the User.
