-- Re-add anonymous_id column to users table  
ALTER TABLE users ADD COLUMN IF NOT EXISTS anonymous_id VARCHAR(36);