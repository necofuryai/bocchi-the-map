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

## Solution 2: Manual File Management (ALTERNATIVE)
```bash
# Alternative approach: manually manage unwanted files before running migrations
# Move or delete unwanted files temporarily
mkdir -p backup_temp
mv migrations/explain_*.sql backup_temp/ 2>/dev/null || true
mv migrations/*.md backup_temp/ 2>/dev/null || true
mv migrations/production backup_temp/ 2>/dev/null || true

# Run migrations
migrate -path migrations -database "mysql://$TEST_DATABASE_URL" up

# Restore files
mv backup_temp/* migrations/ 2>/dev/null || true
mv backup_temp/production migrations/ 2>/dev/null || true
rmdir backup_temp 2>/dev/null || true
```

**Benefits:**
- ✅ Works with any version of migrate tool
- ✅ Full control over which files are included
- ✅ No dependency on specific migrate features

**Drawbacks:**
- ⚠️ More complex file management
- ⚠️ Risk of data loss if backup fails
- ⚠️ Requires careful error handling

## Solution 3: Using Goose Migration Tool (TOOL ALTERNATIVE)
```bash
# Alternative migration tool that supports file exclusion patterns
# Install goose: go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir migrations mysql "$TEST_DATABASE_URL" up
```

**Benefits:**
- ✅ Modern migration tool with better file handling
- ✅ Built-in support for ignoring non-migration files
- ✅ Better error handling and rollback capabilities
- ✅ Active development and community support

**Drawbacks:**
- ⚠️ Requires changing migration tool
- ⚠️ May need migration file format adjustments
- ⚠️ Additional dependency

## Solution 4: Explicit File List (VERBOSE BUT SAFE)
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
# 
# Total files found: 4 (or the exact number of numbered migration files)
```