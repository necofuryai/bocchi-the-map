# Common Patterns

This document contains frequently used command patterns and standard implementation templates.

## Development Commands

### API Development

**Basic Commands:**
```bash
cd api
make deps              # Install Go dependencies
make sqlc              # Generate type-safe SQL code from queries/
make proto             # Generate protobuf files
make test              # Run test suite
make run               # Run server
make dev               # Run with hot reload (requires air)
make build             # Build binary to bin/api
make clean             # Clean generated files
```

**Database Commands:**
```bash
make migrate-up        # Run database migrations
make migrate-down      # Rollback database migrations
make migrate-create NAME=migration_name  # Create new migration
make docker-up         # Start MySQL development environment
make docker-down       # Stop development environment
make dev-setup         # Complete development setup (MySQL + migrations)
```

**Documentation:**
```bash
make docs              # Generate OpenAPI documentation
```

### Web Development

**Basic Commands:**
```bash
cd web
pnpm install           # Install dependencies (auto-installs Playwright)
pnpm dev               # Development server (with Turbopack)
pnpm build             # Production build
pnpm start             # Start production server
pnpm lint              # ESLint + TypeScript checking
```

**Testing Commands:**
```bash
pnpm test              # Run unit/component tests with Vitest
pnpm test:ui           # Run Vitest with UI mode
pnpm test:coverage     # Run tests with coverage report
pnpm test:e2e          # Run E2E tests with Playwright
pnpm test:e2e:ui       # Run Playwright with UI mode
```

**Note:** React 19 dependency conflicts are generally resolved better with pnpm

### Infrastructure Commands

```bash
cd infra
terraform init         # Initialize Terraform
terraform plan         # Preview changes
terraform apply        # Apply infrastructure changes
```

### Protocol Buffers

```bash
# From api/ directory
make proto             # Generate Go files from .proto definitions
```

## Advanced Development Commands

### Single Test Execution

**Backend (Go):**
```bash
cd api
# Run specific test file
go test -v ./infrastructure/grpc/spot_service_test.go

# Run specific test function
go test -v -run TestSpotService_CreateSpot ./infrastructure/grpc/

# Run specific test with pattern matching
go test -v -run "TestSpotService_.*" ./...

# Run tests with coverage for specific package
go test -v -cover ./infrastructure/grpc/
```

**Frontend (Vitest/Playwright):**
```bash
cd web
# Run specific test file with Vitest
pnpm test src/components/map/Map.test.tsx

# Run specific test pattern
pnpm test --run --reporter=verbose Map

# Run specific E2E test file
pnpm test:e2e tests/auth.spec.ts

# Run specific E2E test by name
pnpm test:e2e --grep "should login with Google"
```

### Debugging and Logging

**Backend Debug Mode:**
```bash
cd api
# Run with debug logging
LOG_LEVEL=DEBUG make run

# Run with trace logging (most verbose)
LOG_LEVEL=TRACE make run

# Run individual test with verbose output
go test -v -run TestSpotService_CreateSpot ./infrastructure/grpc/ -test.v
```

**Frontend Debug Mode:**
```bash
cd web
# Run development server with debug info
DEBUG=* pnpm dev

# Run tests with debug output
DEBUG=vitest* pnpm test

# Run E2E tests with debug mode
pnpm test:e2e --debug
```

### Performance and Monitoring

**Backend Performance:**
```bash
cd api
# Run with CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./...

# Memory profiling
go test -memprofile=mem.prof -bench=. ./...

# Benchmark specific functions
go test -bench=BenchmarkSpotService ./infrastructure/grpc/
```

**Database Operations:**
```bash
cd api
# Show database connection status
make docker-logs

# Reset database (caution: deletes all data)
make migrate-down && make migrate-up

# Create and run specific migration
make migrate-create NAME=add_user_preferences
make migrate-up
```

### Environment Management

**Environment Variables:**
```bash
# Backend environment setup
cd api
cp .env.example .env
# Edit .env with your configurations

# Frontend environment setup  
cd web
cp .env.local.example .env.local
# Add OAuth credentials to .env.local
```

