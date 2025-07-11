# TDD+BDD Hybrid Development Methodology

## OVERVIEW

This document outlines the hybrid approach combining Kent Beck's Test-Driven Development (TDD) with Behavior-Driven Development (BDD) principles, tailored for the Bocchi The Map project.

## CORE PHILOSOPHY

The TDD+BDD hybrid approach leverages the strengths of both methodologies:
- **BDD**: Defines high-level behavior and user stories (Outside)
- **TDD**: Drives internal implementation details (Inside)
- **Integration**: Creates a seamless flow from user requirements to clean implementation

## PRODUCTION IMPLEMENTATION: AUTH0 AUTHENTICATION

This methodology has been successfully applied to production features, particularly Auth0 authentication integration:

### Real-World Results
- **34 comprehensive test cases** covering authentication flows
- **97% pass rate** in production testing
- **E2E BDD scenarios** for login/logout, protected routes, and error handling
- **TDD unit tests** for backend auth services and rate limiting
- **Integration tests** for database sessions and API authentication

### Key Benefits Demonstrated
- **User-centric approach**: BDD scenarios ensure authentication flows match user expectations
- **Implementation quality**: TDD ensures robust internal auth logic
- **Maintainability**: Tests serve as documentation for complex auth flows
- **Confidence**: High test coverage enables safe refactoring and feature additions

*See `.claude/tdd-bdd-examples.md` for detailed implementation examples.*

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
     - Write unit test (Red)
     - Write minimal implementation (Green)
     - Refactor while keeping tests green
   - Continue until BDD test passes

3. **Refactor the Whole** (Refactor)
   - Clean up integration between components
   - Ensure all tests remain green
   - Improve overall design

### 2. Double-Loop TDD

This approach creates two nested TDD loops:

```
Outer Loop: BDD Feature Test
Inner Loop: TDD Unit Tests
```

1. **Write failing BDD test** (Outer Red)
2. **Write failing TDD unit test** (Inner Red)
3. **Make TDD test pass** (Inner Green)
4. **Refactor unit implementation** (Inner Refactor)
5. **Repeat Inner Loop** until BDD test passes (Outer Green)
6. **Refactor entire feature** (Outer Refactor)

### 3. Specification by Example + TDD

This approach uses concrete examples to drive both behavior and implementation:

1. **Define Examples** - Create concrete scenarios with real data
2. **Write BDD Tests** - Use examples as basis for Given-When-Then scenarios
3. **Apply TDD** - Use examples to drive unit test creation
4. **Validate** - Ensure examples work in both BDD and TDD contexts

## LAYER-SPECIFIC STRATEGIES

### Interface Layer (Handlers/Controllers)
- **Primary**: BDD for user-facing behavior
- **Secondary**: TDD for error handling and edge cases
- **Focus**: User experience and API contracts

### Application Layer (Services)
- **Primary**: TDD for business logic
- **Secondary**: BDD for workflow validation
- **Focus**: Orchestration and business rules

### Domain Layer (Core Business Logic)
- **Primary**: TDD for pure business logic
- **Secondary**: Example-based testing
- **Focus**: Correctness and business invariants

### Infrastructure Layer (Adapters)
- **Primary**: TDD for implementation details
- **Secondary**: Integration tests
- **Focus**: External system interaction

## COMMIT GUIDELINES FOR HYBRID APPROACH

### Commit Types
- `feat(bdd):` - New BDD scenario or feature specification
- `feat(tdd):` - New TDD implementation or unit test
- `test(e2e):` - End-to-end BDD tests
- `test(unit):` - Unit tests and TDD cycles
- `refactor(bdd):` - Refactoring BDD scenarios
- `refactor(tdd):` - Refactoring implementation with TDD

### Commit Message Format
```
<type>(<scope>): <description>

[optional body explaining the BDD scenario or TDD cycle]

[optional footer with test results or metrics]
```

## TESTING PATTERNS

### BDD Test Structure
```go
// Given-When-Then structure
Describe("User Authentication", func() {
    Context("When user logs in with valid credentials", func() {
        It("should grant access to protected resources", func() {
            // Given: User has valid credentials
            // When: User attempts to log in
            // Then: User gains access to protected resources
        })
    })
})
```

### TDD Test Structure
```go
// Arrange-Act-Assert structure
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    service := NewUserService()
    userData := UserData{Name: "Test User"}
    
    // Act
    user, err := service.CreateUser(userData)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "Test User", user.Name)
}
```

## BEST PRACTICES

1. **Start with BDD for new features**
   - Define user behavior first
   - Create executable specifications
   - Use as acceptance criteria

2. **Use TDD for implementation details**
   - Drive internal API design
   - Ensure comprehensive edge case coverage
   - Maintain fast feedback loops

3. **Maintain clear separation**
   - BDD tests should not know about implementation
   - TDD tests should focus on specific units
   - Integration tests bridge the gap

4. **Keep tests independent**
   - Each test should be isolated
   - Use proper setup/teardown
   - Avoid test interdependencies

5. **Focus on readability**
   - Tests are documentation
   - Use descriptive names
   - Keep test logic simple

## ANTI-PATTERNS TO AVOID

1. **BDD tests that test implementation details**
   - Don't test internal API calls in BDD scenarios
   - Focus on user-observable behavior

2. **TDD tests that test multiple units**
   - Keep unit tests focused on single responsibility
   - Use integration tests for multi-unit scenarios

3. **Skipping the Red step**
   - Always see tests fail for the right reason
   - Ensures tests are actually testing something

4. **Writing tests after implementation**
   - Tests should drive design, not validate it
   - Post-hoc tests miss design benefits

5. **Ignoring refactoring**
   - Clean code is as important as working code
   - Refactor both production and test code

## MEASURING SUCCESS

### BDD Metrics
- **Scenario Coverage**: Percentage of user stories covered by BDD scenarios
- **Living Documentation**: How well BDD scenarios serve as documentation
- **Stakeholder Engagement**: Business stakeholder participation in scenario review

### TDD Metrics
- **Unit Test Coverage**: Percentage of code covered by unit tests
- **Test Speed**: Time to run full test suite
- **Defect Density**: Number of bugs per unit of code

### Integration Metrics
- **Feature Completion**: Time from BDD scenario to working feature
- **Test Pyramid Health**: Ratio of unit/integration/E2E tests
- **Refactoring Confidence**: Ability to safely change code

## CONTINUOUS IMPROVEMENT

1. **Regular retrospectives** on test effectiveness
2. **Refactor test code** as production code
3. **Update methodology** based on team learnings
4. **Share patterns** in common-patterns.md
5. **Monitor test execution time** and optimize slow tests

## NEXT STEPS

- See `.claude/tdd-bdd-implementation-guide.md` for practical implementation patterns
- See `.claude/tdd-bdd-examples.md` for concrete implementation examples
- See `.claude/common-patterns.md` for shared testing patterns and templates

This methodology provides the foundation for building high-quality, maintainable software that meets both user needs and technical requirements.
