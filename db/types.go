package db

import "github.com/bwmarrin/discordgo"

// SavedMessage represents the exact rich structure we want to keep in the DB
type SavedMessage struct {
	ID          string                         `json:"id"`
	Content     string                         `json:"content"`
	AuthorID    string                         `json:"author_id"`
	MessageID   string                         `json:"message_id"`
	Attachments []*discordgo.MessageAttachment `json:"attachments,omitempty"`
}
