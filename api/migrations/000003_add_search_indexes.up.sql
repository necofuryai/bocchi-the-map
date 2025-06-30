-- Add indexes for improved search performance on spots table
-- DEPENDENCY: This migration requires idx_location index from 000001_initial_schema.up.sql

-- Create a temporary procedure to safely check for required dependencies
DELIMITER $$
CREATE PROCEDURE CheckDependencies()
BEGIN
    DECLARE index_count INT DEFAULT 0;
    
    -- Verify that the required idx_location index exists before proceeding
    SELECT COUNT(*) INTO index_count
    FROM information_schema.statistics 
    WHERE table_schema = database() 
      AND table_name = 'spots' 
      AND index_name = 'idx_location'
      AND column_name IN ('latitude', 'longitude');
    
    -- Abort migration if the required index does not exist
    IF index_count = 0 THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = 'Migration aborted: Required index idx_location not found on spots table. Please ensure migration 000001_initial_schema.up.sql has been applied first.';
    END IF;
END$$
DELIMITER ;

-- Execute the dependency check
CALL CheckDependencies();

-- Clean up the temporary procedure
DROP PROCEDURE CheckDependencies;

-- Composite index for name and address LIKE searches
CREATE INDEX `idx_spots_name_address` ON `spots`(`name`, `address`(255));

-- Individual index for name-only searches
CREATE INDEX `idx_spots_name` ON `spots`(`name`);

-- Individual index for address-only searches
CREATE INDEX `idx_spots_address` ON `spots`(`address`(255));

-- Composite index for filtering and sorting
CREATE INDEX `idx_spots_category_rating` ON `spots`(`category`, `average_rating` DESC, `review_count` DESC);
CREATE INDEX `idx_spots_country_rating` ON `spots`(`country_code`, `average_rating` DESC, `review_count` DESC);

-- Regular composite index for location-based queries
-- Note: SPATIAL INDEX requires POINT/GEOMETRY columns, using regular BTREE index for DECIMAL columns
-- DEPENDENCY CONFIRMED: idx_location index exists from 000001_initial_schema.up.sql (line 32)
-- This migration relies on the existing idx_location index for location-based functionality
-- No additional location index needed since idx_location already provides this functionality