**Multi-Environment Testing:**
```bash
# Test against local MySQL
cd api
export DATABASE_URL="mysql://bocchi_user:change_me_too@tcp(localhost:3306)/bocchi_the_map"
make test

# Test against TiDB (production-like)
export DATABASE_URL="mysql://user:pass@tcp(gateway.tidbcloud.com:4000)/bocchi_the_map"
make test
```

### Code Generation and Build

**Protocol Buffers:**
```bash
cd api
# Generate only specific proto file
protoc -I proto --go_out=gen proto/spot.proto

# Validate proto files
protoc --proto_path=proto --lint_out=. proto/*.proto
```

**SQL Code Generation:**
```bash
cd api  
# Regenerate after schema changes
make sqlc

# Validate SQL queries
sqlc vet
```

## Troubleshooting Common Issues

### Port Conflicts
```bash
# Check what's using port 8080 (API)
lsof -i :8080

# Check what's using port 3000 (Frontend)
lsof -i :3000

# Kill process using specific port (try graceful termination first)
kill -15 $(lsof -t -i:8080)
# If process doesn't stop, force kill with:
# kill -9 $(lsof -t -i:8080)
```

### Cache Issues
```bash
# Clear Next.js cache
cd web
rm -rf .next/

# Clear Go module cache
cd api
go clean -modcache

# Clear pnpm cache
cd web
pnpm store prune
```

### Dependency Issues
```bash
# Rebuild Go modules
cd api
rm go.sum && go mod tidy

# Reinstall Node modules
cd web
rm -rf node_modules package-lock.json pnpm-lock.yaml
pnpm install
```

## Standard Implementation Templates

### Backend Templates

#### Repository Interface Template
```go
type EntityRepository interface {
    Create(ctx context.Context, entity *entities.Entity) error
    GetByID(ctx context.Context, id string) (*entities.Entity, error)
    Update(ctx context.Context, entity *entities.Entity) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter *ListFilter) ([]*entities.Entity, error)
}
```

#### Service Client Template
```go
type EntityClient struct {
    repo   domain.EntityRepository
    logger zerolog.Logger
}

func NewEntityClient(repo domain.EntityRepository, logger zerolog.Logger) *EntityClient {
    return &EntityClient{
        repo:   repo,
        logger: logger,
    }
}

func (c *EntityClient) CreateEntity(ctx context.Context, req *CreateEntityRequest) (*CreateEntityResponse, error) {
    // Input validation
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Business logic
    entity := &entities.Entity{
        // Map from request
    }
```

## Production Authentication Security Patterns

### JWT Token Management

```go
// Generate JWT with unique ID for tracking
func GenerateToken(userID, email string) (string, error) {
    jti := uuid.New().String()
    claims := &JWTClaims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ID:        jti,  // Essential for token revocation
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "bocchi-the-map-api",
            Subject:   userID,
        },
    }
    return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}
```

### Secure Cookie Configuration

```go
// Production-ready httpOnly cookie settings
func createSecureCookies(accessToken, refreshToken string, expiresAt time.Time) []string {
    isProduction := os.Getenv("ENVIRONMENT") == "production"
    domain := os.Getenv("COOKIE_DOMAIN")
    
    accessCookie := &http.Cookie{
        Name:     "bocchi_access_token",
        Value:    accessToken,
        Expires:  expiresAt,
        HttpOnly: true,           // XSS protection
        Secure:   isProduction,   // HTTPS only in production
        SameSite: http.SameSiteStrictMode, // CSRF protection
        Domain:   domain,
        Path:     "/",
    }
    // Return cookie strings for response headers
    return []string{accessCookie.String()}
}
```

### Token Blacklist Integration

```go
// Check token blacklist in middleware
func (am *AuthMiddleware) validateToken(ctx context.Context, tokenString string) (*JWTClaims, error) {
    claims, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return am.jwtSecret, nil
    })
    
    // Check if token is blacklisted
    if am.queries != nil && claims.ID != "" {
        isBlacklisted, err := am.queries.IsTokenBlacklisted(ctx, claims.ID)
        if err != nil || isBlacklisted {
            return nil, errors.New("token has been revoked")
        }
    }
    return claims, nil
}
```

