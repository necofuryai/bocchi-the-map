# Frontend TDD + BDD Hybrid Development Methodology

## OVERVIEW

This document outlines the hybrid approach combining Test-Driven Development (TDD) with Behavior-Driven Development (BDD) principles, specifically tailored for React/Next.js frontend development in the Bocchi The Map project.

## CORE PHILOSOPHY

The Frontend TDD+BDD hybrid approach leverages the strengths of both methodologies for web development:
- **BDD**: Defines user interactions and UI behavior (Outside-In)
- **TDD**: Drives component implementation and logic (Inside-Out)
- **Integration**: Creates a seamless flow from user stories to maintainable React components

## THE THREE HYBRID APPROACHES FOR FRONTEND

### 1. Outside-In TDD with BDD (UI-First)

This approach starts from the user interface and works inward:

```
E2E BDD Scenario → Component Integration → Unit TDD
```

1. **Write E2E BDD Scenario** (Red)
   - Define user story using Given-When-Then
   - Create high-level E2E test with Playwright
   - Ensure test fails for the right reason

2. **Implement with Component TDD** (Green)
   - For each React component needed to pass the E2E test:
     - Write component test (TDD Red)
     - Implement minimal component (TDD Green)
     - Refactor component (TDD Refactor)

3. **Validate Integration** (Refactor)
   - Ensure E2E test passes
   - Refactor the overall component structure
   - Keep all tests green

### 2. Double-Loop TDD (User Story + Component Logic)

This creates two feedback loops - one for user behavior, one for component implementation:

```tsx
// Outer Loop (BDD) - E2E Test
test.describe('Spot Search Feature', () => {
  test('When user searches for spots, Then results should be displayed', async ({ page }) => {
    // This drives the inner loop
  })
})

// Inner Loop (TDD) - Component Test
describe('SearchInput Component', () => {
  it('should filter spots based on search query', () => {
    // Red: Write failing test
    // Green: Implement
    // Refactor: Clean up
  })
})
```

### 3. Specification by Example + TDD (Design-First)

This approach uses UI mockups and concrete examples to drive development:

1. **Specification Phase**
   - Create UI mockups/wireframes
   - Write BDD scenarios with concrete examples
   - Define acceptance criteria with visual examples

2. **Implementation Phase**
   - Use TDD to implement each React component
   - Follow Red-Green-Refactor strictly
   - Keep UI specifications as north star

3. **Validation Phase**
   - Run all E2E scenarios
   - Ensure visual and functional examples pass
   - Document edge cases discovered

## PRACTICAL WORKFLOW

### Step-by-Step Process

1. **Feature Request Analysis**
   ```gherkin
   Feature: Solo-friendly spot search
   Scenario: User searches for quiet cafes
   Given I am on the search page
   When I enter "quiet cafe" in the search input
   Then I should see spots filtered by solo-friendliness
   ```

2. **Create E2E BDD Test Structure**
   ```typescript
   // web/e2e/spot-search.spec.ts
   test.describe('Spot Search Feature', () => {
     test.describe('Given I am on the search page', () => {
       test.beforeEach(async ({ page }) => {
         await page.goto('/search')
       })
       
       test('When I search for spots, Then filtered results should appear', async ({ page }) => {
         // E2E test implementation
       })
     })
   })
   ```

3. **TDD Component Implementation**
   ```typescript
   // web/src/components/search/__tests__/search-input.test.tsx
   describe('SearchInput Component', () => {
     it('should call onSearch when user types', () => {
       // TDD cycle for component logic
     })
   })
   
   // web/src/hooks/__tests__/use-spot-search.test.ts
   describe('useSpotSearch Hook', () => {
     it('should filter spots based on query', () => {
       // TDD cycle for custom hook
     })
   })
   ```

## LAYER-SPECIFIC STRATEGIES

### Presentation Layer (React Components)
- **Primary**: BDD approach for user interactions
- **Focus**: Component behavior, props, and user events
- **Tools**: React Testing Library with user-event
- **Pattern**: Given-When-Then for component behavior

### Logic Layer (Custom Hooks)
- **Primary**: TDD with BDD context
- **Focus**: Business logic, state management, side effects
- **Tools**: Vitest with React Testing Library renderHook
- **Pattern**: Red-Green-Refactor for hook logic

### Integration Layer (API Calls)
- **Primary**: TDD with mocking
- **Focus**: Data fetching, caching, error handling
- **Tools**: MSW (Mock Service Worker) for API mocking
- **Pattern**: Test doubles for external dependencies

### E2E Layer (User Journeys)
- **Primary**: Pure BDD
- **Focus**: Complete user workflows
- **Tools**: Playwright for full browser automation
- **Pattern**: User story scenarios

## FRONTEND-SPECIFIC PATTERNS

