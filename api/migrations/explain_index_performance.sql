-- EXPLAIN Analysis for Composite Index Performance
-- Use this file to verify that the composite indexes improve query performance

-- Query Pattern 1: Get reviews for a spot ordered by creation date (most common)
-- This benefits from idx_reviews_spot_created (spot_id, created_at DESC)
EXPLAIN SELECT * 
FROM reviews 
WHERE spot_id = 'spot123' 
ORDER BY created_at DESC 
LIMIT 10;

-- Query Pattern 2: Get user's review history ordered by date
-- This benefits from idx_reviews_user_created (user_id, created_at DESC)
EXPLAIN SELECT r.*, s.name as spot_name 
FROM reviews r 
JOIN spots s ON r.spot_id = s.id 
WHERE r.user_id = 'user456' 
ORDER BY r.created_at DESC 
LIMIT 20;

-- Query Pattern 3: Get rating statistics for a spot
-- This benefits from idx_reviews_spot_rating (spot_id, rating)
EXPLAIN SELECT 
    AVG(rating) as average_rating,
    COUNT(*) as review_count,
    COUNT(CASE WHEN rating = 5 THEN 1 END) as five_star_count,
    COUNT(CASE WHEN rating = 4 THEN 1 END) as four_star_count,
    COUNT(CASE WHEN rating = 3 THEN 1 END) as three_star_count,
    COUNT(CASE WHEN rating = 2 THEN 1 END) as two_star_count,
    COUNT(CASE WHEN rating = 1 THEN 1 END) as one_star_count
FROM reviews 
WHERE spot_id = 'spot123';

-- Query Pattern 4: Find spots with high ratings (4+ stars)
-- This benefits from idx_reviews_rating_spot (rating, spot_id)
EXPLAIN SELECT s.*, AVG(r.rating) AS avg_rating, COUNT(r.id) AS total_reviews
FROM spots s
JOIN reviews r ON s.id = r.spot_id
WHERE r.rating >= 4
GROUP BY s.id
HAVING COUNT(r.id) >= 5
ORDER BY avg_rating DESC, total_reviews DESC
LIMIT 20;

-- Query Pattern 5: Check if user already reviewed a spot
-- This benefits from the existing unique constraint (user_id, spot_id)
EXPLAIN SELECT id FROM reviews 
WHERE user_id = 'user456' AND spot_id = 'spot123';

-- Query Pattern 5b: Alternative EXISTS pattern for checking user review existence
-- Compare execution plans with the direct SELECT approach above
EXPLAIN SELECT EXISTS (
    SELECT 1 FROM reviews 
    WHERE user_id = 'user456' AND spot_id = 'spot123'
) AS has_reviewed;

-- Expected performance improvements:
-- 1. Queries with WHERE spot_id + ORDER BY created_at can use covering index
-- 2. Queries with WHERE user_id + ORDER BY created_at can use covering index  
-- 3. Rating aggregation queries can efficiently scan by spot_id + rating
-- 4. High-rating filtering can efficiently scan by rating + spot_id
-- 5. Unique constraint already optimizes user+spot lookups

-- Note: Run these EXPLAIN queries both before and after applying the migration
-- to compare query execution plans and verify index usage.