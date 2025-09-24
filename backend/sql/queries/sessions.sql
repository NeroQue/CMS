-- name: CreateSession :one
INSERT INTO sessions (id, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    now(),
    now()
)
RETURNING *;

-- name: GetActiveSession :one
SELECT * FROM sessions
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: DeleteAllSessions :exec
DELETE FROM sessions;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = $1;
