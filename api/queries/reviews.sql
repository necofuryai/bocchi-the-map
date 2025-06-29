-- name: CreateReview :exec
INSERT INTO reviews (
    id, spot_id, user_id, rating, comment, rating_aspects
) VALUES (
    ?, ?, ?, ?, ?, ?
);

-- name: GetReviewByID :one
SELECT * FROM reviews 
WHERE id = ?;

-- name: GetReviewByUserAndSpot :one
SELECT * FROM reviews 
WHERE user_id = ? AND spot_id = ?;

-- name: UpdateReview :exec
UPDATE reviews 
SET rating = ?, comment = ?, rating_aspects = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteReview :exec
DELETE FROM reviews 
WHERE id = ?;

-- name: ListReviewsBySpot :many
SELECT
  r.id,
  r.spot_id,
  r.user_id,
  r.rating,
  r.comment,
  r.rating_aspects,
  r.created_at,
  r.updated_at,
  u.display_name  AS user_name,
  u.avatar_url    AS user_avatar
FROM reviews r
JOIN users u ON r.user_id = u.id
WHERE r.spot_id = ?
ORDER BY r.created_at DESC
LIMIT ? OFFSET ?;

-- name: CountReviewsBySpot :one
SELECT COUNT(*) FROM reviews 
WHERE spot_id = ?;

-- name: ListReviewsByUser :many
SELECT r.*, s.name as spot_name, s.category as spot_category
FROM reviews r
JOIN spots s ON r.spot_id = s.id
WHERE r.user_id = ?
ORDER BY r.created_at DESC
LIMIT ? OFFSET ?;

-- name: CountReviewsByUser :one
SELECT COUNT(*) FROM reviews 
WHERE user_id = ?;

-- name: GetSpotRatingStats :one
SELECT 
    AVG(rating) as average_rating,
    COUNT(*) as review_count,
    SUM(CASE WHEN rating = 5 THEN 1 ELSE 0 END) as five_star_count,
    SUM(CASE WHEN rating = 4 THEN 1 ELSE 0 END) as four_star_count,
    SUM(CASE WHEN rating = 3 THEN 1 ELSE 0 END) as three_star_count,
    SUM(CASE WHEN rating = 2 THEN 1 ELSE 0 END) as two_star_count,
    SUM(CASE WHEN rating = 1 THEN 1 ELSE 0 END) as one_star_count
FROM reviews 
WHERE spot_id = ?;

-- name: ListTopRatedSpots :many
SELECT
  s.id,
  s.name,
  s.category,
  s.address,
  s.latitude,
  s.longitude,
  s.country_code,
  s.average_rating,
  s.review_count,
  s.created_at,
  s.updated_at,
  AVG(r.rating)  AS avg_rating,
  COUNT(r.id)    AS total_reviews
FROM spots s
LEFT JOIN reviews r ON s.id = r.spot_id
GROUP BY s.id, s.name, s.category, s.address, s.latitude, s.longitude, s.country_code, s.average_rating, s.review_count, s.created_at, s.updated_at
HAVING COUNT(r.id) >= ?
ORDER BY avg_rating DESC, total_reviews DESC
LIMIT ? OFFSET ?;

-- name: CountTopRatedSpots :one
SELECT COUNT(*) FROM (
  SELECT s.id FROM spots s
  LEFT JOIN reviews r ON s.id = r.spot_id
  GROUP BY s.id
  HAVING COUNT(r.id) >= ?
) AS filtered_spots;