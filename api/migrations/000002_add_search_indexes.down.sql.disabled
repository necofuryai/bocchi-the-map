-- Remove search performance indexes
DROP INDEX `idx_spots_name_address` ON `spots`;
DROP INDEX `idx_spots_name` ON `spots`;
DROP INDEX `idx_spots_address` ON `spots`;
DROP INDEX `idx_spots_category_rating` ON `spots`;
DROP INDEX `idx_spots_country_rating` ON `spots`;

-- Restore original location index
DROP INDEX `idx_spots_location` ON `spots`;
CREATE INDEX `idx_location` ON `spots`(`latitude`, `longitude`);