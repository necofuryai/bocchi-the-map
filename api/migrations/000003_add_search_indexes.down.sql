-- Remove search performance indexes
DROP INDEX IF EXISTS `idx_spots_name_address` ON `spots`;
DROP INDEX IF EXISTS `idx_spots_address` ON `spots`;
DROP INDEX IF EXISTS `idx_spots_category_rating` ON `spots`;
DROP INDEX IF EXISTS `idx_spots_country_rating` ON `spots`;

-- Keep original location index (no restoration needed)
-- Original idx_location remains unchanged from initial schema