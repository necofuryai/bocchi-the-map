-- Token blacklist table for secure logout and token revocation
CREATE TABLE IF NOT EXISTS token_blacklist (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    jti VARCHAR(255) NOT NULL UNIQUE, -- JWT ID for identifying tokens
    user_id VARCHAR(255) NOT NULL,    -- User who owns the token
    token_type ENUM('access', 'refresh') NOT NULL DEFAULT 'access',
    expires_at TIMESTAMP NOT NULL,    -- When the token expires
    revoked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When it was revoked
    reason VARCHAR(255),              -- Reason for revocation (logout, security, etc.)
    
    INDEX idx_user_id (user_id),
    INDEX idx_revoked_at (revoked_at),
    INDEX idx_token_blacklist_jti_expires (jti, expires_at)
);

