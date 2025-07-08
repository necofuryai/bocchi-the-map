# Design Philosophy

## Three Masters' Approach to Bocchi The Map

This document outlines how John Carmack, Robert C. Martin (Uncle Bob), and Rob Pike would approach the design of Bocchi The Map, and how their philosophies are integrated into our architecture.

### John Carmack's Approach

**Core Principles:**
- **Performance First**: Optimize critical paths, especially map rendering and location queries
- **Pragmatic Simplicity**: Choose the simplest solution that performs well
- **Direct Implementation**: Avoid unnecessary abstractions that hurt performance

**Applied to Bocchi The Map:**
```go
// Carmack-style: Direct, efficient implementation for hot paths
type LocationCache struct {
    mu       sync.RWMutex
    data     map[int64]*Location
    spatial  *QuadTree // Spatial indexing for fast queries
}

func (c *LocationCache) FindNearby(lat, lng float64, radius float64) []*Location {
    // Direct spatial query without abstraction layers
    return c.spatial.Query(lat, lng, radius)
}
```

**Key Decisions:**
- Use spatial indexing (QuadTree/R-tree) for location queries
- Implement caching at critical layers
- Profile and optimize database queries early
- Minimize API round trips with efficient data fetching

### Robert C. Martin's Approach

**Core Principles:**
- **Clean Architecture**: Dependency inversion and clear boundaries
- **SOLID Principles**: Single responsibility, open/closed, etc.
- **Test-Driven Development**: Tests drive design

**Applied to Bocchi The Map:**
```go
// Uncle Bob style: Clean architecture with clear boundaries
// Domain Layer (innermost)
type Location struct {
    ID          int64
    Name        string
    Coordinates Coordinates
    SoloRating  float64
}

// Use Case Layer
type FindSoloFriendlyLocationsUseCase interface {
    Execute(ctx context.Context, criteria SearchCriteria) ([]Location, error)
}

// Interface Adapters Layer
type LocationRepository interface {
    FindByArea(ctx context.Context, area Area) ([]Location, error)
    Save(ctx context.Context, location Location) error
}

// Infrastructure Layer (outermost)
type PostgresLocationRepository struct {
    db *sql.DB
}
```

**Key Decisions:**
- Onion Architecture with clear dependency rules
- Domain entities have no external dependencies
- Use cases orchestrate business logic
- Interfaces for all external dependencies
- Comprehensive test coverage at each layer

### Rob Pike's Approach

**Core Principles:**
- **Simplicity**: Clear is better than clever
- **Concurrency**: Design for concurrent execution
- **Composition**: Small interfaces, composed together

**Applied to Bocchi The Map:**
```go
// Pike style: Simple, concurrent, composable
type LocationService struct {
    repo   LocationRepository
    cache  Cache
    events EventBus
}

// Simple interface, does one thing well
type Reviewer interface {
    Review(locationID int64, rating Rating) error
}

// Concurrent processing with clear communication
func (s *LocationService) ProcessReviews(reviews <-chan Review) {
    for i := 0; i < runtime.NumCPU(); i++ {
        go func() {
            for review := range reviews {
                if err := s.processReview(review); err != nil {
                    log.Printf("failed to process review: %v", err)
                }
            }
        }()
    }
}
```

**Key Decisions:**
- Small, focused interfaces
- Embrace Go's concurrency model
- Error handling as values
- Flat package structure where sensible
- "Don't communicate by sharing memory; share memory by communicating"

## Unified Design Principles for Bocchi The Map

Combining the wisdom of these three masters:

### 1. Performance-Aware Clean Architecture
- Follow Clean Architecture principles BUT optimize critical paths
- Profile early and often (Carmack)
- Keep domain logic pure (Martin)
- Use concurrency effectively (Pike)

### 2. Pragmatic Testing Strategy
- TDD for business logic (Martin)
- Performance tests for critical paths (Carmack)
- Simple, table-driven tests (Pike)

### 3. Simple Abstractions, Complex Implementations
- Simple interfaces (Pike)
- Clean boundaries (Martin)
- Optimized implementations (Carmack)

### 4. Implementation Guidelines

```go
// Example: Location search combining all three philosophies

// Simple interface (Pike)
type LocationSearcher interface {
    Search(ctx context.Context, query Query) ([]Location, error)
}

// Clean architecture implementation (Martin)
type locationSearchService struct {
    repo      LocationRepository  // Interface, not concrete
    cache     Cache              // Performance optimization (Carmack)
    metrics   Metrics            // Observability
}

// Efficient implementation with concurrency (Carmack + Pike)
func (s *locationSearchService) Search(ctx context.Context, query Query) ([]Location, error) {
    // Check cache first (Carmack - performance)
    if cached, hit := s.cache.Get(query.Hash()); hit {
        return cached.([]Location), nil
    }
    
    // Concurrent search across shards (Pike)
    results := make(chan []Location, len(s.shards))
    errors := make(chan error, len(s.shards))
    
    for _, shard := range s.shards {
        go func(sh Shard) {
            locs, err := sh.Search(ctx, query)
            if err != nil {
                errors <- err
                return
            }
            results <- locs
        }(shard)
    }
    
    // Collect results
    var locations []Location
    for i := 0; i < len(s.shards); i++ {
        select {
        case locs := <-results:
            locations = append(locations, locs...)
        case err := <-errors:
            return nil, err
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    }
    
    // Cache results
    s.cache.Set(query.Hash(), locations, 5*time.Minute)
    
    return locations, nil
}
```

### 5. Practical Application Checklist

When designing a new feature:

1. **Start Simple** (Pike)
   - What's the simplest interface that could work?
   - Can it be composed from existing components?

2. **Apply Clean Architecture** (Martin)
   - What layer does this belong to?
   - What are the dependencies?
   - How will we test it?

3. **Optimize Critical Paths** (Carmack)
   - Is this on a performance-critical path?
   - What's the expected load?
   - Where should we add caching/optimization?

4. **Design for Concurrency** (Pike)
   - Can this be parallelized?
   - What's the communication pattern?
   - How do we handle errors in concurrent operations?

### 6. Anti-Patterns to Avoid

- **Over-engineering** (violates Carmack's pragmatism)
- **Tight coupling** (violates Martin's principles)
- **Clever code** (violates Pike's simplicity)
- **Premature optimization** without profiling
- **Ignoring error handling** in concurrent code
- **Large interfaces** that do too much

## Conclusion

By combining these three masters' approaches, Bocchi The Map achieves:
- **Performance**: Optimized for real-world usage
- **Maintainability**: Clean architecture and clear boundaries
- **Simplicity**: Easy to understand and modify
- **Scalability**: Designed for concurrent operation

Remember: "Make it work, make it right, make it fast" - in that order, with tests at every step.