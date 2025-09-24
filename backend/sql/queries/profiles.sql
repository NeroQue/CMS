-- name: CreateProfile :one
INSERT INTO profiles (id, created_at, updated_at, name)
VALUES (
    $1,
    now(),
    now(),
    $2
)
RETURNING *;

-- name: GetAllProfiles :many
SELECT * FROM profiles;

-- name: GetProfileById :one
SELECT *
FROM profiles
WHERE id = $1;

-- name: GetProfileByName :one
SELECT *
FROM profiles
WHERE name = $1;

-- name: UpdateProfileByID :one
UPDATE profiles
SET name       = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteProfile :exec
DELETE
FROM profiles
WHERE id = $1;

-- name: GetProfilesByNamePattern :many
SELECT *
FROM profiles
WHERE name LIKE $1;

-- name: GetProfilesCount :one
SELECT COUNT(*)
FROM profiles;