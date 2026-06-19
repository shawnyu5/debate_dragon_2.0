-- name: GetChannelMessages :many
-- GetChannelMessages gets the past number of messages from a specific channel in a specific guild
SELECT *
FROM messages
WHERE channel_id = $1 AND guild_id = $2
ORDER BY created_at DESC
LIMIT $3;
