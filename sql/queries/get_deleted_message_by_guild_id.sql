
-- name: GetDeletedMessagesByGuildID :one
-- GetDeletedMessagesByGuildID gets a specific guild's latest deleted message in a specific guild
SELECT id, guild_id, author_id, metadata, created_at
FROM messages
WHERE guild_id = $1
ORDER BY created_at DESC
LIMIT 1;
