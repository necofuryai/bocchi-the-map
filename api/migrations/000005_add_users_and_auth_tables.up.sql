-- Add users table for Auth0 integration
CREATE TABLE `users` (
    `id` VARCHAR(36) PRIMARY KEY,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `name` VARCHAR(255),
    `nickname` VARCHAR(100),
    `picture` TEXT,
    `provider` VARCHAR(50) NOT NULL DEFAULT 'auth0',
    `provider_id` VARCHAR(255) NOT NULL,
    `email_verified` BOOLEAN NOT NULL DEFAULT FALSE,
    `preferences` JSON,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY `idx_provider_user` (`provider`, `provider_id`),
    INDEX `idx_email` (`email`),
    INDEX `idx_provider_id` (`provider_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Add token blacklist table for logout and security
CREATE TABLE `token_blacklist` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `jti` VARCHAR(255) NOT NULL UNIQUE,
    `token_type` ENUM('access', 'refresh') NOT NULL DEFAULT 'access',
    `expires_at` TIMESTAMP NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    INDEX `idx_jti` (`jti`),
    INDEX `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Update reviews table to use user_id instead of reviewer_name
ALTER TABLE `reviews` 
ADD COLUMN `user_id` VARCHAR(36) AFTER `spot_id`,
ADD CONSTRAINT `fk_reviews_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE;

-- Create index for better query performance
ALTER TABLE `reviews` ADD INDEX `idx_user_id` (`user_id`);
ALTER TABLE `reviews` ADD INDEX `idx_spot_user` (`spot_id`, `user_id`);