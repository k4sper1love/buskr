-- +goose Up
ALTER TABLE locations ADD COLUMN is_veteran BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE locations DROP COLUMN IF EXISTS is_veteran;
