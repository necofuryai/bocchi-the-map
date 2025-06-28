-- Drop performance indexes for reviews table
DROP INDEX `idx_reviews_spot_id` ON `reviews`;
DROP INDEX `idx_reviews_user_id` ON `reviews`;
DROP INDEX `idx_reviews_rating` ON `reviews`;
DROP INDEX `idx_reviews_created_at` ON `reviews`;