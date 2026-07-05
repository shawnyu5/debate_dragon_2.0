package models

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/db"
)

// SavedMessage represents the hydrated message structure saved in the Database
type SavedMessage struct {
	ID          string                         `json:"id"`
	Content     string                         `json:"content"`
	AuthorID    string                         `json:"author_id"`
	MessageID   string                         `json:"message_id"`
	Attachments []*discordgo.MessageAttachment `json:"attachments,omitempty"`
}

// DBMessageToSavedMessage converts a `db.Message` to a `SavedMessage`
func DBMessageToSavedMessage(dbMsg db.Message) (*SavedMessage, error) {
	var richDTO SavedMessage
	err := json.Unmarshal(dbMsg.Metadata, &richDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize message metadata: %w", err)
	}

	richDTO.AuthorID = dbMsg.AuthorID
	return &richDTO, nil

}
