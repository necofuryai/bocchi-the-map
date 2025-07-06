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
- `Auth Provider` - Auth0 Universal Login session management
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

#### Authentication Middleware Pattern

```go
// Enhanced authentication middleware with token blacklisting
type AuthMiddleware struct {
    jwtValidator *JWTValidator
    queries      *database.Queries
    rateLimiter  *RateLimiter
    logger       zerolog.Logger
}

// Middleware function with comprehensive security checks
func (am *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Extract and validate JWT token
        token := extractTokenFromRequest(r)
        claims, err := am.jwtValidator.ValidateToken(r.Context(), token)
        if err != nil {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }
        
        // 2. Check token blacklist
        if err := am.checkTokenBlacklist(r.Context(), claims.ID); err != nil {
            http.Error(w, "token has been revoked", http.StatusUnauthorized)
            return
        }
        
        // 3. Add authentication context
        ctx := am.addAuthContext(r.Context(), claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Safe fail pattern for blacklist checking
func (am *AuthMiddleware) checkTokenBlacklist(ctx context.Context, jti string) error {
    if jti == "" {
        return nil // No JTI means older token format, allow for backward compatibility
    }
    
    isBlacklisted, err := am.queries.IsTokenBlacklisted(ctx, jti)
    if err != nil {
        // Log error but don't fail authentication - availability over security
        am.logger.Warn("Blacklist check failed, allowing authentication", err)
        return nil
    }
    
    if isBlacklisted {
        return errors.New("token has been revoked")
    }
    return nil
}
```

#### Secure Account Deletion Pattern

```go
// Multi-step secure deletion with rollback capability
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
    // 1. Verify user exists and get associated data
    user, err := s.queries.GetUser(ctx, userID)
    if err != nil {
        return status.Error(codes.NotFound, "user not found")
    }
    
    // 2. Begin transaction for atomic deletion
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return status.Error(codes.Internal, "failed to begin transaction")
    }
    defer tx.Rollback()
    
    // 3. Delete user (CASCADE will handle related data)
    if err := s.queries.WithTx(tx).DeleteUser(ctx, userID); err != nil {
        return status.Error(codes.Internal, "failed to delete user")
    }
    
    // 4. Commit transaction
    if err := tx.Commit(); err != nil {
        return status.Error(codes.Internal, "failed to commit deletion")
    }
    
    // 5. Log deletion for audit trail
    s.logger.InfoWithFields("User account deleted", map[string]interface{}{
        "user_id": userID,
        "email":   user.Email,
    })
    
    return nil
}
```

#### Context Management Pattern

```go
// Type-safe context key management
type contextKey string

const (
    userIDKey    contextKey = "user_id"
    jtiKey       contextKey = "jti"
    tokenExpKey  contextKey = "token_exp"
)

// Context helper functions with type safety
func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

func GetUserIDFromContext(ctx context.Context) string {
    if userID, ok := ctx.Value(userIDKey).(string); ok {
        return userID
    }
    return ""
}

func WithJTI(ctx context.Context, jti string) context.Context {
    return context.WithValue(ctx, jtiKey, jti)
}

func GetJTIFromContext(ctx context.Context) string {
    if jti, ok := ctx.Value(jtiKey).(string); ok {
        return jti
    }
    return ""
}
```

### Frontend Patterns

#### TDD+BDD Component Pattern

```typescript
// BDD-style component testing
describe('SearchInput Component', () => {
  describe('Given the SearchInput component is rendered', () => {
    describe('When user types in the search input', () => {
      it('Then the input value should update correctly', async () => {
        // Given
        const user = userEvent.setup()
        render(<SearchInput onSearch={mockOnSearch} onClear={mockOnClear} />)
        
        // When
        await user.type(searchInput, 'quiet cafe')
        
        // Then
        expect(searchInput).toHaveValue('quiet cafe')
      })
    })
  })
})
```

#### Custom Hook TDD Pattern

```typescript
// TDD-driven custom hook development
export function useSpotSearch(): UseSpotSearchReturn {
  const [spots, setSpots] = useState<Spot[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  
  const search = useCallback(async (query: string) => {
    // Implementation driven by tests
    setLoading(true)
    try {
      const response = await searchSpots(query, filters)
      setSpots(response.data)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }, [filters])
  
  return { spots, loading, error, search }
}
```

#### BDD E2E Test Pattern

