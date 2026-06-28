package messagetracking

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shawnyu5/debate_dragon_2.0/config"
	"github.com/shawnyu5/debate_dragon_2.0/db"
)

// PrepareMessageForDB converts a discordgo.Message into a format that can be inserted into Postgres
func PrepareMessageForDB(msg *discordgo.Message) ([]byte, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		log.Fatalf("failed to generate new UUID: %s", err)
	}

	savedMsg := db.SavedMessage{
		ID:          uuid.String(),
		Content:     msg.Content,
		MessageID:   msg.ID,
		Attachments: msg.Attachments,
	}

	jsonData, err := json.Marshal(savedMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal discord message: %w", err)
	}

	return jsonData, nil
}

// DBMessageToRichMessage retrieves a rich message from the database and deserializes its metadata.
//
// Returns:
//   - *db.SavedMessage: A pointer to the deserialized db.SavedMessage if successful.
//   - error: An error if the deserialization fails
func DBMessageToRichMessage(dbMsg db.Message) (*db.SavedMessage, error) {
	var richDTO db.SavedMessage
	err := json.Unmarshal(dbMsg.Metadata, &richDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize message metadata: %w", err)
	}

	richDTO.AuthorID = dbMsg.AuthorID
	return &richDTO, nil
}

// TrackAllSentMessage tracks all sent messages in all guilds, up to 1000 messages total
func TrackAllSentMessage(store *db.Store, msg *discordgo.MessageCreate) {
	log.Infof("Storing message in DB: \"%s\" from user @%s", msg.Content, msg.Author.Username)
	log.Debugf("Got discord message: %+v", msg.Message)

	if msg.Content == "" {
		log.Warn("Message empty, not storing in DB")
		return
	}

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
		GuildID: msg.GuildID,
		ChannelID: pgtype.Text{
			String: msg.ChannelID,
			Valid:  true,
		},
		AuthorID:  msg.Author.ID,
		MessageID: msg.ID,
		Metadata:  json,
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

// GetDeletedMessageByAuthorID returns the last deleted messages in a specific `guildID` for a specific `authorID`
func GetDeletedMessageByAuthorID(store *db.Store, guildID, authorID string) (*db.SavedMessage, error) {
	message, err := store.GetDeletedMessagesByAuthorID(context.Background(), db.GetDeletedMessagesByAuthorIDParams{
		AuthorID: authorID,
		GuildID:  guildID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deleted messages by author ID: %s", err)
	}
	richMsg, err := DBMessageToRichMessage(message)
	if err != nil {
		return nil, fmt.Errorf("failed to convert DB message to rich message: %s", err)
	}

	return richMsg, nil
}

// GetDeletedMessageByGuildID get the last deleted message by guild ID
func GetDeletedMessageByGuildID(store *db.Store, guildID string) (*db.SavedMessage, error) {
	messages, err := store.GetDeletedMessagesByGuildID(context.Background(), guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last deleted message: %s", err)
	}

	if len(messages) == 0 {
		return &db.SavedMessage{}, nil
	}

	richMsg, err := DBMessageToRichMessage(messages[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert DB message to rich message: %s", err)
	}

	return richMsg, nil
}

// Deprecated: does nothing now
func GetLastDeletedMessage() discordgo.Message {
	return discordgo.Message{}
	// return lastDeletedMessage
}

// TrackDeletedMessage marks a message as been deleted in the DB. If the message does not exist, nothing is done, since Discord does not provide the content of the message on messageDelete event
func TrackDeletedMessage(ctx context.Context, msg *discordgo.MessageDelete) {
	log.Infof("Marking message %s as deleted in DB", msg.ID)

	cfg := config.LoadConfig()
	// Ignore this rule in Dev mode, otherwise we cant test this thing...
	if !cfg.DevMode && msg.Author.ID == cfg.BotOwner {
		log.Info("Message sent by bot owner. Not snipable. Ignoring...")
		return
	}

	store, err := db.StoreFromContext(ctx)
	if err != nil {
		log.Fatalf("No db found in context: %s", err)
	}

	err = store.ExecTx(context.Background(), func(q *db.Queries) error {
		err := q.MarkMessageDeleted(context.Background(), db.MarkMessageDeletedParams{
			MessageID: msg.ID,
			GuildID:   msg.GuildID,
		})
		return err
	})

	if err != nil {
		log.Errorf("Failed to mark message as deleted in DB: %s", err)
	}
}
