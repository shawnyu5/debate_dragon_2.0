-- +goose Up
ALTER TABLE messages
ADD COLUMN channel_id TEXT;

-- +goose Down
ALTER TABLE messages
DROP COLUMN channel_id;
