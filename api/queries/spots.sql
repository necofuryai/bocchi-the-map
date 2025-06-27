-- name: CreateSpot :exec
INSERT INTO spots (
    id, name, name_i18n, latitude, longitude, category, address, address_i18n, country_code
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetSpotByID :one
SELECT * FROM spots 
WHERE id = ?;

-- name: UpdateSpot :exec
UPDATE spots 
SET name = ?, name_i18n = ?, latitude = ?, longitude = ?, category = ?, 
    address = ?, address_i18n = ?, country_code = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateSpotRating :exec
UPDATE spots 
SET average_rating = ?, review_count = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: ListSpotsByLocation :many
SELECT * FROM spots 
WHERE (6371 * acos(
    cos(radians(?)) * cos(radians(latitude)) * 
    cos(radians(longitude) - radians(?)) + 
    sin(radians(?)) * sin(radians(latitude))
)) <= ?
ORDER BY (6371 * acos(
    cos(radians(?)) * cos(radians(latitude)) * 
    cos(radians(longitude) - radians(?)) + 
    sin(radians(?)) * sin(radians(latitude))
))
LIMIT ? OFFSET ?;

-- name: CountSpotsByLocation :one
SELECT COUNT(*) FROM spots 
WHERE (6371 * acos(
    cos(radians(?)) * cos(radians(latitude)) * 
    cos(radians(longitude) - radians(?)) + 
    sin(radians(?)) * sin(radians(latitude))
)) <= ?;

-- name: ListSpotsByCategory :many
SELECT * FROM spots 
WHERE category = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountSpotsByCategory :one
SELECT COUNT(*) FROM spots 
WHERE category = ?;

-- name: ListSpotsByCountry :many
SELECT * FROM spots 
WHERE country_code = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountSpotsByCountry :one
SELECT COUNT(*) FROM spots 
WHERE country_code = ?;

-- name: SearchSpots :many
SELECT * FROM spots 
WHERE (name LIKE ? OR address LIKE ?)
  AND (? = '' OR category = ?)
  AND (? = '' OR country_code = ?)
  AND (? = 0 OR (6371 * acos(
      cos(radians(?)) * cos(radians(latitude)) * 
      cos(radians(longitude) - radians(?)) + 
      sin(radians(?)) * sin(radians(latitude))
  )) <= ?)
ORDER BY 
  CASE WHEN name LIKE ? THEN 1 ELSE 2 END,
  average_rating DESC,
  review_count DESC
LIMIT ? OFFSET ?;

-- name: CountSearchSpots :one
SELECT COUNT(*) FROM spots 
WHERE (name LIKE ? OR address LIKE ?)
  AND (? = '' OR category = ?)
  AND (? = '' OR country_code = ?)
  AND (? = 0 OR (6371 * acos(
      cos(radians(?)) * cos(radians(latitude)) * 
      cos(radians(longitude) - radians(?)) + 
      sin(radians(?)) * sin(radians(latitude))
  )) <= ?);

-- name: DeleteSpot :exec
DELETE FROM spots 
WHERE id = ?;