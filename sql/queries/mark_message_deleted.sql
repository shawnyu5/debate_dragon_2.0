-- name: MarkMessageDeleted :exec
-- MarkMessageDeleted marks a message in a guild as deleted
UPDATE messages
SET deleted = TRUE
WHERE message_id = $1 AND guild_id = $2;
