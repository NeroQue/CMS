-- +goose Up
CREATE TABLE IF NOT EXISTS modules (
    id UUID PRIMARY KEY,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    relative_path TEXT NOT NULL,
    "order" INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_modules_course_id ON modules(course_id);

-- +goose Down
DROP INDEX IF EXISTS idx_modules_course_id;
DROP TABLE IF EXISTS modules;