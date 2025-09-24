-- +goose Up
CREATE TABLE IF NOT EXISTS content_items (
    id UUID PRIMARY KEY,
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    relative_path TEXT NOT NULL,
    content_type TEXT NOT NULL,
    duration INT,
    size BIGINT,
    "order" INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_content_items_module_id ON content_items(module_id);

-- +goose Down
DROP INDEX IF EXISTS idx_content_items_module_id;
DROP TABLE IF EXISTS content_items;