```typescript
// Playwright E2E tests with BDD structure
test.describe('Spot Search Feature', () => {
  test.describe('Given I am on the search page', () => {
    test('When I search for "quiet cafe", Then I should see relevant results', async ({ page }) => {
      // Given
      await page.goto('/search')
      
      // When
      await page.getByTestId('search-input').fill('quiet cafe')
      await page.keyboard.press('Enter')
      
      // Then
      await expect(page.getByTestId('search-results')).toBeVisible()
    })
  })
})
```

#### Component Composition Pattern

```typescript
// Composable UI components with testability
export function SpotCard({ spot, onSelect }: SpotCardProps) {
  return (
    <Card onClick={() => onSelect(spot)} data-testid="spot-item">
      <CardHeader>
        <CardTitle data-testid="spot-name">{spot.name}</CardTitle>
        {spot.soloFriendly && (
          <Badge data-testid="solo-friendly-badge">Solo-friendly</Badge>
        )}
      </CardHeader>
    </Card>
  );
}
```

#### Test Utilities Pattern

```typescript
// BDD-style test helpers and assertions
export const BDDAssertions = {
  expectLoadingState: () => {
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument()
  },
  
  expectSearchResults: (count?: number) => {
    const resultsContainer = screen.getByTestId('search-results')
    expect(resultsContainer).toBeInTheDocument()
    
    if (count !== undefined) {
      const spotItems = screen.getAllByTestId('spot-item')
      expect(spotItems).toHaveLength(count)
    }
  }
}

export const BDDActions = {
  performSearch: async (query: string) => {
    const searchInput = screen.getByPlaceholderText('Search for spots...')
    await userEvent.clear(searchInput)
    await userEvent.type(searchInput, query)
    await userEvent.keyboard('{Enter}')
  }
}
```

#### MSW API Mocking Pattern

```typescript
// Realistic API mocking for tests
export const spotHandlers = [
  http.get('/api/spots/search', ({ request }) => {
    const url = new URL(request.url)
    const query = url.searchParams.get('q')?.toLowerCase() || ''
    
    const filteredSpots = mockSpots.filter(spot =>
      spot.name.toLowerCase().includes(query)
    )
    
    return HttpResponse.json({
      data: filteredSpots,
      total: filteredSpots.length,
      hasMore: false,
    })
  }),
]
```

## Anti-Patterns to Avoid

### Backend Anti-Patterns

1. **Direct Database Access in Handlers**: Always use the repository pattern
2. **Business Logic in Infrastructure**: Keep domain logic in the domain layer
3. **Tight Coupling Between Layers**: Respect onion architecture dependencies
4. **Inconsistent Error Handling**: Use structured error types consistently

#### Authentication & Security Anti-Patterns

5. **Storing JTI in Database Without Checking**: 
   - ❌ Don't add tokens to blacklist without verifying they aren't already there
   - ✅ Use proper unique constraints and handle duplicate key errors gracefully
   
6. **Failing Authentication on Database Errors**:
   - ❌ Don't block all authentication when blacklist database is unavailable
   - ✅ Implement graceful degradation - log errors but allow authentication
   
7. **Missing Token Context Information**:
   - ❌ Don't store only user ID in request context
   - ✅ Store JTI, expiration, and other token metadata for security operations
   
8. **Inconsistent Account Deletion Scope**:
   - ❌ Don't leave orphaned data when deleting user accounts
   - ✅ Use CASCADE constraints and verify all related data is cleaned up
   
9. **Hard Failure on Token Blacklisting**:
   - ❌ Don't fail logout/deletion operations if blacklisting fails
   - ✅ Log blacklist failures but complete the primary operation
   
10. **Missing Audit Trail for Security Operations**:
    - ❌ Don't perform account deletion without logging
    - ✅ Always log security-sensitive operations with sufficient context

#### Context Management Anti-Patterns

11. **String-based Context Keys**:
    - ❌ `ctx.Value("user_id")`
    - ✅ Use typed context keys: `ctx.Value(userIDKey)`
    
12. **Missing Context Value Type Assertions**:
    - ❌ `userID := ctx.Value(userIDKey).(string)` (can panic)
    - ✅ `userID, ok := ctx.Value(userIDKey).(string); if !ok { return "" }`
    
13. **Overloading Context with Non-Request Data**:
    - ❌ Don't store application configuration in request context
    - ✅ Only store request-scoped authentication and user data

### Frontend Anti-Patterns

#### React Development Anti-Patterns
1. **Direct State Mutation**: Use immutable updates with React state
2. **Prop Drilling**: Use context or state management for deep component trees
3. **Uncontrolled Side Effects**: Use useEffect dependencies properly
4. **Missing Error Boundaries**: Implement error handling for async operations

