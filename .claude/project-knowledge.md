# Project Knowledge

This document contains implementation patterns, design decisions, architecture choices, and patterns to avoid.

## Development Philosophy

### Behavior-Driven Development (BDD)

Follow Behavior-Driven Development principles throughout the project:

- Start with BDD for all new features and bug fixes
- Write behavior specifications first using Given-When-Then format
- Focus on describing the behavior from user's perspective
- Use Ginkgo framework for Go testing with descriptive specs
- Only write specification code initially, no implementation
- Run specs to verify they fail as expected
- Commit specs once verified correct
- Then implement code to make specs pass
- Never modify specs during implementation - only fix the code
- Repeat until all specs pass

## API Architecture (Onion Architecture)

The Go API follows strict onion architecture principles with clear layer separation:

### Domain Layer (`/domain/`)

- `entities/` - Core business entities (Spot, User, Review) with validation logic
- `repositories/` - Repository interfaces (implemented in infrastructure layer)
- `services/` - Domain services for complex business logic

### Application Layer (`/application/`)

- `clients/` - Application services orchestrating domain entities

### Infrastructure Layer (`/infrastructure/`)

- `database/` - sqlc-generated database models and queries
- `grpc/` - gRPC service implementations (TiDB/MySQL)
- `external/` - Third-party service integrations

### Interface Layer (`/interfaces/`)

- `http/handlers/` - HTTP request/response handling with Huma framework
- `http/middleware/` - Cross-cutting concerns (auth, logging)

### Protocol Buffers (`/proto/`)

- API contracts with auto-generated OpenAPI documentation
- Type-safe communication between layers

## Frontend Architecture (web/)

### Component Structure

- `src/app/` - Next.js 15 App Router pages and layouts
- `src/components/ui/` - Reusable Shadcn/ui components
- `src/components/map/` - MapLibre GL JS integration components
- `src/hooks/` - Custom React hooks for map interactions and state
- `src/lib/` - Utilities and shared configurations
- `src/types/` - TypeScript type definitions

### Key Components

- `Map component` - Main MapLibre GL JS wrapper with PMTiles support
- `POI Features` - Point of interest rendering and interaction logic
- `Auth Provider` - Auth.js session management
- `Theme Provider` - Dark/light mode using next-themes

## Key Design Principles

1. **Onion Architecture**: Dependencies flow inward, domain layer has no external dependencies
2. **Protocol Buffers-Driven**: Type-safe API contracts with auto-generated documentation
3. **Microservice-Ready**: Loose coupling for future service extraction
4. **Type Safety**: Protocol Buffers for API, TypeScript for frontend
5. **Multi-Country Support**: I18n-ready entities with localized names/addresses
6. **Structured Logging**: JSON format with zerolog (ERROR, WARN, INFO, DEBUG)
7. **Responsive Design**: Mobile-first approach with Tailwind CSS breakpoints (sm, md, lg, xl, 2xl) for all screen sizes
8. **English-Only Comments**: All code comments must be written in English for international collaboration
9. **English-Only Commit Messages**: All git commit messages must be written in English for international collaboration

## Fundamental Development Principles

The project strictly adheres to these core software development principles:

### SOLID Principles

#### Single Responsibility Principle (SRP)
- Each class/struct should have only one reason to change
- Handlers only handle HTTP concerns, Services only contain business logic
- Example: `UserHandler` only handles HTTP requests, `UserService` only contains user business logic

#### Open/Closed Principle (OCP)
- Software entities should be open for extension, closed for modification
- Use interfaces and dependency injection to extend functionality
- Example: `SpotRepository` interface allows different storage implementations

#### Liskov Substitution Principle (LSP)
- Objects should be replaceable with instances of their subtypes
- All repository implementations must fulfill their interface contracts
- Example: Memory and database repositories are interchangeable

#### Interface Segregation Principle (ISP)
- Many client-specific interfaces are better than one general-purpose interface
- Create focused interfaces rather than large, monolithic ones
- Example: Separate `SpotReader` and `SpotWriter` instead of one large interface

#### Dependency Inversion Principle (DIP)
- Depend on abstractions, not concretions
- High-level modules should not depend on low-level modules
- Example: Application layer depends on repository interfaces, not database implementations

### KISS (Keep It Simple, Stupid)

- Prefer simple solutions over complex ones
- Avoid over-engineering and premature optimization
- Use clear, readable code over clever tricks
- Example: Straightforward error handling with `if err != nil` instead of complex error wrapping

### YAGNI (You Aren't Gonna Need It)

