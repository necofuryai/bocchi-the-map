# Database URL Consistency Fix

## Problem Fixed
The `TEST_DATABASE_URL` environment variable was duplicated and hardcoded in multiple places, leading to:
- ❌ Code duplication and maintenance burden
- ❌ Potential inconsistencies between action and workflow
- ❌ Risk of configuration drift over time

## Changes Made

### 1. Composite Action (`.github/actions/setup-go-test-env/action.yml`)
✅ **Added environment variable export**:
```bash
# Export TEST_DATABASE_URL for use in subsequent workflow steps
echo "TEST_DATABASE_URL=$TEST_DATABASE_URL" >> $GITHUB_ENV
```

**Benefits:**
- Single source of truth for database URL construction
- Consistent password handling via `${{ inputs.mysql-root-password }}`
- Environment variable available to all subsequent workflow steps

### 2. Workflow Configuration (`.github/workflows/bdd-tests.yml`)

#### BDD Tests Job (lines 51-56)
❌ **Removed duplicated TEST_DATABASE_URL**:
```yaml
# BEFORE
env:
  TEST_DATABASE_URL: >
    root:${{ secrets.MYSQL_ROOT_PASSWORD }}@tcp(localhost:3306)/bocchi_test?parseTime=true&multiStatements=true
  CGO_ENABLED: 1

# AFTER  
env:
  CGO_ENABLED: 1
```

#### Integration Tests Job (lines 108-114)
❌ **Removed duplicated TEST_DATABASE_URL**:
```yaml
# BEFORE
env:
  TEST_DATABASE_URL: >
    root:${{ secrets.MYSQL_ROOT_PASSWORD }}@tcp(localhost:3306)/bocchi_test?parseTime=true&multiStatements=true
  CGO_ENABLED: 1

# AFTER
env:
  CGO_ENABLED: 1
```

## How It Works Now

1. **Setup Action** creates `TEST_DATABASE_URL` with secure password
2. **Environment Export** makes it available to all subsequent steps
3. **Test Steps** automatically inherit the consistent database URL
4. **No Duplication** - single source of truth maintained

## Flow Diagram

```
[Setup Action] 
    ↓ Generates TEST_DATABASE_URL
    ↓ Exports to $GITHUB_ENV
    ↓
[BDD Tests] ← Uses exported TEST_DATABASE_URL
    ↓
[Integration Tests] ← Uses exported TEST_DATABASE_URL
```

## Security Benefits

✅ **Centralized password handling** - only in composite action  
✅ **Consistent URL format** - no copy-paste errors  
✅ **Reduced maintenance** - change in one place  
✅ **No hardcoded credentials** - uses secrets properly  

## Testing

The workflow will now:
1. Generate `TEST_DATABASE_URL` once in the setup action
2. Export it for use in all subsequent steps  
3. Both BDD and integration tests will use the same URL
4. Password comes from `MYSQL_ROOT_PASSWORD` secret consistently

## Rollback (if needed)

If you need to temporarily rollback:
```yaml
# Add back to individual test steps - TEMPORARY ONLY
env:
  TEST_DATABASE_URL: "root:${{ secrets.MYSQL_ROOT_PASSWORD }}@tcp(localhost:3306)/bocchi_test?parseTime=true&multiStatements=true"
  CGO_ENABLED: 1
```

**Note**: Always maintain the single source of truth approach for production!