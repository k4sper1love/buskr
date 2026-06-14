-- +goose Up
ALTER TABLE invites 
ADD COLUMN created_by UUID REFERENCES users(id) ON DELETE SET NULL,
ADD COLUMN used_by UUID REFERENCES users(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE invites
DROP COLUMN created_by,
DROP COLUMN used_by;
