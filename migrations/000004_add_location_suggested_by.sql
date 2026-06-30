-- +goose Up
ALTER TABLE locations ADD COLUMN suggested_by UUID REFERENCES users(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE locations DROP COLUMN suggested_by;
