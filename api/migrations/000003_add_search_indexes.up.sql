-- Add indexes for improved search performance on spots table
-- Composite index for name and address LIKE searches
CREATE INDEX `idx_spots_name_address` ON `spots`(`name`, `address`(255));

-- Individual index for address-only searches
CREATE INDEX `idx_spots_address` ON `spots`(`address`(255));

-- Composite index for filtering and sorting
CREATE INDEX `idx_spots_category_rating` ON `spots`(`category`, `average_rating` DESC, `review_count` DESC);
CREATE INDEX `idx_spots_country_rating` ON `spots`(`country_code`, `average_rating` DESC, `review_count` DESC);

-- Regular composite index for location-based queries
-- Note: SPATIAL INDEX requires POINT/GEOMETRY columns, using regular BTREE index for DECIMAL columns
-- Keep original idx_location index as-is (already exists from initial schema)
-- No additional location index needed since idx_location already provides this functionality