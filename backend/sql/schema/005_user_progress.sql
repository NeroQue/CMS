-- +goose Up
CREATE TABLE IF NOT EXISTS user_progress (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    content_item_id UUID NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    completed BOOLEAN NOT NULL DEFAULT false,
    progress_pct REAL NOT NULL DEFAULT 0,
    last_position INT,
    last_accessed TIMESTAMP,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE(user_id, content_item_id)
);

CREATE INDEX idx_user_progress_user_id ON user_progress(user_id);
CREATE INDEX idx_user_progress_content_item_id ON user_progress(content_item_id);

-- +goose Down
DROP INDEX IF EXISTS idx_user_progress_content_item_id;
DROP INDEX IF EXISTS idx_user_progress_user_id;

DROP TABLE IF EXISTS user_progress;