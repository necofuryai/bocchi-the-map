-- Add indexes for improved search performance on spots table
-- Composite index for name and address LIKE searches
CREATE INDEX `idx_spots_name_address` ON `spots`(`name`, `address`(255));

-- Individual indexes for better query optimization
CREATE INDEX `idx_spots_name` ON `spots`(`name`);
CREATE INDEX `idx_spots_address` ON `spots`(`address`(255));

-- Composite index for filtering and sorting
CREATE INDEX `idx_spots_category_rating` ON `spots`(`category`, `average_rating` DESC, `review_count` DESC);
CREATE INDEX `idx_spots_country_rating` ON `spots`(`country_code`, `average_rating` DESC, `review_count` DESC);

-- Spatial index for better location-based queries (if MySQL 5.7+)
-- Note: This replaces the existing idx_location with a spatial index for better performance
DROP INDEX `idx_location` ON `spots`;
ALTER TABLE `spots` ADD SPATIAL INDEX `idx_spots_location` (`latitude`, `longitude`);