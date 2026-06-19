package ai

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/ollama/ollama/api"
	"github.com/shawnyu5/debate_dragon_2.0/config"
)

// RespondToPing when a user pings the bot, the pinged message is sent to AI for a response
func RespondToPing(ctx context.Context, sess *discordgo.Session, msg *discordgo.MessageCreate) {
	ollama, err := initOllamaClient()
	if err != nil {
		log.Errorf("failed to init ollama client: %s", err)
		return
	}

	cfg := config.LoadConfig()
	stream := false
	req := &api.GenerateRequest{
		Model: cfg.Ollama.Model,
		// TODO: provide conversation context
		Prompt:  fmt.Sprintf("Respond to this user's message: %s", msg.Content),
		System:  personalitySystemPrompt,
		Context: []int{},
		Stream:  &stream,
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
