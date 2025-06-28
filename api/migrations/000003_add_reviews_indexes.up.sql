-- Add performance indexes for reviews table
CREATE INDEX `idx_reviews_spot_id` ON `reviews` (`spot_id`);
CREATE INDEX `idx_reviews_user_id` ON `reviews` (`user_id`);
CREATE INDEX `idx_reviews_rating` ON `reviews` (`rating`);
CREATE INDEX `idx_reviews_created_at` ON `reviews` (`created_at`);