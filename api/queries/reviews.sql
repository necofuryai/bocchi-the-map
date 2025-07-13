-- name: CreateReview :exec
INSERT INTO reviews (
    id, spot_id, user_id, rating, comment
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: GetReviewByID :one
SELECT * FROM reviews 
WHERE id = ?;

-- name: GetReviewByUserAndSpot :one
SELECT * FROM reviews 
WHERE user_id = ? AND spot_id = ?;

-- name: UpdateReview :exec
UPDATE reviews 
SET rating = ?, comment = ?, updated_at = CURRENT_TIMESTAMP
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
  r.created_at,
  r.updated_at,
  u.name          AS user_name,
  u.picture       AS user_avatar
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
    COUNT(*) as review_count
FROM reviews 
WHERE spot_id = ?;

