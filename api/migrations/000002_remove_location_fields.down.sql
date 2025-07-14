-- Restore location-related fields and index
ALTER TABLE `spots` ADD COLUMN `latitude` DECIMAL(10, 8) NOT NULL AFTER `name`;
ALTER TABLE `spots` ADD COLUMN `longitude` DECIMAL(11, 8) NOT NULL AFTER `latitude`;
ALTER TABLE `spots` ADD INDEX `idx_location` (`latitude`, `longitude`);