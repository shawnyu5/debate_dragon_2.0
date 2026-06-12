package messagetracking

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shawnyu5/debate_dragon_2.0/db"
)

// contains all discord messages sent since the bot startup, map of guild id to message id to discord message
var allMessagesMap = make(map[string]map[string]discordgo.Message)

// map of guild id to author ID to list of deleted messages
var deletedMessagesMap = make(map[string]map[string][]discordgo.Message)

var lastDeletedMessage = discordgo.Message{}

// PrepareMessageForDB converts a discordgo.Message into a format that can be inserted into Postgres
func PrepareMessageForDB(msg *discordgo.Message) ([]byte, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		log.Fatalf("failed to generate new UUID: %s", err)
	}

	savedMsg := db.SavedMessage{
		ID:          uuid.String(),
		Content:     msg.Content,
		AuthorID:    msg.Author.ID,
		Attachments: msg.Attachments,
	}

	jsonData, err := json.Marshal(savedMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal discord message: %w", err)
	}

	return jsonData, nil
}

// TrackAllSentMessage tracks all sent messages in all guilds, up to 1000 messages total
func TrackAllSentMessage(store *db.Store, msg *discordgo.MessageCreate) {
	log.Infof("Storing message in DB: %s", msg.Content)
	log.Debugf("Got discord message: %+v", msg.Message)
	uuid, err := uuid.NewV7()
	if err != nil {
		log.Error("failed to create UUID: %s", err)
	}

	json, err := PrepareMessageForDB(msg.Message)
	if err != nil {
		log.Fatalf("failed to convert Discord message to JSON: %s", err)
	}

	err = store.InsertMessage(context.Background(), db.InsertMessageParams{
		ID: pgtype.UUID{
			Bytes: uuid,
			Valid: true,
		},
		GuildID:  msg.GuildID,
		AuthorID: msg.Author.ID,
		Metadata: json,
	})
	if err != nil {
		log.Errorf("failed to insert new message into DB: %s", err)
	}

	go func() {
		log.Info("Cleaning up messages table")
		ctx := context.Background()
		err := store.CleanupMessagesTable(ctx)
		if err != nil {
			log.Errorf("failed to clean up messages table: %s", err)
		}

	}()
}

// TrackDeletedMessage tracks the last 10 deleted messages metadata for a guild, excluding their content. Use allMessagesMap to get the contents of those messages
//
// When there are 10 messages, the oldest message will be deleted.
//
// Deprecated: no longer need to specifically track deleted messages anymore
func TrackDeletedMessage(guildID string, messageID string) {
	if deletedMessagesMap[guildID] == nil {
		log.Debugf("Creating deleted messages guild map for guild %s", guildID)
		deletedMessagesMap[guildID] = make(map[string][]discordgo.Message)

	}

	// We are only able to find deleted message that was deleted when the bot is alive.
	mess := GetMessageByID(guildID, messageID)
	// No contents means the message was deleted before the bot was alive, nothing we can do about it...
	if mess.Content == "" {
		return
	}

	if len(deletedMessagesMap[mess.GuildID][mess.Author.ID]) == 10 {
		deletedMessagesMap[mess.GuildID][mess.Author.ID] = deletedMessagesMap[mess.GuildID][mess.Author.ID][1:]
	}
	deletedMessagesMap[mess.GuildID][mess.Author.ID] = append(deletedMessagesMap[mess.GuildID][mess.Author.ID], mess)
	lastDeletedMessage = mess
}

// GetMessageByID returns the message with the given `messageID` from the given `guildID`.
func GetMessageByID(guildID, messageID string) discordgo.Message {
	return allMessagesMap[guildID][messageID]
}

// GetDeletedMessagesByAuthorID returns a list of deleted messages in a specific `guildID` for a specific `authorID`
func GetDeletedMessagesByAuthorID(guildID, authorID string) []discordgo.Message {
	usrMsg := deletedMessagesMap[guildID][authorID]
	return usrMsg
}

func GetLastDeletedMessage() discordgo.Message {
	return lastDeletedMessage
}
