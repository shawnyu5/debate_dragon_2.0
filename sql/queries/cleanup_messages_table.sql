-- name: CleanupMessagesTable :exec
DELETE FROM messages WHERE id NOT IN (SELECT id FROM messages ORDER BY created_at DESC LIMIT 10000);
