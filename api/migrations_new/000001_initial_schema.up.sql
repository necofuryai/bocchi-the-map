-- Initial schema for bocchi-the-map
CREATE TABLE
  `users` (
    `id` VARCHAR(36) PRIMARY KEY,
    `anonymous_id` VARCHAR(36),
    `email` VARCHAR(255) UNIQUE NOT NULL,
    `display_name` VARCHAR(100) NOT NULL,
    `avatar_url` TEXT,
    `auth_provider` ENUM('google', 'x') NOT NULL,
    `auth_provider_id` VARCHAR(255) NOT NULL,
    `preferences` JSON,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `unique_provider_user` UNIQUE KEY (`auth_provider`, `auth_provider_id`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE spots (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    name_i18n JSON,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    category VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    address_i18n JSON,
    country_code CHAR(2) NOT NULL,
    average_rating DECIMAL(3, 1) NOT NULL DEFAULT 0.0,
    review_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_location (latitude, longitude),
    INDEX idx_category (category),
    INDEX idx_country (country_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE reviews (
    id VARCHAR(36) PRIMARY KEY,
    spot_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    rating_aspects JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_reviews_spot_id FOREIGN KEY (spot_id) REFERENCES spots(id) ON DELETE CASCADE,
    CONSTRAINT fk_reviews_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_spot_review UNIQUE KEY (user_id, spot_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;