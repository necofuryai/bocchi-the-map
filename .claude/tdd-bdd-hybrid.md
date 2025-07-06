# TDD + BDD Hybrid Development Methodology

## OVERVIEW

This document outlines the hybrid approach combining Kent Beck's Test-Driven Development (TDD) with Behavior-Driven Development (BDD) principles, tailored for the Bocchi The Map project.

## CORE PHILOSOPHY

The TDD+BDD hybrid approach leverages the strengths of both methodologies:
- **BDD**: Defines high-level behavior and user stories (Outside)
- **TDD**: Drives internal implementation details (Inside)
- **Integration**: Creates a seamless flow from user requirements to clean implementation

## THE THREE HYBRID APPROACHES

### 1. Outside-In TDD with BDD

This approach starts from the user's perspective and works inward:

```
BDD (External) → TDD (Internal)
```

1. **Write BDD Scenario** (Red)
   - Define user story using Given-When-Then
   - Create high-level E2E test with Ginkgo
   - Ensure test fails for the right reason

2. **Implement with TDD** (Green)
   - For each component needed to pass the BDD test:
     - Write unit test (TDD Red)
     - Implement minimal code (TDD Green)
     - Refactor (TDD Refactor)

3. **Validate Integration** (Refactor)
   - Ensure BDD test passes
   - Refactor the overall design
   - Keep all tests green

### 2. Double-Loop TDD

This creates two feedback loops - one for behavior, one for implementation:

```go
// Outer Loop (BDD) - api/tests/e2e/auth_test.go
Describe("User Authentication", func() {
    Context("When a new user signs up", func() {
        It("should create a user profile", func() {
            // This drives the inner loop
        })
    })
})

// Inner Loop (TDD) - api/internal/domain/user/user_test.go
func TestUserCreation(t *testing.T) {
    // Red: Write failing test
    // Green: Implement
    // Refactor: Clean up
}
```

### 3. Specification by Example + TDD

This approach uses concrete examples to drive development:

1. **Specification Phase**
   - Write BDD scenarios with concrete examples
   - Define acceptance criteria
   - Create test data fixtures

2. **Implementation Phase**
   - Use TDD to implement each component
   - Follow Red-Green-Refactor strictly
   - Keep scenarios as north star

3. **Validation Phase**
   - Run all BDD scenarios
   - Ensure examples pass
   - Document edge cases discovered

## PRACTICAL WORKFLOW

### Step-by-Step Process

1. **Feature Request Analysis**
   ```gherkin
   Feature: Solo-friendly spot review
   Scenario: User reviews a cafe for solo dining
   Given I am authenticated
   When I submit a review with solo-friendliness rating
   Then the review should be saved with my rating
   ```

2. **Create BDD Test Structure**
   ```go
   // api/tests/e2e/review_test.go
   var _ = Describe("Review API", func() {
       Context("Given I am authenticated", func() {
           BeforeEach(func() {
               // Setup authentication
           })
           
           Context("When I submit a review", func() {
               It("Then the review should be saved", func() {
                   // E2E test implementation
               })
           })
       })
   })
   ```

3. **TDD Implementation**
   ```go
   // api/internal/domain/review/review_test.go
   func TestReviewCreation(t *testing.T) {
       // TDD cycle for domain logic
   }
   
   // api/internal/application/review_service_test.go
   func TestReviewService_Create(t *testing.T) {
       // TDD cycle for service layer
   }
   ```

## LAYER-SPECIFIC STRATEGIES

### Interface Layer (Handlers/Controllers)
- **Primary**: BDD approach
- **Focus**: User interactions and API contracts
- **Tools**: Ginkgo for E2E tests

### Application Layer (Services)
- **Primary**: TDD with BDD context
- **Focus**: Business logic orchestration
- **Tools**: Standard Go testing with clear scenarios