- Don't implement features until they are actually needed
- Remove unused code and dependencies
- Focus on current requirements, not hypothetical future needs
- Example: Don't create complex caching systems until performance issues are proven

### DRY (Don't Repeat Yourself)

- Every piece of knowledge should have a single, unambiguous representation
- Extract common code into shared utilities
- Use code generation (sqlc, Protocol Buffers) to eliminate repetition
- Example: Common error handling patterns in `pkg/errors/` package

## Implementation Patterns

### Backend Patterns

#### Repository Pattern

```go
// Domain layer - interface
type SpotRepository interface {
    Create(ctx context.Context, spot *entities.Spot) error
    GetByID(ctx context.Context, id string) (*entities.Spot, error)
}

// Infrastructure layer - implementation
type spotRepository struct {
    db *sql.DB
}
```

#### Service Layer Pattern

```go
// Application layer
type SpotClient struct {
    spotRepo domain.SpotRepository
    logger   zerolog.Logger
}

func (c *SpotClient) CreateSpot(ctx context.Context, req *CreateSpotRequest) error {
    // Business logic orchestration
}
```

#### Error Handling Pattern

```go
// Use custom error types with context
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}
```

### Frontend Patterns

#### Custom Hook Pattern

```typescript
// Map interaction hooks
export function useMapInteraction() {
  const [selectedSpot, setSelectedSpot] = useState<Spot | null>(null);
  
  const handleSpotClick = useCallback((spot: Spot) => {
    setSelectedSpot(spot);
  }, []);
  
  return { selectedSpot, handleSpotClick };
}
```

#### Component Composition Pattern

```typescript
// Composable UI components
export function SpotCard({ spot, onSelect }: SpotCardProps) {
  return (
    <Card onClick={() => onSelect(spot)}>
      <CardHeader>
        <CardTitle>{spot.name}</CardTitle>
      </CardHeader>
    </Card>
  );
}
```

## Anti-Patterns to Avoid

### Backend Anti-Patterns

1. **Direct Database Access in Handlers**: Always use the repository pattern
2. **Business Logic in Infrastructure**: Keep domain logic in the domain layer
3. **Tight Coupling Between Layers**: Respect onion architecture dependencies
4. **Inconsistent Error Handling**: Use structured error types consistently

### Frontend Anti-Patterns

1. **Direct State Mutation**: Use immutable updates with React state
2. **Prop Drilling**: Use context or state management for deep component trees
3. **Uncontrolled Side Effects**: Use useEffect dependencies properly
4. **Missing Error Boundaries**: Implement error handling for async operations

## Testing Strategies

### Backend Testing

- **Unit Tests**: Test domain entities and services in isolation
- **Integration Tests**: Test repository implementations with real database
- **BDD Tests**: Use Ginkgo for behavior specifications
- **Contract Tests**: Verify Protocol Buffer contracts

### Frontend Testing

- **Unit Tests**: Test utility functions and custom hooks with Vitest
- **Component Tests**: Test React components in isolation
- **Integration Tests**: Test component interactions and data flow
- **E2E Tests**: Test complete user workflows with Playwright

## Performance Considerations

### Backend Performance

- Use connection pooling for database connections
- Implement proper indexing strategies
- Use context cancellation for timeout handling
- Profile performance with Go's built-in profiler

### Frontend Performance

- Implement virtual scrolling for large lists
- Use React.memo for expensive component renders
- Optimize map rendering with clustering for many points
- Implement proper image lazy loading

## Cloud Run & Monitoring Integration (Latest Implementation)

### Monitoring Architecture

#### New Relic Integration
- **Application Performance Monitoring**: Custom metrics, distributed tracing, and performance insights
- **Middleware Integration**: HTTP request monitoring with transaction tracking
- **Custom Metrics**: Business metrics recording (spot creations, user activities)
- **Background Transactions**: Non-web operation monitoring
- **Graceful Shutdown**: Proper flush handling on application termination

```go
// Initialize monitoring in main()
if err := monitoring.InitMonitoring(
    cfg.Monitoring.NewRelicLicenseKey,
    cfg.Monitoring.SentryDSN,
    "bocchi-the-map-api",
    cfg.App.Environment,
    "1.0.0",
); err != nil {
    logger.Error("Failed to initialize monitoring", err)
    // Don't exit - monitoring is not critical for basic functionality
}
```

#### Sentry Integration
- **Error Tracking**: Context-aware error capturing with breadcrumbs
- **Performance Monitoring**: Transaction tracking and bottleneck identification
- **Release Tracking**: Version-based error attribution
- **User Context**: Request-specific error attribution
- **Sensitive Data Filtering**: Automatic removal of sensitive headers and data

