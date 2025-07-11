# TDD+BDD Hybrid Implementation Guide

## OVERVIEW

This guide provides practical implementation patterns for applying TDD+BDD hybrid methodology in both frontend (React/Next.js) and backend (Go) development for the Bocchi The Map project.

## PRACTICAL WORKFLOW

### Step-by-Step Process

1. **Define BDD Scenario** (5-10 minutes)
   ```bash
   # Create or update BDD test file
   touch api/tests/e2e/feature_name_test.go
   # OR for frontend
   touch web/e2e/feature-name.spec.ts
   ```

2. **Run BDD Test** (Red)
   ```bash
   # Backend
   cd api && go test ./tests/e2e/feature_name_test.go
   
   # Frontend
   cd web && npm run test:e2e -- feature-name.spec.ts
   ```

3. **Implement with TDD Cycles** (30-60 minutes)
   - Write unit test (Red)
   - Write minimal implementation (Green)
   - Refactor (Refactor)
   - Repeat until BDD test passes

4. **Final Integration** (5-10 minutes)
   ```bash
   # Run all tests
   make test-all
   # OR
   npm run test && go test ./...
   ```

## BACKEND IMPLEMENTATION (Go)

### Tooling Setup

#### Ginkgo (BDD Framework)
```bash
# Install Ginkgo
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install github.com/onsi/gomega@latest

# Initialize test suite
ginkgo bootstrap
```

#### Go Test (TDD Framework)
```bash
# Standard Go testing
go test ./...

# With coverage
go test -cover ./...

# With race detection
go test -race ./...
```

### BDD Test Structure (Backend)
```go
// api/tests/e2e/user_authentication_test.go
package e2e_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("User Authentication", func() {
    Context("When user provides valid credentials", func() {
        It("should authenticate successfully", func() {
            // Given: Valid user credentials
            credentials := UserCredentials{
                Email:    "test@example.com",
                Password: "validpassword",
            }
            
            // When: User attempts to authenticate
            token, err := authService.Authenticate(credentials)
            
            // Then: Authentication succeeds
            Expect(err).ToNot(HaveOccurred())
            Expect(token).ToNot(BeEmpty())
        })
    })
})
```

### TDD Test Structure (Backend)
```go
// api/internal/domain/user/user_test.go
package user_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestUser_Create(t *testing.T) {
    // Arrange
    userData := UserData{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Act
    user, err := NewUser(userData)
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, "John Doe", user.Name)
    assert.Equal(t, "john@example.com", user.Email)
}
```

### Layer-Specific Patterns (Backend)

#### Interface Layer (Handlers)
```go
// BDD-focused: Test HTTP endpoints
func TestUserHandler_CreateUser(t *testing.T) {
    // Given: Valid user data
    userData := `{"name":"John","email":"john@example.com"}`
    req := httptest.NewRequest("POST", "/users", strings.NewReader(userData))
    
    // When: Making request
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    
    // Then: User created successfully
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

#### Application Layer (Services)
```go
// TDD-focused: Test business logic
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    // Act
    user, err := service.CreateUser(userData)
    
    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

#### Domain Layer (Core Logic)
```go
// Pure TDD: Test business rules
func TestUser_ValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "invalid", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## FRONTEND IMPLEMENTATION (React/Next.js)

### Tooling Setup

#### Vitest (Unit Testing)
```bash
# Install Vitest
npm install --save-dev vitest @vitest/ui

# Vitest config (vitest.config.ts)
import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
  },
})
```

#### Playwright (E2E Testing)
```bash
# Install Playwright
npm install --save-dev @playwright/test

# Initialize
npx playwright install
```

#### React Testing Library
```bash
# Install RTL
npm install --save-dev @testing-library/react @testing-library/jest-dom
```

### BDD Test Structure (Frontend E2E)
```typescript
// web/e2e/user-authentication.spec.ts
import { test, expect } from '@playwright/test'

test.describe('User Authentication', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
  })

  test('should authenticate user with valid credentials', async ({ page }) => {
    // Given: User is on login page
    await expect(page.getByText('Login')).toBeVisible()
    
    // When: User enters valid credentials
    await page.fill('[data-testid="email"]', 'test@example.com')
    await page.fill('[data-testid="password"]', 'validpassword')
    await page.click('[data-testid="login-button"]')
    
    // Then: User is redirected to dashboard
    await expect(page.getByText('Dashboard')).toBeVisible()
  })
})
```

### TDD Test Structure (Frontend Component)
```typescript
// web/src/components/LoginForm/__tests__/LoginForm.test.tsx
import { render, screen, fireEvent } from '@testing-library/react'
import { LoginForm } from '../LoginForm'

