# Pull Request

## 📋 Description
<!-- Briefly describe what this PR does -->

## 🔄 Type of Change
<!-- Mark with 'x' the type that applies -->
- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 📚 Documentation only changes
- [ ] 🔧 Code refactoring (no functional changes, no api changes)
- [ ] ⚡ Performance improvements
- [ ] 🧪 Adding or updating tests
- [ ] 🚀 Infrastructure/CI/CD changes

## 🧪 Testing
<!-- Describe the tests you ran to verify your changes -->
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] E2E tests pass (if applicable)
- [ ] Manual testing completed

## 📚 Documentation Checklist
<!-- Automated check will suggest documentation updates, but please review manually -->

### 📝 Core Documentation
- [ ] Updated `README.md` if setup/installation instructions changed
- [ ] Updated `docs/IMPLEMENTATION_LOG.md` for significant features or milestones
- [ ] Updated `CLAUDE.md` if AI development guidance changed

### 🔧 Component-Specific Documentation
- [ ] Updated `api/README.md` if API endpoints or database schema changed
- [ ] Updated `web/README.md` if frontend setup or component structure changed  
- [ ] Updated `infra/README.md` if deployment procedures or infrastructure changed

### 📊 Technical Documentation  
- [ ] Added/updated code comments for complex logic
- [ ] Updated API documentation (if API changes)
- [ ] Updated type definitions or schemas
- [ ] Added performance benchmarks (if performance-related changes)
  <!-- Run benchmarks using `make benchmark` in the api/ directory or check docs/BENCHMARKS.md for detailed instructions -->

### 🎯 Quality Assurance
- [ ] All commit messages follow [conventional commits](https://www.conventionalcommits.org/)
- [ ] Code follows project style guidelines
- [ ] No hardcoded secrets or sensitive information
- [ ] Security vulnerability scan completed (if security-related changes)
- [ ] Breaking changes are clearly documented

## 🔗 Related Issues
<!-- Link related issues or feature requests -->
Closes #

## 📸 Screenshots (if applicable)
<!-- Add screenshots for UI changes -->

## 🎯 Review Focus Areas
<!-- Help reviewers by highlighting specific areas that need attention -->
- [ ] Security considerations
- [ ] Performance impact
- [ ] Breaking changes
- [ ] Database migration safety (including rollback testing for schema changes)
- [ ] Error handling

## 🤖 Automated Checks
<!-- These will be verified automatically -->
- [ ] All CI/CD checks passing
- [ ] Code coverage maintained or improved
- [ ] No new linting errors
- [ ] Documentation sync check passed

---

**💡 Tip:** The Documentation Sync Check workflow will automatically suggest documentation updates based on your changes. Please review the suggestions and update accordingly.

**🎉 Thank you for contributing to Bocchi The Map!**