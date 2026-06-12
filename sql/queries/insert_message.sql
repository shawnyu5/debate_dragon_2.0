-- name: InsertMessage :exec
INSERT INTO messages (id, guild_id, author_id, metadata)
VALUES ($1, $2, $3, $4);