```go
// Unified error handling with Sentry integration
logger.ErrorWithContext(ctx, "Database operation failed", err)
logger.ErrorWithContextAndFields(ctx, "User operation failed", err, map[string]interface{}{
    "user_id": userID,
    "operation": "create_spot",
})
```

### Docker Containerization

#### Multi-Stage Build Pattern
```dockerfile
# Build stage - Full Go development environment
FROM golang:1.21-alpine AS builder
# ... build steps

# Production stage - Minimal runtime
FROM alpine:latest
# Security: non-root user, ca-certificates, health checks
```

#### Security Best Practices
- **Non-root execution**: Application runs as dedicated user (uid 1001)
- **Minimal base image**: Alpine Linux for reduced attack surface
- **Health checks**: Container orchestration support
- **Efficient layering**: Optimized for Docker layer caching

### Terraform Infrastructure

#### Secret Management
```hcl
# Google Secret Manager integration
resource "google_secret_manager_secret" "new_relic_license_key" {
  secret_id = "new-relic-license-key-${var.environment}"
  replication { auto {} }
}

# Service account with minimal permissions
resource "google_service_account" "cloud_run_service_account" {
  account_id   = "bocchi-cloud-run-${var.environment}"
  display_name = "Bocchi Cloud Run Service Account"
}
```

#### Cloud Run Configuration
- **Service Account Integration**: Dedicated IAM with minimal required permissions
- **Environment-Specific Scaling**: Production (min 1, max 10), Development (min 0, max 3)
- **Resource Optimization**: CPU/memory requests and limits based on usage patterns
- **Health Check Integration**: Kubernetes-ready probes

### Configuration Management

#### Environment-Based Pattern
```go
type MonitoringConfig struct {
    NewRelicLicenseKey string // From Secret Manager
    SentryDSN          string // From Secret Manager
}

// Graceful degradation when monitoring unavailable
func (c *Config) Validate() error {
    if c.Database.Password == "" {
        return errors.New("TIDB_PASSWORD is required")
    }
    // Monitoring is optional - application continues without it
    return nil
}
```

### Build and Deployment

#### Automated Build Script
```bash
# Environment-aware build and deployment
./scripts/build.sh dev YOUR_PROJECT_ID asia-northeast1

# Features:
# - Multi-arch Docker builds
# - Automatic image tagging with timestamps
# - Optional Cloud Run deployment
# - GCR authentication handling
# - Environment validation
```

## Security Best Practices

### Authentication & Authorization

- Use JWT tokens with proper expiration
- Implement CSRF protection
- Validate all user inputs
- Use HTTPS everywhere

### Data Protection

- Sanitize database queries (sqlc helps prevent SQL injection)
- Implement rate-limiting
- Log security events appropriately
- Never log sensitive information

### Cloud Security

- **Secret Management**: Use Google Secret Manager, never environment variables for secrets
- **Service Accounts**: Minimal IAM permissions following principle of least privilege
- **Container Security**: Non-root user execution, minimal base images
- **Network Security**: Cloud Run managed HTTPS with automatic certificate management

## Unified gRPC Architecture Implementation (2025-06-27)

### Architecture Overview

The project has been successfully refactored to use a **unified gRPC architecture pattern** throughout all service layers, eliminating the previous mixed approach of direct database access and gRPC calls.

### Key Architectural Changes

#### 1. **Consistent Service Layer Pattern**

All services now follow the same pattern:
```
HTTP Handler → gRPC Client (internal) → gRPC Service → Database Layer
```

**Before (Mixed Pattern):**
- UserHandler: Direct database access via `queries *database.Queries`
- SpotHandler: gRPC client via `spotClient *clients.SpotClient`
- Inconsistent patterns led to maintenance complexity

**After (Unified Pattern):**
- UserHandler: gRPC client via `userClient *clients.UserClient`
- SpotHandler: gRPC client via `spotClient *clients.SpotClient`
- ReviewHandler: gRPC client via `reviewClient *clients.ReviewClient` (planned)
- All handlers follow identical patterns

#### 2. **gRPC Service Implementations**

**UserService Enhanced:**
- Replaced dummy data with real database operations
- Added consistent gRPC request/response types
- Implemented `CreateUserGRPC`, `UpdateUserGRPC`, `GetUserByAuthProviderGRPC`
- Proper error handling with gRPC status codes

