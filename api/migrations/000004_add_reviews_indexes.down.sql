-- Drop performance indexes for reviews table
DROP INDEX IF EXISTS `idx_reviews_spot_created` ON `reviews`;
DROP INDEX IF EXISTS `idx_reviews_user_created` ON `reviews`;
DROP INDEX IF EXISTS `idx_reviews_spot_rating` ON `reviews`;
DROP INDEX IF EXISTS `idx_reviews_rating_spot` ON `reviews`;