#### Frontend Testing Anti-Patterns
1. **Testing Implementation Details**: 
   - ❌ Don't test internal component state
   - ✅ Test user-visible behavior and interactions
   
2. **Over-mocking in E2E Tests**: 
   - ❌ Don't mock internal components in E2E tests
   - ✅ Mock only external APIs and services
   
3. **Skipping TDD Cycles**: 
   - ❌ Don't write component code before tests
   - ✅ Always follow Red-Green-Refactor cycle
   
4. **Missing Accessibility Testing**:
   - ❌ Don't ignore screen reader compatibility
   - ✅ Always test with semantic queries and a11y tools
   
5. **Writing Tests After Implementation**:
   - ❌ Don't write tests just to increase coverage
   - ✅ Let tests drive your component design (TDD)
   
6. **Mixing Test Concerns**:
   - ❌ Don't test BDD scenarios in unit tests
   - ✅ Keep E2E (BDD) and component (TDD) tests separate
   
7. **Unrealistic Test Data**:
   - ❌ Don't use overly simplified mock data
   - ✅ Use realistic API responses with MSW

## Testing Strategies

### TDD+BDD Hybrid Methodology

This project employs a sophisticated hybrid testing approach that combines the strengths of both Test-Driven Development (TDD) and Behavior-Driven Development (BDD):

#### Core Philosophy
- **Outside-In Development**: Start with BDD scenarios to define user behavior, then use TDD for implementation
- **Double-Loop Testing**: Outer loop (BDD) drives user stories, inner loop (TDD) drives implementation details
- **Specification by Example**: Use concrete examples to drive both behavior and implementation

#### Layer-Specific Testing Strategy

**Interface Layer (Handlers/Controllers)**
- Primary: BDD approach with Ginkgo
- Focus: User interactions and API contracts
- Tools: E2E tests with Given-When-Then structure

**Application Layer (Services)**
- Primary: TDD with BDD context
- Focus: Business logic orchestration
- Tools: Standard Go testing with clear scenarios

**Domain Layer (Core Business Logic)**
- Primary: Pure TDD
- Focus: Business rules and entities
- Tools: Table-driven tests, property-based testing

**Infrastructure Layer (Adapters)**
- Primary: TDD with integration tests
- Focus: External system interactions
- Tools: Mocks and test containers

#### Test Flow Pattern
1. **Feature Request** → BDD scenario definition
2. **E2E Test** → High-level behavior specification
3. **TDD Cycle** → Implementation of supporting components
4. **Integration** → Ensure all layers work together
5. **Validation** → BDD scenarios pass end-to-end

#### Benefits
- **User-Centric**: BDD ensures features meet user needs
- **Clean Implementation**: TDD ensures robust, testable code
- **Comprehensive Coverage**: Both behavior and implementation are tested
- **Maintainable**: Clear separation of concerns in test structure

#### Security Feature Implementation Pattern

When implementing security-critical features like authentication, token management, and account deletion, follow this specialized TDD+BDD approach:

##### BDD Security Scenarios
```gherkin
Feature: Token Blacklist Management
  Scenario: User logs out and token is blacklisted
    Given a user is authenticated with a valid JWT
    When the user logs out
    Then the token should be added to blacklist
    And subsequent requests with that token should be rejected

Feature: Account Deletion Security
  Scenario: User deletes account with data cleanup
    Given a user is authenticated
    When the user requests account deletion
    Then the user data should be removed from database
    And all user tokens should be blacklisted
    And related content should be cleaned up via CASCADE
```

##### TDD Security Implementation
```go
// 1. Red Phase - Write failing test
func TestAuthMiddleware_CheckTokenBlacklist(t *testing.T) {
    tests := []struct {
        name          string
        jti           string
        setupMock     func(*mocks.MockQuerier)
        expectedError bool
        errorMessage  string
    }{
        {
            name: "revoked token should be rejected",
            jti:  "test-jti-123",
            setupMock: func(m *mocks.MockQuerier) {
                m.EXPECT().IsTokenBlacklisted(gomock.Any(), "test-jti-123").Return(true, nil)
            },
            expectedError: true,
            errorMessage:  "token has been revoked",
        },
        {
            name: "database error should not block authentication",
            jti:  "test-jti-456",
            setupMock: func(m *mocks.MockQuerier) {
                m.EXPECT().IsTokenBlacklisted(gomock.Any(), "test-jti-456").Return(false, errors.New("db error"))
            },
            expectedError: false, // Graceful degradation
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}

// 2. Green Phase - Implement minimal code to pass
func (am *AuthMiddleware) checkTokenBlacklist(ctx context.Context, jti string) error {
    // Implementation that makes tests pass
}

// 3. Refactor Phase - Optimize and clean up
func (am *AuthMiddleware) checkTokenBlacklist(ctx context.Context, jti string) error {
    // Refactored implementation with proper error handling
}
```