**SpotService Enhanced:**
- Replaced dummy data with comprehensive database integration
- Added spot creation, retrieval, location-based search
- Implemented proper coordinate handling and i18n support
- Full CRUD operations with database persistence

**Database Layer Expansion:**
- Created `spots.sql` with comprehensive queries (location search, filtering, pagination)
- Created `reviews.sql` with rating statistics and user/spot associations
- Generated `spots.sql.go` and `reviews.sql.go` with type-safe operations
- Updated `Querier` interface with all new methods

#### 3. **Type-Safe Database Integration**

**Query Pattern:**
```sql
-- name: ListSpotsByLocation :many
SELECT * FROM spots 
WHERE (6371 * acos(cos(radians(?)) * cos(radians(latitude)) * ...)) <= ?
ORDER BY distance
LIMIT ? OFFSET ?
```

**Generated Go Code:**
```go
type ListSpotsByLocationParams struct {
    Latitude  string `json:"latitude"`
    Longitude string `json:"longitude"`
    RadiusKm  string `json:"radius_km"`
    Limit     int32  `json:"limit"`
    Offset    int32  `json:"offset"`
}

func (q *Queries) ListSpotsByLocation(ctx context.Context, arg ListSpotsByLocationParams) ([]Spot, error)
```

#### 4. **Client Pattern Standardization**

**Consistent Client Structure:**
```go
type UserClient struct {
    service *grpcSvc.UserService
    conn    *grpc.ClientConn
}

func NewUserClient(serviceAddr string, db *sql.DB) (*UserClient, error) {
    if serviceAddr == "internal" {
        return &UserClient{
            service: grpcSvc.NewUserService(db),
        }, nil
    }
    // External gRPC connection for microservices
}
```

**Conversion Pattern:**
```go
func (c *UserClient) convertGRPCUserToEntity(grpcUser *grpcSvc.User) *entities.User {
    // Standard conversion from gRPC types to domain entities
}
```

### Benefits Achieved

#### 1. **Architectural Consistency**
- All handlers use identical gRPC client patterns
- Uniform error handling and response formatting
- Predictable code structure across all modules

#### 2. **Microservice Readiness**
- Internal mode: `service := grpcSvc.NewUserService(db)`
- External mode: `conn, err := grpc.Dial("user-service:9090")`
- Zero code changes required for microservice migration

#### 3. **Type Safety**
- Protocol Buffers ensure contract consistency
- sqlc generates type-safe database operations
- Compile-time verification of all service calls

#### 4. **Scalability Preparation**
- Each service can be extracted independently
- Database operations properly abstracted
- Load balancing and service discovery ready

### Implementation Patterns

#### Service Creation Pattern
```go
// main.go - Dependency injection
userClient, err := clients.NewUserClient("internal", db)
spotClient, err := clients.NewSpotClient("internal", db)

// Handler initialization
userHandler := handlers.NewUserHandler(userClient)
spotHandler := handlers.NewSpotHandler(spotClient)
```

#### Request Flow Pattern
```go
// HTTP Request
POST /api/v1/users
↓
// Handler
userHandler.CreateUser(ctx, input)
↓
// Client
userClient.CreateUser(ctx, domainEntity)
↓
// gRPC Service
userService.CreateUserGRPC(ctx, grpcRequest)
↓
// Database
queries.CreateUser(ctx, dbParams)
```

#### Error Handling Pattern
```go
// gRPC Service level
if req.Email == "" {
    return nil, status.Error(codes.InvalidArgument, "email is required")
}

// Client level - convert gRPC errors to domain errors
if err != nil {
    return nil, err // gRPC status errors propagate correctly
}

// Handler level - convert to HTTP errors
if err != nil {
    return nil, huma.Error500InternalServerError("failed to create user")
}
```

### Database Schema Utilization

**Complete Table Coverage:**
- `users` - OAuth authentication, preferences, profile data
- `spots` - Location data with i18n, ratings, geographic indexing
- `reviews` - User reviews with rating aspects, foreign key constraints

**Advanced Query Features:**
- Geographic distance calculations using Haversine formula
- Full-text search with relevance ranking
- Pagination with count optimization
- JSON field handling for i18n and preferences

### Future Microservice Migration Path

**Phase 1: Internal gRPC (Current)**
```go
userClient := NewUserClient("internal", db)
```

**Phase 2: Service Extraction**
```go
userClient := NewUserClient("user-service:9090", nil)
```

**Phase 3: Service Mesh**
```go
userClient := NewUserClient("user-service.default.svc.cluster.local:9090", nil)
```

The unified architecture provides a clear path for horizontal scaling and service decomposition as the application grows.
