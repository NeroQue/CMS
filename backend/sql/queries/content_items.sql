-- name: GetContentItem :one
SELECT * FROM content_items
WHERE id = $1;

-- name: ListContentItemsByModule :many
SELECT * FROM content_items
WHERE module_id = $1
ORDER BY "order" ASC;

-- name: CreateContentItem :one
INSERT INTO content_items (
    id,
    module_id,
    title,
    description,
    relative_path,
    content_type,
    duration,
    size,
    "order"
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateContentItem :one
UPDATE content_items
SET
    title = $2,
    description = $3,
    content_type = $4,
    duration = $5,
    "order" = $6,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteContentItem :exec
DELETE FROM content_items
WHERE id = $1;

