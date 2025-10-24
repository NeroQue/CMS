-- name: FactoryResetDatabase :exec
-- Reset the database to its initial state, RESTART IDENTITY CASCADE, means that the primary keys will be reset and any
-- incrementing columns will be reset to 1.
TRUNCATE user_progress, sessions, content_items, modules, courses, profiles RESTART IDENTITY CASCADE;
