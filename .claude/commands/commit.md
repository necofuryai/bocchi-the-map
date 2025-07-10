## Git & Commit Standards

- **Commit Format**: Use conventional commit format with scopes: `type(scope): description`
- **No Auto-Signatures**: Do NOT include Claude Code attribution or co-authorship signatures in commit messages
- **Clean History**: Keep commit messages concise and focused on the actual changes
- **English Only**: All commit messages must be written in English

### Conventional Commit Format with Scopes

**Format**: `type(scope): description`

**Examples**:
```bash
feat(api): add user authentication endpoint
fix(web): resolve header styling issue
docs(claude): update development guidelines
chore(deps): update Go dependencies
refactor(scripts): consolidate test scripts
```

**Commit Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `chore`: Maintenance tasks, dependency updates
- `refactor`: Code refactoring without feature changes
- `test`: Adding or updating tests
- `ci`: CI/CD pipeline changes
- `perf`: Performance improvements

**Recommended Scopes for This Project**:
- `api`: Go backend changes
- `web`: Next.js frontend changes
- `infra`: Terraform infrastructure
- `scripts`: Build/test scripts
- `docs`: Documentation files (README, guides)
- `claude`: CLAUDE.md updates
- `deps`: Dependency updates
- `config`: Configuration files
- `auth`: Authentication-related changes
- `db`: Database migrations and changes

## Linus Torvalds' Ideal Commit Message Principles

For non-trivial changes, follow Linus Torvalds' commit message principles:

### 1. **Explain WHY, not WHAT**
The code shows WHAT changed. The commit message should explain:
- **Why** was this change necessary?
- What **problem** did it solve?
- What **effect** does it have?

### 2. **Structure for Complex Changes**
```
Short (50 chars or less) summary line

More detailed explanatory text, if necessary. Wrap it to
about 72 characters or so. The blank line separating the
summary from the body is critical.

Explain the problem that this commit is solving. Focus on
why you are making this change as opposed to how (the code
explains that). Are there side effects or other unintuitive
consequences of this change? Here's the place to explain them.

Further paragraphs come after blank lines.

- Bullet points are okay, too
- Typically a hyphen or asterisk is used for the bullet,
  preceded by a single space
```

### 3. **Good vs Bad Examples**

**❌ Bad:** `fix: update code`  
**✅ Good:** `fix(api): prevent race condition in user session handling`

**❌ Bad:** `feat: add new feature`  
**✅ Good:** `feat(web): add solo-friendly rating to improve UX for target users`

### 4. **When to Apply These Principles**

- **Always apply** for: Features, bug fixes, refactoring, performance improvements
- **Can be simplified** for: Typo fixes, dependency updates, formatting changes
- **Rule of thumb**: If you need to think about the change, explain your thinking

## Commit Workflow

When using this command file, follow these steps:

1. **Review Changes**: Check all modified files with `git status` and `git diff`
2. **Apply Standards**: Follow the Git & Commit Standards above when creating commits
3. **Stage Files**: Use `git add` to stage appropriate files for each commit
4. **Create Commit**: Use the conventional commit format with appropriate type and scope

## Multiple File Commits

When multiple files have been updated:

1. **Group by Feature/Scope**: Group related changes together
2. **Split Logical Changes**: Create separate commits for:
   - Different features or bug fixes
   - Different layers (frontend/backend/infrastructure)
   - Documentation updates vs code changes
   - Dependency updates vs feature changes

3. **Commit Order**: Consider dependencies between changes
   - Infrastructure changes before application changes
   - Core changes before dependent features
   - Breaking changes clearly marked

Example workflow for multiple files:
```bash
# First commit: Backend API changes
git add api/*.go
git commit -m "feat(api): add user profile endpoint"

# Second commit: Frontend integration
git add web/src/services/* web/src/components/*
git commit -m "feat(web): integrate user profile API"

# Third commit: Documentation
git add docs/*.md README.md
git commit -m "docs: update API documentation for user profiles"
```
