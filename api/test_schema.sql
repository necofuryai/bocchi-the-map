-- Complete test schema for bocchi-the-map

-- Create spots table
CREATE TABLE `spots` (
    `id` VARCHAR(36) PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `name_i18n` JSON,
    `latitude` DECIMAL(10, 8) NOT NULL,
    `longitude` DECIMAL(11, 8) NOT NULL,
    `category` VARCHAR(100) NOT NULL,
    `address` TEXT NOT NULL,
    `address_i18n` JSON,
    `country_code` CHAR(2) NOT NULL,
    `average_rating` DECIMAL(3, 1) NOT NULL DEFAULT 0.0,
    `review_count` INT NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX `idx_location` (`latitude`, `longitude`),
    INDEX `idx_category` (`category`),
    INDEX `idx_country` (`country_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create users table
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

-- Create token blacklist table
CREATE TABLE `token_blacklist` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `jti` VARCHAR(255) NOT NULL UNIQUE,
    `token_type` ENUM('access', 'refresh') NOT NULL DEFAULT 'access',
    `expires_at` TIMESTAMP NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    INDEX `idx_jti` (`jti`),
    INDEX `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create reviews table with user_id reference
CREATE TABLE `reviews` (
    `id` VARCHAR(36) PRIMARY KEY,
    `spot_id` VARCHAR(36) NOT NULL,
    `user_id` VARCHAR(36) NOT NULL,
    `reviewer_name` VARCHAR(100) NOT NULL,
    `rating` INT NOT NULL CHECK (`rating` >= 1 AND `rating` <= 5),
    `comment` TEXT,
    `rating_aspects` JSON,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT `fk_reviews_spot_id` FOREIGN KEY (`spot_id`) REFERENCES `spots`(`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_reviews_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_spot_user` (`spot_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;