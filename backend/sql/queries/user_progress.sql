-- name: GetUserProgressByContentItem :one
SELECT * FROM user_progress
WHERE user_id = $1 AND content_item_id = $2;

-- name: UpsertUserProgress :one
INSERT INTO user_progress (
    id, user_id, content_item_id, completed, progress_pct, last_position, last_accessed, created_at, updated_at
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4, $5, $6, now(), now()
)
ON CONFLICT (user_id, content_item_id)
DO UPDATE SET
    completed = EXCLUDED.completed,
    progress_pct = EXCLUDED.progress_pct,
    last_position = EXCLUDED.last_position,
    last_accessed = EXCLUDED.last_accessed,
    updated_at = now()
RETURNING *;

-- name: ListUserProgressByCourse :many
SELECT up.* FROM user_progress up
JOIN content_items ci ON up.content_item_id = ci.id
JOIN modules m ON ci.module_id = m.id
WHERE m.course_id = $1 AND up.user_id = $2
ORDER BY m."order", ci."order";

-- name: GetModuleProgressStats :one
SELECT
    COUNT(*) as total_items,
    COUNT(*) FILTER (WHERE up.completed = true) as completed_items,
    COALESCE(AVG(up.progress_pct), 0) as avg_progress
FROM content_items ci
LEFT JOIN user_progress up ON ci.id = up.content_item_id AND up.user_id = $2
WHERE ci.module_id = $1;

-- name: GetCourseProgressStats :one
SELECT
    COUNT(DISTINCT m.id) as total_modules,
    COUNT(DISTINCT ci.id) as total_items,
    COUNT(DISTINCT ci.id) FILTER (WHERE up.completed = true) as completed_items,
    MAX(up.last_accessed) as last_accessed
FROM modules m
LEFT JOIN content_items ci ON m.id = ci.module_id
LEFT JOIN user_progress up ON ci.id = up.content_item_id AND up.user_id = $2
WHERE m.course_id = $1;