### Domain Layer (Core Business Logic)
- **Primary**: Pure TDD
- **Focus**: Business rules and entities
- **Tools**: Table-driven tests, property-based testing

### Infrastructure Layer (Adapters)
- **Primary**: TDD with integration tests
- **Focus**: External system interactions
- **Tools**: Mocks and test containers

## COMMIT GUIDELINES FOR HYBRID APPROACH

### Commit Types

1. **BDD Commits** (Behavior)
   ```
   feat(review): add BDD scenario for solo-friendly ratings
   
   - Added E2E test for review submission
   - Defined acceptance criteria
   ```

2. **TDD Commits** (Implementation)
   ```
   feat(review): implement review domain model
   
   - Added Review entity with TDD
   - Implemented validation rules
   - All tests passing
   ```

3. **Refactoring Commits** (Structure)
   ```
   refactor(review): extract rating calculation logic
   
   - No behavior change
   - All tests still passing
   ```

## TESTING PATTERNS

### BDD Test Structure
```go
Describe("Feature", func() {
    Context("Given <precondition>", func() {
        BeforeEach(func() {
            // Setup
        })
        
        When("When <action>", func() {
            It("Then <outcome>", func() {
                // Assertion
            })
        })
    })
})
```

### TDD Test Structure
```go
func TestComponent_Method(t *testing.T) {
    // Arrange
    sut := NewComponent()
    
    // Act
    result := sut.Method()
    
    // Assert
    assert.Equal(t, expected, result)
}
```

## TOOLING INTEGRATION

### Ginkgo (BDD)
- Use for E2E and integration tests
- Focus on user scenarios
- Run with: `ginkgo -v`

### Go Test (TDD)
- Use for unit tests
- Focus on implementation details
- Run with: `go test ./...`

### Test Helpers
- CommonTestSuite for shared setup
- FixtureManager for test data
- Authentication helpers for secure endpoints

## BEST PRACTICES

1. **Start with BDD for new features**
   - Define behavior first
   - Create executable specifications
   - Use as acceptance criteria

2. **Switch to TDD for implementation**
   - Small, focused tests
   - Fast feedback loop
   - High code coverage

3. **Maintain test independence**
   - Each test should be isolated
   - Use proper setup/teardown
   - Avoid test interdependencies

4. **Keep tests readable**
   - Use descriptive names
   - Follow AAA pattern
   - Document complex scenarios

5. **Balance test levels**
   - Few E2E tests (BDD)
   - Many unit tests (TDD)
   - Strategic integration tests

## ANTI-PATTERNS TO AVOID

1. **Writing BDD tests for internal logic**
   - BDD is for user-facing behavior
   - Use TDD for implementation details

2. **Skipping TDD cycles**
   - Always follow Red-Green-Refactor
   - Don't write code without tests

3. **Mixing test concerns**
   - Keep BDD and TDD tests separate
   - Don't test implementation in BDD

4. **Over-mocking in BDD**
   - BDD tests should be as real as possible
   - Mock only external dependencies

## EXAMPLE: IMPLEMENTING A NEW FEATURE

Let's implement "Find solo-friendly spots nearby":

1. **BDD Scenario** (.claude/features/find-spots.feature)
2. **E2E Test** (api/tests/e2e/spot_search_test.go)
3. **Domain TDD** (api/internal/domain/spot/search_test.go)
4. **Service TDD** (api/internal/application/spot_service_test.go)
5. **Handler BDD** (api/interfaces/http/handlers/spot_handler_test.go)

This creates a complete test pyramid from user story to implementation.

## MEASURING SUCCESS

- **BDD Metrics**: Scenario coverage, acceptance test pass rate
- **TDD Metrics**: Unit test coverage, test execution time
- **Combined Metrics**: Feature completion rate, defect density

## CONTINUOUS IMPROVEMENT

1. Regular retrospectives on test effectiveness
2. Refactor test code as production code
3. Update this guide based on team learnings
4. Share testing patterns in common-patterns.md