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
