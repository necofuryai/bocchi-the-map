# Implementor

Split complex tasks into sequential steps, where each step can contain multiple parallel subtasks.

## Process

1. **Initial Analysis**
   - First, analyze the entire task to understand scope and requirements
   - Identify dependencies and execution order
   - Plan sequential steps based on dependencies

2. **Step Planning**
   - Break down into 2-4 sequential steps
   - Each step can contain multiple parallel subtasks
   - Define what context from previous steps is needed

3. **Step-by-Step Execution**
   - Execute all subtasks within a step in parallel
   - Wait for all subtasks in current step to complete
   - Pass relevant results to next step
   - Request concise summaries (100-200 words) from each subtask

4. **Step Review and Adaptation**
   - After each step completion, review results
   - Validate if remaining steps are still appropriate
   - Adjust next steps based on discoveries
   - Add, remove, or modify subtasks as needed

5. **Progressive Aggregation**
   - Synthesize results from completed step
   - Use synthesized results as context for next step
   - Build comprehensive understanding progressively
   - Maintain flexibility to adapt plan

## Example Usage

When given "analyze test lint and commit":

**Step 1: Initial Analysis** (1 subtask)
- Analyze project structure to understand test/lint setup

**Step 2: Quality Checks** (parallel subtasks)
- Run tests and capture results
- Run linting and type checking
- Check git status and changes

**Step 3: Fix Issues** (parallel subtasks, using Step 2 results)
- Fix linting errors found in Step 2
- Fix type errors found in Step 2
- Prepare commit message based on changes
*Review: If no errors found in Step 2, skip fixes and proceed to commit*

**Step 4: Final Validation** (parallel subtasks)
- Re-run tests to ensure fixes work
- Re-run lint to verify all issues resolved
- Create commit with verified changes
*Review: If Step 3 had no fixes, simplify to just creating commit*

## Key Benefits

- **Sequential Logic**: Steps execute in order, allowing later steps to use earlier results
- **Parallel Efficiency**: Within each step, independent tasks run simultaneously
- **Memory Optimization**: Each subtask gets minimal context, preventing overflow
- **Progressive Understanding**: Build knowledge incrementally across steps
- **Clear Dependencies**: Explicit flow from analysis → execution → validation

## TDD+BDD Integration

When implementing features, integrate the TDD+BDD hybrid methodology into the sequential step framework:

### BDD-Driven Step Planning
- **Step 1: BDD Analysis** - Define user behavior with Given-When-Then scenarios
- **Step 2: TDD Implementation** - Drive internal components with Red-Green-Refactor cycles
- **Step 3: Integration Testing** - Validate end-to-end behavior matches BDD scenarios
- **Step 4: Documentation Update** - Update relevant `.claude/` files with new patterns

### Layer-Specific Testing Within Steps
- **Interface Layer**: Use BDD scenarios to define expected behavior
- **Domain Layer**: Apply TDD for business logic and validation
- **Infrastructure Layer**: Use TDD for database and external service integration

### Test-First Execution Pattern
Within each step, always follow:
1. **Red**: Write failing test first
2. **Green**: Implement minimal code to make test pass  
3. **Refactor**: Improve code while keeping tests green
4. **Validate**: Ensure BDD scenarios still pass

### Reference Documentation Integration
- Use `.claude/tdd-bdd-methodology.md` for theoretical guidance
- Apply patterns from `.claude/tdd-bdd-implementation-guide.md`
- Reference examples from `.claude/tdd-bdd-examples.md`
- Update `.claude/common-patterns.md` with new testing patterns

## Implementation Notes

- Always start with a single analysis task to understand the full scope
- Group related parallel tasks within the same step
- Pass only essential findings between steps (summaries, not full output)
- Use TodoWrite to track both steps and subtasks for visibility
- **Apply TDD+BDD methodology within each step's execution**
- After each step, explicitly reconsider the plan:
  - Are the next steps still relevant?
  - Did we discover something that requires new tasks?
  - Can we skip or simplify upcoming steps?
  - Should we add new validation steps?
  - Do the BDD scenarios still match the implemented behavior?

## Adaptive Planning Example

```
Initial Plan: Step 1 → Step 2 → Step 3 → Step 4

After Step 2: "No errors found in tests or linting"
Adapted Plan: Step 1 → Step 2 → Skip Step 3 → Simplified Step 4 (just commit)

After Step 2: "Found critical architectural issue"
Adapted Plan: Step 1 → Step 2 → New Step 2.5 (analyze architecture) → Modified Step 3
```
