-- +goose Up
-- Add gamification columns to profiles table
ALTER TABLE profiles
    ADD COLUMN experience       INT NOT NULL DEFAULT 0,
    ADD COLUMN gems             INT NOT NULL DEFAULT 0,
    ADD COLUMN level            INT NOT NULL DEFAULT 1,
    ADD COLUMN streak           INT NOT NULL DEFAULT 0,
    ADD COLUMN last_active_date TIMESTAMP;

-- Add XP value column to content_items for manual overrides
ALTER TABLE content_items
    ADD COLUMN xp_value INT;

-- Add tracking to prevent double XP awards after resets
ALTER TABLE user_progress
    ADD COLUMN xp_awarded BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN xp_amount  INT     NOT NULL DEFAULT 0;

CREATE INDEX idx_profiles_level ON profiles (level);
CREATE INDEX idx_profiles_experience ON profiles (experience);

-- +goose Down
DROP INDEX IF EXISTS idx_profiles_experience;
DROP INDEX IF EXISTS idx_profiles_level;

ALTER TABLE user_progress
    DROP COLUMN xp_amount,
    DROP COLUMN xp_awarded;

ALTER TABLE content_items
    DROP COLUMN xp_value;

ALTER TABLE profiles
    DROP COLUMN last_active_date,
    DROP COLUMN streak,
    DROP COLUMN level,
    DROP COLUMN gems,
    DROP COLUMN experience;

