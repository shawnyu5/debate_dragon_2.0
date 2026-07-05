-- name: GetUserMessages :many
-- GetUserMessages get messages sent by a user in a specific guild
SELECT *
FROM messages
WHERE author_id = $1 AND guild_id = $2
ORDER BY created_at DESC
LIMIT 10;
