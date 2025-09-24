-- name: GetCourse :one
SELECT * FROM courses
WHERE id = $1;

-- name: ListCourses :many
SELECT * FROM courses
ORDER BY created_at DESC;

-- name: ListCoursesByCreator :many
SELECT * FROM courses
WHERE creator_id = $1
ORDER BY created_at DESC;

-- name: CreateCourse :one
INSERT INTO courses (
    id,
    title,
    description,
    creator_id,
    relative_path
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateCourse :one
UPDATE courses
SET
    title = $2,
    description = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteCourse :exec
DELETE FROM courses
WHERE id = $1;

