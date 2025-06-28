-- Re-add anonymous_id column to users table  
ALTER TABLE users ADD COLUMN anonymous_id VARCHAR(36);