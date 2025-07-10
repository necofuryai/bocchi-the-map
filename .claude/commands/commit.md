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