### Rate Limiting Pattern

```go
// Memory-efficient rate limiter
type RateLimiter struct {
    requests map[string][]time.Time
    mutex    sync.RWMutex
    limit    int           // Max requests per window
    window   time.Duration
}

func (rl *RateLimiter) Allow(clientIP string) bool {
    rl.mutex.Lock()
    defer rl.mutex.Unlock()
    
    now := time.Now()
    cutoff := now.Add(-rl.window)
    
    // Clean old requests
    requests := rl.requests[clientIP]
    var validRequests []time.Time
    for _, req := range requests {
        if req.After(cutoff) {
            validRequests = append(validRequests, req)
        }
    }
    
    if len(validRequests) >= rl.limit {
        return false
    }
    
    validRequests = append(validRequests, now)
    rl.requests[clientIP] = validRequests
    return true
}
```

### Database Cleanup Commands

```bash
# Token blacklist maintenance
cd api

# Manual cleanup of expired tokens
mysql -e "DELETE FROM token_blacklist WHERE expires_at < NOW() - INTERVAL 24 HOUR;"

# Check blacklist table status  
mysql -e "SELECT COUNT(*) as total_tokens, COUNT(CASE WHEN expires_at > NOW() THEN 1 END) as active_tokens FROM token_blacklist;"

# Enable MySQL event scheduler for automatic cleanup
mysql -e "SET GLOBAL event_scheduler = ON;"

# Create cleanup event (from scripts/token_cleanup_event.sql)
mysql < scripts/token_cleanup_event.sql
```

### Security Monitoring Commands

```bash
cd api

# Check authentication endpoint performance
curl -w "@curl-format.txt" -o /dev/null -s "http://localhost:8080/api/v1/auth/token"

# Test rate limiting (should return 429 after 5 requests)
for i in {1..6}; do curl -I "http://localhost:8080/api/v1/auth/token"; done

# Verify CORS headers
curl -H "Origin: http://localhost:3000" -H "Access-Control-Request-Method: POST" -H "Access-Control-Request-Headers: X-Requested-With" -X OPTIONS "http://localhost:8080/api/v1/auth/token"
```

    return &CreateEntityResponse{
        // Map from entity
    }, nil
}
```

#### HTTP Handler Template
```go
type EntityHandler struct {
    client *application.EntityClient
    logger zerolog.Logger
}

func NewEntityHandler(client *application.EntityClient, logger zerolog.Logger) *EntityHandler {
    return &EntityHandler{
        client: client,
        logger: logger,
    }
}

