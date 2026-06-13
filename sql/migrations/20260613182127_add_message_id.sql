-- +goose Up
ALTER TABLE messages
ADD COLUMN message_id text NOT NULL;

-- +goose Down
ALTER TABLE messages
DROP COLUMN message_id;
