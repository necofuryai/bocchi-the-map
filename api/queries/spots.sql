-- name: CreateSpot :exec
INSERT INTO spots (
    id, name, category, address, country_code
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: GetSpotByID :one
SELECT id, name, category, address, country_code, created_at, updated_at FROM spots 
WHERE id = ?;

-- name: UpdateSpot :exec
UPDATE spots 
SET name = ?, category = ?, 
    address = ?, country_code = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;


-- name: ListSpotsByCategory :many
SELECT id, name, category, address, country_code, created_at, updated_at FROM spots 
WHERE category = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountSpotsByCategory :one
SELECT COUNT(*) FROM spots 
WHERE category = ?;

-- name: ListSpotsByCountry :many
SELECT id, name, category, address, country_code, created_at, updated_at FROM spots 
WHERE country_code = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountSpotsByCountry :one
SELECT COUNT(*) FROM spots 
WHERE country_code = ?;

-- name: ListSpots :many
SELECT id, name, category, address, country_code, created_at, updated_at FROM spots 
WHERE (? = '' OR name LIKE ?)
  AND (? = '' OR category = ?)
  AND (? = '' OR country_code = ?)
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountSpots :one
SELECT COUNT(*) FROM spots 
WHERE (? = '' OR name LIKE ?)
  AND (? = '' OR category = ?)
  AND (? = '' OR country_code = ?);

-- name: DeleteSpot :exec
DELETE FROM spots 
WHERE id = ?;