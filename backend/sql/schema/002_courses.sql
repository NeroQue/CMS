-- +goose Up
CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    creator_id UUID REFERENCES profiles(id),
    relative_path TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_courses_creator_id ON courses(creator_id);

-- +goose Down
DROP INDEX IF EXISTS idx_courses_creator_id;
DROP TABLE IF EXISTS courses;