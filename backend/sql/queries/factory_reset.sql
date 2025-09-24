-- name: FactoryResetDatabase :exec
-- Clear all data from tables in the correct order (respecting foreign keys)
-- Start with dependent tables first, then parent tables

-- Clear user progress (depends on profiles and content_items)
DELETE FROM user_progress;

-- Clear sessions (depends on profiles)
DELETE FROM sessions;

-- Clear content items (depends on modules)
DELETE FROM content_items;

-- Clear modules (depends on courses)
DELETE FROM modules;

-- Clear courses (depends on profiles via creator_id)
DELETE FROM courses;

-- Clear profiles last (no dependencies)
DELETE FROM profiles;