func (h *EntityHandler) CreateEntity(ctx context.Context, req *CreateEntityRequest) (*CreateEntityResponse, error) {
    return h.client.CreateEntity(ctx, req)
}
```

### Frontend Templates

#### Custom Hook Template
```typescript
export function useEntityManagement() {
    const [entities, setEntities] = useState<Entity[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const createEntity = useCallback(async (data: CreateEntityData) => {
        setLoading(true);
        setError(null);
        
        try {
            const response = await api.post('/entities', data);
            setEntities(prev => [...prev, response.data]);
            return response.data;
        } catch (err) {
            const message = err instanceof Error ? err.message : 'Unknown error';
            setError(message);
            throw err;
        } finally {
            setLoading(false);
        }
    }, []);

    const fetchEntities = useCallback(async () => {
        setLoading(true);
        setError(null);
        
        try {
            const response = await api.get('/entities');
            setEntities(response.data);
        } catch (err) {
            const message = err instanceof Error ? err.message : 'Unknown error';
            setError(message);
        } finally {
            setLoading(false);
        }
    }, []);

    return {
        entities,
        loading,
        error,
        createEntity,
        fetchEntities,
    };
}
```

#### Component Template
```typescript
interface EntityListProps {
    entities: Entity[];
    onSelect?: (entity: Entity) => void;
    loading?: boolean;
    error?: string | null;
}

export function EntityList({ entities, onSelect, loading, error }: EntityListProps) {
    if (loading) {
        return <div className="flex justify-center p-4">Loading...</div>;
    }

    if (error) {
        return (
            <div className="text-red-500 p-4">
                Error: {error}
            </div>
        );
    }

    if (entities.length === 0) {
        return (
            <div className="text-gray-500 p-4 text-center">
                No entities found
            </div>
        );
    }

    return (
        <div className="space-y-2">
            {entities.map((entity) => (
                <EntityCard
                    key={entity.id}
                    entity={entity}
                    onClick={() => onSelect?.(entity)}
                />
            ))}
        </div>
    );
}
```

### Testing Templates

#### Backend Test Template (Ginkgo)
```go
var _ = Describe("EntityService", func() {
    var (
        service    *EntityService
        mockRepo   *MockEntityRepository
        ctx        context.Context
    )

    BeforeEach(func() {
        mockRepo = NewMockEntityRepository()
        service = NewEntityService(mockRepo)
        ctx = context.Background()
    })

    Describe("CreateEntity", func() {
        Context("when given valid input", func() {
            It("should create entity successfully", func() {
                // Given
                request := &CreateEntityRequest{
                    Name: "Test Entity",
                }
                
                // When
                result, err := service.CreateEntity(ctx, request)
                
                // Then
                Expect(err).ToNot(HaveOccurred())
                Expect(result).ToNot(BeNil())
                Expect(result.Name).To(Equal("Test Entity"))
            })
        })

        Context("when given invalid input", func() {
            It("should return validation error", func() {
                // Given
                request := &CreateEntityRequest{
                    Name: "", // Invalid empty name
                }
                
                // When
                result, err := service.CreateEntity(ctx, request)
                
                // Then
                Expect(err).To(HaveOccurred())
                Expect(result).To(BeNil())
            })
        })
    })
})
```

#### Frontend Test Template (Vitest)
```typescript
describe('useEntityManagement', () => {
    it('should create entity successfully', async () => {
        // Given
        const mockApi = vi.mocked(api);
        const mockEntity = { id: '1', name: 'Test Entity' };
        mockApi.post.mockResolvedValueOnce({ data: mockEntity });

        const { result } = renderHook(() => useEntityManagement());

        // When
        await act(async () => {
            await result.current.createEntity({ name: 'Test Entity' });
        });

        // Then
        expect(result.current.entities).toContain(mockEntity);
        expect(result.current.error).toBeNull();
        expect(mockApi.post).toHaveBeenCalledWith('/entities', { name: 'Test Entity' });
    });

    it('should handle create entity error', async () => {
        // Given
        const mockApi = vi.mocked(api);
        const error = new Error('Creation failed');
        mockApi.post.mockRejectedValueOnce(error);

        const { result } = renderHook(() => useEntityManagement());

        // When
        let thrownError;
        await act(async () => {
            try {
                await result.current.createEntity({ name: 'Test Entity' });
            } catch (err) {
                thrownError = err;
            }
        });

        // Then
        expect(thrownError).toBe(error);
        expect(result.current.error).toBe('Creation failed');
        expect(result.current.entities).toHaveLength(0);
    });
});

## File Formatting Rules

### Markdown File End-of-File Requirements

**Rule**: All markdown files must end with exactly one blank line.

**Rationale**: 
- Ensures consistent file formatting across the project
- Prevents Git diff issues related to "No newline at end of file" warnings
- Maintains compatibility with POSIX standards

**Target Files**:
- All README.md files in project directories
- Documentation files (.md) in project root and subdirectories
- All files in `.claude/` directory
- Project-specific markdown files (CLAUDE.md, HANDOFF.md, etc.)

**Excluded Files**:
- Auto-generated test result files in `test-results/` directories
- Auto-generated report files in `playwright-report/` directories
- Temporary or cache-generated markdown files

**Implementation**: Use automated tools or manual checks to ensure compliance during development.

## Commit Message Format Rules

**IMPORTANT**: All commit messages MUST follow the Conventional Commit format. The project uses commitlint with husky hooks to enforce these rules.

### Required Configuration Files
- `.commitlintrc.js` - Commitlint configuration
- `.husky/commit-msg` - Pre-commit hook for message validation
- `package.json` - Dependencies for commitlint and husky

### Allowed Commit Types
- `feat:` A new feature
- `fix:` A bug fix
- `docs:` Documentation only changes
- `style:` Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- `refactor:` A code change that neither fixes a bug nor adds a feature
- `perf:` A code change that improves performance
- `test:` Adding missing tests or correcting existing tests
- `build:` Changes that affect the build system or external dependencies
- `ci:` Changes to our CI configuration files and scripts
- `chore:` Other changes that don't modify src or test files
- `revert:` Reverts a previous commit
- `improve:` Improvements to existing features

### Format Requirements
- Type must be lowercase
- Subject must not be empty and should not end with a period
- Header must not exceed 100 characters
- Body and footer lines must not exceed 100 characters
- Body and footer should have leading blank lines

### Valid Examples
```bash
feat: add user authentication system
fix: resolve map rendering issue on mobile devices
docs: update API documentation for user endpoints
refactor: restructure database connection logic
improve: enhance map performance with virtualization
test: add unit tests for user service
build: update dependencies and fix vulnerabilities
ci: add automated testing workflow
chore: update .gitignore and cleanup temp files
```

### Invalid Examples (Will be rejected)
```bash
Add new feature              # Missing type
FEAT: add feature           # Type not lowercase
feat:                       # Empty subject
feat: add feature.          # Subject ends with period
feat: this is a very long commit message that exceeds the maximum character limit of 100 characters # Too long
```

**Note**: Commit messages that don't follow this format will be rejected by the pre-commit hook.

## Cloud Run and Docker Deployment Patterns

### Docker Build Commands
```bash
# Local development build
docker build -t bocchi-api:dev .

# Production multi-arch build
docker buildx build --platform linux/amd64,linux/arm64 -t bocchi-api:latest .

# Build with specific tag for Cloud Run
docker build -t gcr.io/YOUR_PROJECT_ID/bocchi-api:v1.0.0 .
```

### Cloud Run Deployment Commands

#### Automated Deployment Script
```bash
# Development environment
cd api
./scripts/build.sh dev YOUR_PROJECT_ID asia-northeast1

# Production environment  
./scripts/build.sh prod YOUR_PROJECT_ID asia-northeast1

# Custom region
./scripts/build.sh staging YOUR_PROJECT_ID us-central1
```

#### Manual Cloud Run Commands
```bash
# Deploy with minimal configuration
gcloud run deploy bocchi-api-dev \
  --image=gcr.io/YOUR_PROJECT_ID/bocchi-api:latest \
  --platform=managed \
  --region=asia-northeast1 \
  --allow-unauthenticated

# Deploy with production settings
gcloud run deploy bocchi-api-prod \
  --image=gcr.io/YOUR_PROJECT_ID/bocchi-api:latest \
  --platform=managed \
  --region=asia-northeast1 \
  --allow-unauthenticated \
  --port=8080 \
  --memory=1Gi \
  --cpu=2 \
  --max-instances=10 \
  --min-instances=1 \
  --timeout=300 \
  --concurrency=100

# Update existing service with new image
gcloud run services update bocchi-api-dev \
  --image=gcr.io/YOUR_PROJECT_ID/bocchi-api:latest \
  --region=asia-northeast1
```

### Google Container Registry Commands
```bash
# Configure Docker for GCR
gcloud auth configure-docker

# Push image to GCR
docker push gcr.io/YOUR_PROJECT_ID/bocchi-api:latest

# List images in registry
gcloud container images list --repository=gcr.io/YOUR_PROJECT_ID

# List image tags
gcloud container images list-tags gcr.io/YOUR_PROJECT_ID/bocchi-api

# Delete old images (cleanup)
gcloud container images delete gcr.io/YOUR_PROJECT_ID/bocchi-api:old-tag --force-delete-tags
```

### Secret Management Commands
```bash
# Create secrets in Google Secret Manager
echo "your-tidb-password" | gcloud secrets create tidb-password-dev --data-file=-
echo "your-new-relic-key" | gcloud secrets create new-relic-license-key-dev --data-file=-
echo "your-sentry-dsn" | gcloud secrets create sentry-dsn-dev --data-file=-

# Update existing secret
echo "new-password" | gcloud secrets versions add tidb-password-dev --data-file=-

# Access secret value (for testing)
gcloud secrets versions access latest --secret="tidb-password-dev"

# List all secrets
gcloud secrets list

# Grant service account access to secret
gcloud secrets add-iam-policy-binding tidb-password-dev \
  --member="serviceAccount:bocchi-cloud-run-dev@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

### Terraform Infrastructure Commands
```bash
# Initialize Terraform with backend
cd infra
terraform init

# Plan infrastructure changes
terraform plan -var="gcp_project_id=YOUR_PROJECT_ID" -var="environment=dev"

# Apply infrastructure
terraform apply -var="gcp_project_id=YOUR_PROJECT_ID" -var="environment=dev"

# Destroy infrastructure (caution!)
terraform destroy -var="gcp_project_id=YOUR_PROJECT_ID" -var="environment=dev"

# Format Terraform files
terraform fmt -recursive

# Validate Terraform configuration
terraform validate

# Show current state
terraform show

# Import existing resources
terraform import google_cloud_run_v2_service.api projects/YOUR_PROJECT_ID/locations/asia-northeast1/services/bocchi-api-dev
```

### Monitoring and Health Check Commands
```bash
# Health check endpoints
curl https://bocchi-api-dev-xxx.a.run.app/health
curl https://bocchi-api-dev-xxx.a.run.app/health/detailed

# View Cloud Run logs
gcloud run services logs read bocchi-api-dev --region=asia-northeast1 --limit=50

# Follow real-time logs
gcloud run services logs tail bocchi-api-dev --region=asia-northeast1

# Get service details
gcloud run services describe bocchi-api-dev --region=asia-northeast1

# List all Cloud Run services
gcloud run services list

# Check service URL
gcloud run services describe bocchi-api-dev \
  --region=asia-northeast1 \
  --format="value(status.url)"
```

### Environment Variable Patterns
```bash
# Development environment variables
export ENV=development
export LOG_LEVEL=DEBUG
export TIDB_HOST=localhost
export PORT=8080

# Production environment variables
export ENV=production
export LOG_LEVEL=INFO
export PORT=8080
export NEW_RELIC_LICENSE_KEY=$(gcloud secrets versions access latest --secret="new-relic-license-key-prod")
export SENTRY_DSN=$(gcloud secrets versions access latest --secret="sentry-dsn-prod")
```

### Docker Compose for Local Development
```bash
# Start complete development environment
cd api
docker-compose up -d

# Start only database
docker-compose up -d mysql

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Remove volumes (reset data)
docker-compose down -v

# Rebuild services
docker-compose up --build
```

### Troubleshooting Cloud Run Deployments
```bash
# Check deployment status
gcloud run services describe bocchi-api-dev --region=asia-northeast1

# View recent revisions
gcloud run revisions list --service=bocchi-api-dev --region=asia-northeast1

# Rollback to previous revision
gcloud run services update-traffic bocchi-api-dev \
  --to-revisions=bocchi-api-dev-00002-abc=100 \
  --region=asia-northeast1

# Delete failed revisions
gcloud run revisions delete bocchi-api-dev-00003-def --region=asia-northeast1

# Test service connectivity
curl -H "Authorization: Bearer $(gcloud auth print-access-token)" \
  https://bocchi-api-dev-xxx.a.run.app/health

# Check IAM permissions
gcloud run services get-iam-policy bocchi-api-dev --region=asia-northeast1
```
