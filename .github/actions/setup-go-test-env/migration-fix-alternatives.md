# Migration File Selection Fix - Alternative Approaches

## Problem
The original approach in `.github/actions/setup-go-test-env/action.yml` lines 89-100 copied all `.sql` files from the migrations directory, which included:
- Utility files like `explain_index_performance.sql`
- Documentation files like `token_blacklist_index_analysis.md`
- Risk of accidentally including production subdirectory files

## Solution 1: Pattern-Based File Selection (IMPLEMENTED)
```bash
# Create temporary directory with only numbered migration files (exclude production subdirectory and utility files)
mkdir -p migrations_test
find migrations -maxdepth 1 -name "[0-9][0-9][0-9][0-9][0-9][0-9]_*.sql" -exec cp {} migrations_test/ \;
migrate -path migrations_test -database "mysql://$TEST_DATABASE_URL" up
rm -rf migrations_test
```

**Benefits:**
- ✅ Only includes properly numbered migration files (000001_*.sql, 000002_*.sql, etc.)
- ✅ Excludes utility and documentation files
- ✅ Excludes production subdirectory
- ✅ Clean and predictable test environment
- ✅ Maintains idempotency

## Solution 2: Using migrate -ignore-files Option (ALTERNATIVE)
```bash
# Alternative approach using migrate's built-in ignore functionality
migrate -path migrations -database "mysql://$TEST_DATABASE_URL" -ignore-files "production/*,*.md,explain_*.sql" up
```

**Benefits:**
- ✅ No temporary directory needed
- ✅ Built-in migrate tool feature
- ✅ Explicit ignore patterns

**Drawbacks:**
- ⚠️ Requires newer version of migrate tool
- ⚠️ May not be available in all migrate versions
- ⚠️ Less explicit control over included files

## Solution 3: Explicit File List (VERBOSE BUT SAFE)
```bash
# Most explicit approach - copy only specific migration files
mkdir -p migrations_test
for file in migrations/[0-9][0-9][0-9][0-9][0-9][0-9]_*.sql; do
  if [ -f "$file" ]; then
    cp "$file" migrations_test/
  fi
done
migrate -path migrations_test -database "mysql://$TEST_DATABASE_URL" up
rm -rf migrations_test
```

## Recommendation
**Solution 1 (Pattern-Based)** is recommended because:
- Uses standard Unix tools (find, cp)
- Clear and explicit pattern matching
- No dependency on specific migrate tool versions
- Proven approach across different environments
- Easy to understand and maintain

## Testing the Fix
To verify the fix works correctly:
```bash
# Test what files are selected
find migrations -maxdepth 1 -name "[0-9][0-9][0-9][0-9][0-9][0-9]_*.sql"

# Should return only:
# migrations/000001_initial_schema.down.sql
# migrations/000001_initial_schema.up.sql
# migrations/000002_token_blacklist.down.sql
# migrations/000002_token_blacklist.up.sql
# ... (numbered migration files only)
```