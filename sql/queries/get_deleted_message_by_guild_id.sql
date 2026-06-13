
-- name: GetDeletedMessagesByGuildID :many
-- GetDeletedMessagesByGuildID gets a specific guild's latest deleted message in a specific guild
--
-- There may not be any deleted messages in a guild, so this query may return 0 rows
SELECT *
FROM messages
WHERE guild_id = $1 AND deleted = true
ORDER BY created_at DESC
LIMIT 1;
