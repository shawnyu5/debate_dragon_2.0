-- +goose Up
CREATE TABLE messages (
   id UUID PRIMARY KEY NOT NULL,
   guild_id TEXT NOT NULL,
   author_id TEXT NOT NULL,
   metadata JSONB NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
