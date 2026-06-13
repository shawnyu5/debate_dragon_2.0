-- +goose Up
ALTER TABLE messages ADD COLUMN deleted BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE messages DROP COLUMN deleted;
