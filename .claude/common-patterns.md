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

    // Persistence
    if err := c.repo.Create(ctx, entity); err != nil {
        c.logger.Error().Err(err).Msg("failed to create entity")
        return nil, fmt.Errorf("failed to create entity: %w", err)
    }

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