describe('LoginForm', () => {
  test('should call onSubmit when form is submitted', () => {
    // Arrange
    const mockOnSubmit = vi.fn()
    render(<LoginForm onSubmit={mockOnSubmit} />)
    
    // Act
    fireEvent.change(screen.getByLabelText('Email'), {
      target: { value: 'test@example.com' }
    })
    fireEvent.change(screen.getByLabelText('Password'), {
      target: { value: 'password' }
    })
    fireEvent.click(screen.getByRole('button', { name: 'Login' }))
    
    // Assert
    expect(mockOnSubmit).toHaveBeenCalledWith({
      email: 'test@example.com',
      password: 'password'
    })
  })
})
```

### Layer-Specific Patterns (Frontend)

#### Component Layer (UI Components)
```typescript
// TDD-focused: Test component behavior
describe('SearchInput', () => {
  test('should emit search event when Enter is pressed', () => {
    // Given: Component with search handler
    const mockOnSearch = vi.fn()
    render(<SearchInput onSearch={mockOnSearch} />)
    
    // When: User types and presses Enter
    const input = screen.getByRole('textbox')
    fireEvent.change(input, { target: { value: 'tokyo' } })
    fireEvent.keyDown(input, { key: 'Enter' })
    
    // Then: Search handler is called
    expect(mockOnSearch).toHaveBeenCalledWith('tokyo')
  })
})
```

#### Zustand (State Management)
```typescript
// web/src/stores/__tests__/spot-store.test.ts
import { renderHook, act } from '@testing-library/react'
import { useSpotStore } from '../spot-store'

describe('SpotStore', () => {
  beforeEach(() => {
    useSpotStore.getState().reset()
  })

  test('should add spot to favorites', () => {
    // Arrange
    const { result } = renderHook(() => useSpotStore())
    const spot = { id: '1', name: 'Test Spot' }
    
    // Act
    act(() => {
      result.current.addToFavorites(spot)
    })
    
    // Assert
    expect(result.current.favorites).toContain(spot)
  })
})
```

#### MSW (Mock Service Worker) [PLANNED]
```typescript
// web/src/mocks/handlers.ts
// NOTE: MSW is planned but not yet installed. Run `npm install msw --save-dev` first.
import { rest } from 'msw'

export const handlers = [
  rest.get('/api/spots', (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json([
        { id: '1', name: 'Test Spot', rating: 4.5 }
      ])
    )
  }),
  
  rest.post('/api/auth/login', (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({ token: 'mock-jwt-token' })
    )
  })
]
```

### Test Helpers
```typescript
// web/src/test/test-utils.tsx
import { render, RenderOptions } from '@testing-library/react'
import { ReactElement } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

// Custom render function with providers
const AllTheProviders = ({ children }: { children: React.ReactNode }) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  })
  
  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  )
}

const customRender = (ui: ReactElement, options?: RenderOptions) =>
  render(ui, { wrapper: AllTheProviders, ...options })

export * from '@testing-library/react'
export { customRender as render }
```

## BEST PRACTICES

### Backend Best Practices

1. **Use Ginkgo for BDD scenarios**
   - Focus on user-facing behavior
   - Use descriptive Context and It blocks
   - Keep scenarios independent

2. **Use standard Go testing for TDD**
   - Fast unit tests for business logic
   - Table-driven tests for multiple scenarios
   - Mock external dependencies

3. **Organize tests by domain**
   - Mirror production code structure
   - Separate unit and integration tests
   - Use test helpers for common setup

### Frontend Best Practices

1. **Start with E2E BDD for new features**
   - Define user behavior first
   - Create executable UI specifications
   - Use as acceptance criteria

2. **Switch to TDD for component implementation**
   - Small, focused component tests
   - Fast feedback loop
   - High component coverage

3. **Maintain test independence**
   - Each test should be isolated
   - Use proper setup/cleanup
   - Mock external dependencies

4. **Keep tests accessible**
   - Use semantic queries (getByRole, getByLabelText)
   - Follow accessibility best practices
   - Test with screen readers in mind

5. **Balance test levels**
   - Few E2E tests (expensive but valuable)
   - Many component tests (fast and focused)
   - Strategic integration tests (stores, API calls)

## NEXT STEPS

1. **Implement Frontend Auth0 Unit Tests** (High Priority)
   - Create unit tests for Auth0 authentication hooks
   - Test protected route components
   - Mock Auth0 provider for component testing

2. **Complete MSW Setup** (Currently Planned)
   - Install MSW package (`npm install msw --save-dev`)
   - Configure MSW handlers
   - Create mock data factories
   - Integrate with test setup

3. **Create Component Test Templates**
   - Basic component test template
   - Component with hooks test template
   - Component with state management test template

4. **Expand Testing Infrastructure**
   - Visual regression testing setup
   - Performance testing integration
   - Accessibility testing automation

## MEASURING SUCCESS

### Frontend Metrics
- **E2E Coverage**: Critical user journey coverage
- **Component Coverage**: UI component test coverage
- **Performance**: Page load time, interaction responsiveness
- **Accessibility**: WCAG compliance, keyboard navigation

### Backend Metrics
- **BDD Coverage**: User story coverage by BDD scenarios
- **Unit Coverage**: Business logic test coverage
- **Integration Coverage**: API endpoint test coverage
- **Performance**: Response time, throughput

### Quality Metrics
- **Bug Density**: Defects per feature/component
- **Test Execution Time**: Time to run test suites
- **Refactoring Confidence**: Safe code changes enabled by tests
- **Team Productivity**: Feature delivery speed and quality

This implementation guide provides the practical foundation for applying TDD+BDD hybrid methodology effectively in both frontend and backend development.
