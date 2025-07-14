-- Remove location-related fields and index
ALTER TABLE `spots` DROP INDEX `idx_location`;
ALTER TABLE `spots` DROP COLUMN `latitude`;
ALTER TABLE `spots` DROP COLUMN `longitude`;