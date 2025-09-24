-- name: GetModule :one
SELECT * FROM modules
WHERE id = $1;

-- name: ListModulesByCourse :many
SELECT * FROM modules
WHERE course_id = $1
ORDER BY "order" ASC;

-- name: CreateModule :one
INSERT INTO modules (
    id,
    course_id,
    title,
    description,
    relative_path,
    "order"
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateModule :one
UPDATE modules
SET
    title = $2,
    description = $3,
    "order" = $4,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteModule :exec
DELETE FROM modules
WHERE id = $1;

