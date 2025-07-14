-- Rollback MVP schema for bocchi-the-map

-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS `reviews`;
DROP TABLE IF EXISTS `spots`;
DROP TABLE IF EXISTS `users`;