-- Add performance indexes for reviews table
-- Composite indexes based on actual query patterns for optimal performance
ALTER TABLE `reviews`
  ADD INDEX `idx_reviews_spot_created` (`spot_id`, `created_at` DESC),
  ADD INDEX `idx_reviews_user_created` (`user_id`, `created_at` DESC),
  ADD INDEX `idx_reviews_spot_rating` (`spot_id`, `rating`),
  ADD INDEX `idx_reviews_rating_spot` (`rating`, `spot_id`)
  ALGORITHM=INPLACE, LOCK=NONE;