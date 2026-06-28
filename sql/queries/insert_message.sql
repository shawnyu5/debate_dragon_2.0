-- name: InsertMessage :exec
INSERT INTO messages (id, guild_id, channel_id, author_id, message_id, metadata)
VALUES ($1, $2, $3, $4, $5, $6);
