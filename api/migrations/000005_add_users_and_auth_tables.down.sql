-- Reverse the changes from 000005_add_users_and_auth_tables.up.sql

-- Remove foreign key constraint and user_id column from reviews table
ALTER TABLE `reviews` 
DROP FOREIGN KEY `fk_reviews_user_id`,
DROP INDEX `idx_user_id`,
DROP INDEX `idx_spot_user`,
DROP COLUMN `user_id`;

-- Drop token blacklist table
DROP TABLE IF EXISTS `token_blacklist`;

-- Drop users table
DROP TABLE IF EXISTS `users`;