##### E2E Security Validation
```go
// BDD E2E test using Ginkgo
var _ = Describe("Token Blacklist Security", func() {
    Context("Given a user is authenticated", func() {
        BeforeEach(func() {
            // Setup authenticated user context
            user := authHelper.CreateTestUser()
            token = authHelper.GenerateTokenForUser(user)
        })
        
        Context("When the user logs out", func() {
            It("Then subsequent requests should be rejected", func() {
                By("Successfully logging out")
                response := authHelper.Logout(token)
                Expect(response.StatusCode).To(Equal(200))
                
                By("Rejecting subsequent authenticated requests")
                response = authHelper.MakeAuthenticatedRequest(token, "/api/v1/users/me")
                Expect(response.StatusCode).To(Equal(401))
                Expect(response.Body).To(ContainSubstring("token has been revoked"))
            })
        })
    })
})
```

##### Security Testing Best Practices

1. **Always Test Security Boundaries**:
   - Test unauthorized access attempts
   - Test with expired tokens
   - Test with malformed tokens
   - Test account deletion edge cases

2. **Test Graceful Degradation**:
   - Database unavailability scenarios
   - Partial system failures
   - Network timeouts and retries

3. **Audit Trail Verification**:
   - Verify security operations are logged
   - Test log content and format
   - Ensure sensitive data is not logged

4. **Data Consistency Validation**:
   - Verify CASCADE deletions work correctly
   - Test transaction rollback scenarios
   - Validate orphaned data cleanup

### Traditional Testing Approaches

#### Backend Testing

- **Unit Tests**: Test domain entities and services in isolation
- **Integration Tests**: Test repository implementations with real database
- **BDD Tests**: Use Ginkgo for behavior specifications
- **Contract Tests**: Verify Protocol Buffer contracts

#### Frontend Testing

The frontend adopts the same TDD+BDD hybrid methodology as the backend, adapted for React/Next.js development:

**TDD+BDD Frontend Hybrid Approach:**

**Presentation Layer (React Components)**
- Primary: BDD approach for user interactions
- Focus: Component behavior, props, and user events  
- Tools: React Testing Library with user-event
- Pattern: Given-When-Then for component behavior

**Logic Layer (Custom Hooks)**
- Primary: TDD with BDD context
- Focus: Business logic, state management, side effects
- Tools: Vitest with React Testing Library renderHook
- Pattern: Red-Green-Refactor for hook logic

**Integration Layer (API Calls)**  
- Primary: TDD with mocking
- Focus: Data fetching, caching, error handling
- Tools: MSW (Mock Service Worker) for API mocking
- Pattern: Test doubles for external dependencies

**E2E Layer (User Journeys)**
- Primary: Pure BDD
- Focus: Complete user workflows
- Tools: Playwright for full browser automation
- Pattern: User story scenarios

**Frontend Testing Infrastructure:**
- **Test Utilities**: Comprehensive test helpers with BDD-style assertions
- **API Mocking**: MSW integration for realistic API responses
- **Accessibility Testing**: Automated a11y validation with jest-axe
- **Visual Testing**: Component screenshot comparison capabilities

**Example Frontend TDD+BDD Workflow:**
1. Write E2E BDD scenario (e.g., `spot-search.spec.ts`)
2. Implement components with TDD (e.g., `SearchInput`, `useSpotSearch`)
3. Create integration tests for component cooperation
4. Validate E2E scenarios pass end-to-end

**Frontend Test Types:**
- **Unit Tests**: Test utility functions and custom hooks with Vitest
- **Component Tests**: Test React components in isolation with BDD structure
- **Integration Tests**: Test component interactions and data flow
- **E2E Tests**: Test complete user workflows with Playwright
- **Accessibility Tests**: Automated accessibility compliance testing
- **Visual Regression Tests**: UI consistency validation

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
# Build stage - Full Go development environment with multi-arch support
FROM --platform=$BUILDPLATFORM golang:1.24-alpine@sha256:... AS builder

