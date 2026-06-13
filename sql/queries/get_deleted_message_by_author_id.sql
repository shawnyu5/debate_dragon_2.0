
-- name: GetDeletedMessagesByAuthorID :one
-- GetDeletedMessagesByAuthorID gets a specific author's latest deleted message in a specific guild
SELECT *
FROM messages
WHERE author_id = $1 AND guild_id = $2
ORDER BY created_at DESC
LIMIT 1;
