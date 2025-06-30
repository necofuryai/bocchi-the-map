# Token Blacklist Index Analysis

## Current Index Configuration
```sql
-- File: api/migrations/000005_add_token_blacklist_composite_index.up.sql:4-5
ALTER TABLE token_blacklist
  ADD INDEX idx_token_blacklist_jti_expires (jti, expires_at);
```

## Column Specifications
```sql
-- File: api/migrations/000002_token_blacklist.up.sql:4
jti VARCHAR(255) NOT NULL UNIQUE, -- JWT ID for identifying tokens
expires_at TIMESTAMP NOT NULL,    -- When the token expires
```

## Index Size Calculation
- **Charset**: utf8mb4 (4 bytes per character maximum)
- **jti column**: VARCHAR(255) × 4 bytes + 2 bytes length = 1,022 bytes maximum
- **expires_at column**: TIMESTAMP = 8 bytes
- **Total index size**: 1,030 bytes
- **MySQL limit**: 3,072 bytes for InnoDB
- **Result**: ✅ **SAFE** - Well within limits

**Note**: VARCHAR fields in MySQL InnoDB require 1-2 additional bytes to store the length of the variable-length data. For VARCHAR(255), 2 bytes are used for length storage.

## Query Pattern Analysis

### 1. Primary Query (Hot Path)
**Pattern**: `WHERE jti = ? AND expires_at > NOW()`
**Frequency**: Every authenticated API request
**Files**: 
- `api/queries/token_blacklist.sql:6-7`
- `api/pkg/auth/middleware.go:94, 162`

**Index Usage**: ✅ **OPTIMAL**
- Uses composite index `(jti, expires_at)` efficiently
- Exact match on `jti` (first column)
- Range condition on `expires_at` (second column)

### 2. Cleanup Query
**Pattern**: `WHERE expires_at < NOW() - INTERVAL 24 HOUR`
**Frequency**: Hourly (scheduled cleanup)
**Files**:
- `api/queries/token_blacklist.sql:18-20`
- `api/scripts/token_cleanup_event.sql:11-13`

**Index Usage**: ⚠️ **SUBOPTIMAL but ACCEPTABLE**
- Cannot use composite index efficiently (no jti filter)
- Would benefit from separate `expires_at` index
- However, cleanup frequency is low (hourly vs. per-request)

## Recommendations

### ✅ **KEEP Current Index Order**
The composite index `(jti, expires_at)` is **correctly ordered** for the primary use case:
- Hot path query performance is critical
- Index covers both equality and range conditions optimally

### 🤔 **Optional: Add Separate `expires_at` Index**
Consider adding for cleanup query optimization:
```sql
ADD INDEX idx_token_blacklist_expires (expires_at);
```

**Trade-offs**:
- **Pros**: Faster cleanup queries, better for large blacklist tables
- **Cons**: Additional storage overhead, slower inserts
- **Decision**: Not critical due to low cleanup frequency

### 📋 **Final Verification**
Before applying in production, run EXPLAIN on both queries:
```sql
-- Verify primary query uses composite index
EXPLAIN SELECT COUNT(*) > 0 FROM token_blacklist 
WHERE jti = 'test-jti' AND expires_at > NOW();

-- Check cleanup query performance
EXPLAIN DELETE FROM token_blacklist 
WHERE expires_at < NOW() - INTERVAL 24 HOUR LIMIT 1000;
```

### 📊 **Performance Analysis with Execution Metrics**
For comprehensive performance insights, run EXPLAIN ANALYZE to capture actual runtime metrics:
```sql
-- Analyze primary query with actual execution time and row processing
EXPLAIN ANALYZE SELECT COUNT(*) > 0 FROM token_blacklist 
WHERE jti = 'test-jti' AND expires_at > NOW();

-- Analyze cleanup query with actual execution time and affected rows
EXPLAIN ANALYZE DELETE FROM token_blacklist 
WHERE expires_at < NOW() - INTERVAL 24 HOUR LIMIT 1000;
```

**Key metrics to monitor:**
- **Execution time**: Actual time vs. estimated cost
- **Rows examined**: Actual vs. estimated row count
- **Index usage**: Confirm composite index utilization
- **Buffer pool**: Memory efficiency during operations

## Conclusion
**✅ Current index configuration is OPTIMAL** for the primary use case and safe regarding size limits. The composite index `(jti, expires_at)` correctly prioritizes the hot path query performance over cleanup query performance.
