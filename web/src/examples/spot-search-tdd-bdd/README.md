# TDD+BDD Hybrid Implementation Example: Spot Search Feature

This directory contains a complete example implementation of the **TDD+BDD Hybrid methodology** applied to a spot search feature for the Bocchi The Map application.

## 🎯 Feature Overview

**User Story**: As a solo traveler, I want to search for spots suitable for solo activities, so that I can find comfortable places to visit alone.

**Key Scenarios**:
- Search for spots with keywords
- Filter by solo-friendliness
- View search results with ratings
- Handle errors gracefully
- Mobile-responsive interface

## 🔄 TDD+BDD Hybrid Methodology Demonstrated

This example follows the **Outside-In TDD with BDD** approach:

```
BDD (E2E Tests) → TDD (Component Tests) → Integration → E2E Validation
     ↓                    ↓                   ↓              ↓
   RED PHASE           RED PHASE         GREEN PHASE    GREEN PHASE
```

### 1. **BDD Phase (Outside-In) - RED**
- **File**: `e2e/spot-search.spec.ts`
- **Purpose**: Define user behavior and acceptance criteria
- **Result**: Failing E2E tests that describe what users expect

### 2. **TDD Phase (Inside-Out) - RED → GREEN → REFACTOR**
- **Components**: Individual component implementations driven by unit tests
- **Custom Hooks**: Business logic with comprehensive test coverage
- **Result**: Working components that can fulfill user scenarios

### 3. **Integration Phase - GREEN**
- **File**: `integration/__tests__/search-integration.test.tsx`
- **Purpose**: Verify components work together correctly
- **Result**: Coordinated system behavior

### 4. **E2E Validation - GREEN**
- **Result**: Original BDD tests now pass, confirming user scenarios work

## 📁 Project Structure

```
spot-search-tdd-bdd/
├── README.md                           # This file
├── e2e/
│   └── spot-search.spec.ts            # BDD E2E tests (Outside-In)
├── components/
│   ├── __tests__/
│   │   └── search-input.test.tsx      # TDD unit tests
│   └── search-input.tsx               # Component implementation
├── hooks/
│   ├── __tests__/
│   │   └── use-spot-search.test.ts    # TDD hook tests
│   └── use-spot-search.ts             # Custom hook implementation
└── integration/
    └── __tests__/
        └── search-integration.test.tsx # Integration tests
```

## 🧪 Test Types and Their Roles

### 1. **BDD E2E Tests** (`e2e/spot-search.spec.ts`)
- **Framework**: Playwright
- **Purpose**: User acceptance testing
- **Focus**: Complete user journeys and business value
- **Pattern**: Given-When-Then scenarios

```typescript
test('When I search for "quiet cafe", Then I should see relevant solo-friendly results', async ({ page }) => {
  // Given - Search page is loaded
  await expect(page.getByTestId('search-page')).toBeVisible()
  
  // When - I enter a search query
  const searchInput = page.getByTestId('search-input')
  await searchInput.fill('quiet cafe')
  await page.keyboard.press('Enter')
  
  // Then - Search results should be displayed
  await expect(page.getByTestId('search-results')).toBeVisible()
})
```

### 2. **TDD Unit Tests** (`components/__tests__/`, `hooks/__tests__/`)
- **Framework**: Vitest + React Testing Library
- **Purpose**: Component behavior and business logic
- **Focus**: Implementation details and edge cases
- **Pattern**: RED → GREEN → REFACTOR cycles

```typescript
describe('Given the SearchInput component is rendered', () => {
  describe('When user types in the search input', () => {
    it('Then the input value should update correctly', async () => {
      // TDD implementation test
    })
  })
})
```

### 3. **Integration Tests** (`integration/__tests__/`)
- **Framework**: Vitest + React Testing Library
- **Purpose**: Component interaction and state flow
- **Focus**: How components work together
- **Pattern**: Mock dependencies, test integration points

## 🛠 Implementation Process

### Step 1: BDD E2E Tests (RED)
1. **Analyze user requirements** and create user stories
2. **Write failing E2E tests** that describe expected behavior
3. **Define test data** and expected outcomes
4. **Run tests** - they should fail (RED)

### Step 2: TDD Component Implementation (RED → GREEN → REFACTOR)

#### 2a. SearchInput Component
1. **RED**: Write failing unit tests for SearchInput
2. **GREEN**: Implement minimal SearchInput to pass tests
3. **REFACTOR**: Clean up code while keeping tests green

#### 2b. useSpotSearch Hook
1. **RED**: Write failing tests for search logic
2. **GREEN**: Implement search functionality
3. **REFACTOR**: Optimize and improve code quality