### Component Testing Pattern
```typescript
describe('Component Name', () => {
  describe('Given initial props', () => {
    describe('When user interacts', () => {
      it('Then component should behave correctly', () => {
        // Test implementation
      })
    })
  })
})
```

### Hook Testing Pattern
```typescript
describe('useCustomHook', () => {
  it('should return expected initial state', () => {
    // Arrange
    const { result } = renderHook(() => useCustomHook())
    
    // Act & Assert
    expect(result.current.state).toBe(expectedInitialState)
  })
})
```

### E2E Testing Pattern
```typescript
test.describe('Feature Name', () => {
  test.describe('Given user context', () => {
    test('When user action, Then expected outcome', async ({ page }) => {
      // Arrange
      await page.goto('/path')
      
      // Act
      await page.click('button')
      
      // Assert
      await expect(page.locator('result')).toBeVisible()
    })
  })
})
```

## COMMIT GUIDELINES FOR FRONTEND HYBRID APPROACH

### Commit Types

1. **BDD Commits** (User Behavior)
   ```
   feat(web): add BDD scenario for spot search functionality
   
   - Added E2E test for search user journey
   - Defined acceptance criteria for search results
   ```

2. **TDD Commits** (Component Implementation)
   ```
   feat(web): implement SearchInput component with TDD
   
   - Added SearchInput component with unit tests
   - Implemented search input validation
   - All tests passing
   ```

3. **Refactoring Commits** (Component Structure)
   ```
   refactor(web): extract search logic into custom hook
   
   - No behavior change
   - Improved component separation
   - All tests still passing
   ```

## TESTING PATTERNS

### BDD E2E Test Structure
```typescript
test.describe('Feature', () => {
  test.describe('Given <precondition>', () => {
    test.beforeEach(async ({ page }) => {
      // Setup
    })
    
    test('When <action>, Then <outcome>', async ({ page }) => {
      // Test implementation
    })
  })
})
```

### TDD Component Test Structure
```typescript
describe('Component', () => {
  it('should behave correctly when props change', () => {
    // Arrange
    const props = { /* test props */ }
    
    // Act
    render(<Component {...props} />)
    
    // Assert
    expect(screen.getByRole('button')).toBeInTheDocument()
  })
})
```

## TOOLING INTEGRATION

### Playwright (E2E BDD)
- Use for end-to-end user scenarios
- Focus on complete user journeys
- Run with: `npm run test:e2e`

### Vitest + React Testing Library (TDD)
- Use for component unit and integration tests
- Focus on component behavior and hooks
- Run with: `npm run test`

### MSW (Mock Service Worker)
- Use for API mocking in tests
- Provides realistic API responses
- Integrates with both unit and E2E tests

### Test Helpers
- Custom render functions for providers
- Test utilities for common patterns
- Accessibility testing helpers

## BEST PRACTICES

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
   - Strategic integration tests (hooks, API calls)

## ANTI-PATTERNS TO AVOID

1. **Testing implementation details**
   - Don't test internal component state
   - Focus on user-visible behavior
   - Avoid testing CSS classes or internal methods

2. **Over-mocking in E2E tests**
   - E2E tests should be as real as possible
   - Mock only external APIs, not internal components
   - Use test databases for realistic data

3. **Skipping accessibility testing**
   - Always test with screen reader queries
   - Ensure keyboard navigation works
   - Check color contrast and focus indicators

4. **Writing tests after implementation**
   - Always follow TDD Red-Green-Refactor cycle
   - Don't write tests just to increase coverage
   - Let tests drive your component design

## EXAMPLE: IMPLEMENTING A NEW FEATURE

Let's implement "Spot Search with Filters":

1. **E2E BDD Scenario** (web/e2e/spot-search.spec.ts)
2. **Search Component TDD** (web/src/components/search/__tests__/search-input.test.tsx)
3. **Filter Hook TDD** (web/src/hooks/__tests__/use-search-filters.test.ts)
4. **Integration Tests** (web/src/components/search/__tests__/search-page.test.tsx)
5. **API Integration** (web/src/services/__tests__/spot-api.test.ts)

This creates a complete test pyramid from user story to component implementation.

## MEASURING SUCCESS

- **E2E Metrics**: User journey coverage, critical path testing
- **Component Metrics**: Component test coverage, accessibility compliance
- **Performance Metrics**: Page load time, interaction responsiveness
- **Quality Metrics**: Bug density, user satisfaction scores

## CONTINUOUS IMPROVEMENT

1. Regular retrospectives on test effectiveness
2. Refactor test code as production code
3. Update this guide based on team learnings
4. Share testing patterns in common-patterns.md
5. Monitor test execution time and optimize slow tests

## NEXT STEPS

1. Create React component test templates
2. Set up MSW for API mocking
3. Implement example feature using this methodology
4. Create visual regression testing setup
5. Add performance testing to the pipeline

This frontend-specific TDD+BDD hybrid approach ensures both user satisfaction and technical quality for the Bocchi The Map web application.