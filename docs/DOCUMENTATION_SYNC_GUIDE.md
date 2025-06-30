# 📚 Documentation Sync Guide

This guide helps developers maintain up-to-date documentation when making changes to the Bocchi The Map project.

## 🎯 Overview

Bocchi The Map uses an automated documentation sync system to ensure that implementation changes are properly reflected in documentation. This guide explains how to work with this system effectively.

## 🤖 Automated Systems

### 1. **Git Hooks (Local Development)**

#### Pre-commit Hook
- **Location**: `.husky/pre-commit`
- **Function**: Analyzes staged changes and suggests documentation updates
- **Triggers**: Before every commit
- **Suggestions**:
  - API changes → Update `api/README.md`
  - Frontend changes → Update `web/README.md`
  - Infrastructure changes → Update `infra/README.md`
  - Feature commits → Update `docs/IMPLEMENTATION_LOG.md`

#### Commitlint Integration
- **Location**: `.commitlintrc.js`
- **Function**: Enforces conventional commit messages
- **Documentation-related types**:
  - `docs:` - Documentation only changes
  - `feat:` - New features (triggers documentation suggestions)
  - `improve:` - Improvements to existing features

### 2. **GitHub Actions (CI/CD)**

#### Documentation Sync Check
- **Workflow**: `.github/workflows/doc-sync-check.yml`
- **Triggers**: Pull requests to main/develop branches
- **Function**: 
  - Analyzes PR changes for documentation requirements
  - Posts automated comments with specific suggestions
  - Warns about critical changes without documentation updates

#### Milestone Recorder
- **Workflow**: `.github/workflows/milestone-recorder.yml`
- **Triggers**: Pushes to main/develop, merged PRs
- **Function**: 
  - Automatically records significant changes in `docs/IMPLEMENTATION_LOG.md`
  - Commits milestone updates directly to the repository

#### E2E Progress Tracker
- **Workflow**: `.github/workflows/e2e-progress-tracker.yml`
- **Triggers**: E2E test completion, weekly schedule
- **Function**: 
  - Tracks E2E test success rate progress
  - Records milestones when reaching 90%, 95%, or 100% success
  - Maintains `docs/e2e-test-progress.md`

## 📋 Documentation Update Checklist

### For New Features (`feat:` commits)

#### **Always Update:**
1. **`docs/IMPLEMENTATION_LOG.md`** - Add milestone entry
2. **Component-specific README** - Update relevant setup/usage instructions

#### **Consider Updating:**
- Main `README.md` if setup instructions changed
- `CLAUDE.md` if AI development guidance changed
- API documentation if new endpoints added
- Performance benchmarks if applicable

### For Bug Fixes (`fix:` commits)

#### **Update if applicable:**
- Component README if fix affects setup/usage
- Implementation log for significant fixes
- Troubleshooting sections

### For Infrastructure Changes (`ci:`, `build:`)

#### **Always Update:**
1. **`infra/README.md`** - Deployment procedures
2. **Main `README.md`** - Setup requirements

#### **Consider Updating:**
- CI/CD documentation
- Environment setup guides

## 🔧 Manual Documentation Update Process

### 1. **Identify Required Updates**
```bash
# Check what files changed
git diff --name-only

# Pre-commit hook will suggest documentation updates
git commit -m "feat: add user authentication system"
```

### 2. **Update Relevant Documentation**
```bash
# Update component documentation
vim api/README.md

# Update implementation log
vim docs/IMPLEMENTATION_LOG.md

# Commit documentation changes
git add docs/ api/README.md
git commit -m "docs: update documentation for user authentication system"
```

### 3. **Verify with Automated Checks**
- Create PR and review GitHub Actions suggestions
- Automated milestone recording will handle implementation log updates for significant changes

## 📊 Documentation Structure

### **Primary Documentation Files**
- **`README.md`** - Project overview, setup, quick start
- **`CLAUDE.md`** - AI development guidelines
- **`docs/IMPLEMENTATION_LOG.md`** - Implementation milestones and history

### **Component Documentation**
- **`api/README.md`** - Backend development guide
- **`web/README.md`** - Frontend development guide  
- **`infra/README.md`** - Infrastructure and deployment

### **Specialized Documentation**
- **`docs/e2e-test-progress.md`** - E2E test progress tracking (auto-generated)
- **`.claude/project-improvements.md`** - AI development knowledge base
- **`.claude/project-knowledge.md`** - Architecture patterns and decisions

## 🎯 Best Practices

### **Commit Message Patterns**
```bash
# Feature with documentation
feat: add user profile management
docs: update API documentation for profile endpoints

# Infrastructure change
ci: improve GitHub Actions performance  
docs: update deployment guide for new CI optimizations

# Bug fix with docs
fix: resolve authentication token expiry issue
docs: add troubleshooting section for token issues
```

### **PR Description Guidelines**
- Use the PR template checklist
- Mark documentation updates in the checklist
- Reference automated suggestions in PR comments
- Link related documentation files

### **Documentation Quality Standards**
1. **Accuracy**: Keep documentation in sync with implementation
2. **Completeness**: Cover setup, usage, and troubleshooting
3. **Clarity**: Use clear, concise language
4. **Examples**: Provide practical code examples
5. **Maintenance**: Regular reviews and updates

## 🔄 Automated vs Manual Updates

### **Automated (No action required)**
- ✅ Milestone recording for significant commits
- ✅ E2E test progress tracking
- ✅ PR documentation suggestions
- ✅ Pre-commit documentation reminders

### **Manual (Developer action required)**
- 📝 Component-specific README updates
- 📝 API documentation changes
- 📝 Setup instruction modifications
- 📝 Troubleshooting guide additions

## 🚨 When Documentation Updates are Critical

### **Always Required:**
- New API endpoints or breaking changes
- Changes to setup/installation procedures
- New dependencies or environment requirements
- Database schema modifications

### **Recommended:**
- New features or significant improvements
- Performance optimizations
- Security enhancements
- Bug fixes affecting user workflows

### **Optional:**
- Minor code refactoring
- Style/formatting changes
- Internal implementation changes without external impact

## 💡 Tips for Efficient Documentation

1. **Use the automated suggestions** - GitHub Actions will guide you
2. **Update documentation in the same PR** - Don't defer documentation updates
3. **Keep changes focused** - One feature/fix per PR with corresponding docs
4. **Review existing documentation** - Ensure consistency with project style
5. **Test documentation** - Verify that setup instructions actually work

## 🔍 Troubleshooting

### **Pre-commit hook not running**
```bash
# Reinstall Husky hooks
npm run prepare
chmod +x .husky/pre-commit
```

### **GitHub Actions not suggesting documentation**
- Check if PR includes changes to tracked file patterns
- Ensure PR is not in draft mode
- Verify workflow permissions in repository settings

### **Documentation out of sync**
- Use the automated milestone recorder to catch up
- Review recent commits for undocumented changes
- Run documentation sync check manually

---

**🎉 Remember**: Good documentation is a gift to your future self and other developers. The automated systems are here to help, but your judgment and attention to detail make the difference!

**💡 Pro tip**: When in doubt, over-document rather than under-document. The automated systems will help you maintain consistency and quality.