### Step 3: Integration Testing (GREEN)
1. **Write integration tests** to verify component cooperation
2. **Mock external dependencies** (APIs, services)
3. **Test state flow** between components
4. **Verify error handling** across the system

### Step 4: E2E Validation (GREEN)
1. **Run original E2E tests** to verify they now pass
2. **Fix any integration issues** discovered
3. **Ensure all scenarios work** end-to-end

## 🏃‍♂️ Running the Tests

### Run All Tests
```bash
cd web

# Run unit tests
pnpm test

# Run E2E tests
pnpm test:e2e

# Run with coverage
pnpm test:coverage
```

### Run Specific Test Categories
```bash
# Run only SearchInput tests
pnpm test src/examples/spot-search-tdd-bdd/components/__tests__/search-input.test.tsx

# Run only useSpotSearch tests
pnpm test src/examples/spot-search-tdd-bdd/hooks/__tests__/use-spot-search.test.ts

# Run integration tests
pnpm test src/examples/spot-search-tdd-bdd/integration

# Run E2E tests for this feature
pnpm test:e2e src/examples/spot-search-tdd-bdd/e2e/spot-search.spec.ts
```

### Test in Watch Mode
```bash
# Watch mode for rapid TDD cycles
pnpm test --watch src/examples/spot-search-tdd-bdd
```

## 📊 Test Coverage Analysis

### Test Pyramid Distribution
- **E2E Tests**: 5 comprehensive user scenarios
- **Integration Tests**: 8 component interaction tests  
- **Unit Tests**: 25+ detailed component and hook tests

### Coverage Metrics (Target)
- **Component Tests**: 95%+ line coverage
- **Hook Tests**: 100% line coverage
- **Integration Tests**: All critical user flows
- **E2E Tests**: Primary user journeys

## 🎯 Key TDD+BDD Principles Demonstrated

### 1. **Outside-In Development**
- Start with user needs (BDD)
- Drive implementation from external behavior
- Build components to satisfy user scenarios

### 2. **RED-GREEN-REFACTOR Cycles**
- Write failing tests first (RED)
- Implement minimal code to pass (GREEN)
- Improve code quality (REFACTOR)

### 3. **Test-Driven Design**
- Tests drive component APIs
- Implementation follows test requirements
- Clean interfaces emerge naturally

### 4. **Comprehensive Coverage**
- User behavior (E2E)
- Component behavior (Unit)
- System integration (Integration)

### 5. **Confidence in Changes**
- Safe refactoring with test coverage
- Regression detection
- Reliable deployment pipeline

## 🔧 Development Workflow

### Adding New Features
1. **Write BDD scenario** for user behavior
2. **Create failing E2E test**
3. **Implement with TDD**:
   - Write failing unit tests
   - Implement component
   - Refactor code
4. **Add integration tests**
5. **Verify E2E tests pass**

### Debugging Test Failures
1. **E2E failure**: Check integration and component tests
2. **Integration failure**: Examine component interactions
3. **Unit failure**: Fix component implementation
4. **Work from inside-out** to identify root cause

### Refactoring Safely
1. **Ensure all tests are green**
2. **Refactor implementation**
3. **Keep tests green throughout**
4. **Update tests only if behavior changes**

## 📚 Learning Outcomes

This example demonstrates:

### TDD Benefits
- ✅ **Test-driven component design**
- ✅ **High confidence in code changes**
- ✅ **Comprehensive edge case coverage**
- ✅ **Clean, focused implementations**

### BDD Benefits
- ✅ **User-centric development**
- ✅ **Living documentation**
- ✅ **Stakeholder communication**
- ✅ **Acceptance criteria validation**

### Hybrid Approach Benefits
- ✅ **Complete test coverage**
- ✅ **Multiple feedback loops**
- ✅ **Quality at all levels**
- ✅ **Maintainable codebase**

## 🚀 Next Steps

To extend this example:

1. **Add more components** (FilterPanel, SpotItem, SearchResults)
2. **Implement real API integration**
3. **Add performance testing**
4. **Include accessibility testing**
5. **Add visual regression tests**

## 📖 Related Documentation

- [Frontend TDD+BDD Methodology Guide](../../../.claude/frontend-tdd-bdd-hybrid.md)
- [Common Testing Patterns](../../../.claude/common-patterns.md)
- [Project Testing Guidelines](../../../README.md)

---

**Note**: This example serves as a practical template for implementing the TDD+BDD hybrid methodology in React/Next.js applications. The patterns and practices demonstrated here can be applied to any frontend feature development.