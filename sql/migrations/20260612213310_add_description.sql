-- +goose Up
COMMENT ON TABLE messages IS 'Table containing all messages send in all guilds this bot is in';

COMMENT ON COLUMN messages.guild_id IS 'Guild ID this message came from';
COMMENT ON COLUMN messages.author_id IS 'Author ID this message came from';
COMMENT ON COLUMN messages.metadata IS 'JSON encoded string, containg other properties, such as image and attachments';
COMMENT ON COLUMN messages.created_at IS 'Timestamp the message was created at';
