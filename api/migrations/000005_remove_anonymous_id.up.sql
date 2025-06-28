-- Remove anonymous_id column from users table
ALTER TABLE users DROP COLUMN IF EXISTS anonymous_id;