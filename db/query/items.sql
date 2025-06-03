-- Items Queries
-- name: CreateItem :one
INSERT INTO items (
    name,
    description
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetItemByID :one
SELECT * FROM items
WHERE id = $1 LIMIT 1;

-- name: UpdateItem :one
UPDATE items
SET
    name = $2,
    description = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = $1;

-- name: ListItems :many
SELECT * FROM items
ORDER BY id
LIMIT $2 OFFSET $1;

-- name: CountItems :one
SELECT COUNT(*) FROM items;

