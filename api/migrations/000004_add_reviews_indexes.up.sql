-- Add performance indexes for reviews table
-- Note: Production version with ALGORITHM=INPLACE, LOCK=NONE is in migrations/production/
ALTER TABLE `reviews`
  ADD INDEX `idx_reviews_spot_id` (`spot_id`),
  ADD INDEX `idx_reviews_user_id` (`user_id`),
  ADD INDEX `idx_reviews_rating` (`rating`),
  ADD INDEX `idx_reviews_created_at` (`created_at`);