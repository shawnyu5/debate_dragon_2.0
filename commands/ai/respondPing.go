package ai

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ollama/ollama/api"
	messagetracking "github.com/shawnyu5/debate_dragon_2.0/commands/messageTracking"
	"github.com/shawnyu5/debate_dragon_2.0/config"
	"github.com/shawnyu5/debate_dragon_2.0/db"
)

// RespondToPing when a user pings the bot, the pinged message is sent to AI for a response
func RespondToPing(ctx context.Context, sess *discordgo.Session, msg *discordgo.MessageCreate) {
	ollama, err := initOllamaClient()
	if err != nil {
		log.Errorf("failed to init ollama client: %s", err)
		return
	}

	dbstore, err := db.StoreFromContext(ctx)
	if err != nil {
		log.Fatalf("failed to get DB from context: %s", err)
	}

	discordMsgs, err := dbstore.GetChannelMessages(ctx, db.GetChannelMessagesParams{
		ChannelID: pgtype.Text{
			String: msg.ChannelID,
			Valid:  true,
		},
		GuildID: msg.GuildID,
		Limit:   10,
	})
	if err != nil {
		log.Errorf("failed to get discord messages from channel: %s", err)
		return
	}

	for _, msg := range discordMsgs {
		richMsg, _ := messagetracking.DBMessageToRichMessage(msg)
		log.Infof("Message: %s", richMsg.Content)
	}

	cfg := config.LoadConfig()
	stream := false
	// api.ChatRequest{
	// 	Model:           "",
	// 	Messages:        []api.Message{},
	// 	Stream:          new(bool),
	// 	Format:          json.RawMessage{},
	// 	KeepAlive:       &api.Duration{},
	// 	Tools:           api.Tools{},
	// 	Options:         map[string]any{},
	// 	Think:           &api.ThinkValue{},
	// 	Truncate:        new(bool),
	// 	Shift:           new(bool),
	// 	DebugRenderOnly: false,
	// 	Logprobs:        false,
	// 	TopLogprobs:     0,
	// }
	req := &api.GenerateRequest{
		Model: cfg.Ollama.Model,
		// TODO: provide conversation context
		Prompt:  fmt.Sprintf("Respond to this user's message: %s", msg.Content),
		System:  personalitySystemPrompt,
		Context: []int{},
		Stream:  &stream,
		Options: map[string]interface{}{
			"temperature":      0.3, // Lower temperature stops it from drifting back to safety text
			"presence_penalty": 0.6, // Discourages repeating polite, preachy phrases
		},
		Think: &api.ThinkValue{
			Value: false,
		},
	}

	err = ollama.Generate(ctx, req, func(gr api.GenerateResponse) error {
		_, err := sess.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("<@%s>, %s", msg.Author.ID, gr.Response),
		})

		return err
	})
	if err != nil {
		log.Errorf("failed to get ollama response: %s", err)
		log.Error("Not responding to ping")
	}

}