# Build arguments for cross-compilation
ARG TARGETOS
ARG TARGETARCH

# Cross-compilation build step
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o main ./cmd/api

# Production stage - Minimal runtime with target platform
FROM --platform=$TARGETPLATFORM alpine:3.19@sha256:...
# Security: non-root user, ca-certificates, health checks
```

**Multi-Architecture Build Options**:
- `--platform=$BUILDPLATFORM`: Uses the native platform for the build stage to maximize build performance
- `--platform=$TARGETPLATFORM`: Sets the target platform for the final runtime image
- `TARGETOS`/`TARGETARCH`: Build arguments automatically set by Docker buildx for Go cross-compilation
- Cross-compilation is handled natively by Go with `GOOS` and `GOARCH` environment variables

**Build Commands**:
```bash
# Single architecture build
docker build -t bocchi-the-map .

# Multi-architecture build with buildx
docker buildx build --platform linux/amd64,linux/arm64 -t bocchi-the-map .
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

#### Core Authentication Features
- **JWT Token Management**: Use JWT tokens with proper expiration and Auth0 validation
- **Token Blacklisting**: Implement JWT ID (JTI) based token revocation for logout security
- **Account Management**: Secure account deletion with CASCADE data removal
- **CSRF Protection**: Implement CSRF protection for state-changing operations
- **Input Validation**: Validate all user inputs at multiple layers
- **HTTPS Enforcement**: Use HTTPS everywhere in production

#### Advanced Security Patterns

##### Token Blacklist Implementation
```go
// JWT ID (JTI) extraction and blacklist management
type Claims struct {
    ID          string   `json:"jti,omitempty"`     // JWT ID for blacklist tracking
    Audience    []string `json:"aud,omitempty"`
    Subject     string   `json:"sub,omitempty"`
    Email       string   `json:"email,omitempty"`
    Permissions []string `json:"permissions,omitempty"`
    jwt.RegisteredClaims
}

// Blacklist check in authentication middleware
func (am *AuthMiddleware) checkTokenBlacklist(ctx context.Context, jti string) error {
    isBlacklisted, err := am.queries.IsTokenBlacklisted(ctx, jti)
    if err != nil {
        // Fail safely - don't block authentication on DB issues
        am.logger.Warn("Failed to check token blacklist", err)
        return nil
    }
    if isBlacklisted {
        return errors.New("token has been revoked")
    }
    return nil
}
```

##### Secure Account Deletion Pattern
```go
// Account deletion with security checks and cascade operations
func (h *UserHandler) DeleteCurrentUser(ctx context.Context, input *DeleteCurrentUserInput) (*DeleteCurrentUserOutput, error) {
    // 1. Extract authenticated user from context
    userID := auth.GetUserIDFromContext(ctx)
    if userID == "" {
        return nil, huma.Error401Unauthorized("authentication required")
    }
    
    // 2. Delete user via gRPC service (triggers CASCADE deletion)
    err := h.userClient.DeleteUser(ctx, userID)
    if err != nil {
        return nil, huma.Error500InternalServerError("failed to delete user account")
    }
    
    // 3. Blacklist current token to prevent further access
    if jti := auth.GetJTIFromContext(ctx); jti != "" {
        if err := h.authMiddleware.Logout(ctx); err != nil {
            h.logger.Warn("Failed to blacklist token during account deletion", err)
        }
    }
    
    return &DeleteCurrentUserOutput{}, nil
}
```

#### Authentication Context Management
```go
// Context keys for authentication data
const (
    UserIDContextKey    = "user_id"
    JTIContextKey      = "jti"
    TokenExpContextKey = "token_exp"
)

// Helper functions for context management
func GetJTIFromContext(ctx context.Context) string {
    if jti, ok := ctx.Value(JTIContextKey).(string); ok {
        return jti
    }
    return ""
}

func GetTokenExpirationFromContext(ctx context.Context) *time.Time {
    if exp, ok := ctx.Value(TokenExpContextKey).(*time.Time); ok {
        return exp
    }
    return nil
}
```

#### Database Security Patterns
```sql
-- Secure user deletion with CASCADE constraints
DELETE FROM users WHERE id = ?;
-- Related reviews are automatically deleted via CASCADE

-- Token blacklist with automatic cleanup
INSERT INTO token_blacklist (jti, token_type, expires_at) 
VALUES (?, 'access', ?);

-- Cleanup expired tokens (run periodically)
DELETE FROM token_blacklist WHERE expires_at < NOW();
```

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

```text
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
