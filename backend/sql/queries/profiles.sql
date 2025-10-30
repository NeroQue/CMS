-- name: CreateProfile :one
INSERT INTO profiles (id, created_at, updated_at, name)
VALUES ($1,
        now(),
        now(),
        $2)
RETURNING *;

-- name: GetAllProfiles :many
SELECT *
FROM profiles;

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

-- name: UpdateProfileGamification :one
UPDATE profiles
SET experience       = $2,
    gems             = $3,
    level            = $4,
    last_active_date = $5,
    updated_at       = now()
WHERE id = $1
RETURNING *;

-- name: AddExperienceToProfile :one
UPDATE profiles
SET experience = experience + $2,
    level      = (experience + $2) / 100 + 1,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: GetProfileWithStats :one
SELECT p.*,
       (p.experience % 100)         as xp_in_level,
       (100 - (p.experience % 100)) as xp_to_next_level
FROM profiles p
WHERE p.id = $1;
