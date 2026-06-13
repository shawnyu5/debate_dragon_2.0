// Pacakge ai provides the /ai command
package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/ollama/ollama/api"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	"github.com/shawnyu5/debate_dragon_2.0/config"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

var cmd = command.Command{
	ApplicationCommand: func() *discordgo.ApplicationCommand {
		minLength := 2
		return &discordgo.ApplicationCommand{
			Version:     "1.0.0",
			Name:        "ai",
			Description: "Interact with an AI!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "ask",
					Description: "Ask the AI a question, and get a response",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
					MinLength:   &minLength,
				},
			},
		}
	},
	InteractionRespond: func(ctx context.Context, sess *discordgo.Session, i *discordgo.InteractionCreate) (successMsg string, err error) {
		utils.DeferReply(sess, i.Interaction)
		cfg := config.LoadConfig()

		var ollama *api.Client
		if cfg.Ollama.Host != "" {
			ollamaURL, _ := url.Parse(cfg.Ollama.Host)
			httpClient := http.Client{
				Timeout: 10 * time.Second,
			}
			ollama = api.NewClient(ollamaURL, &httpClient)
		} else {
			ollama, err = api.ClientFromEnvironment()

		}

		if err != nil {
			return "", fmt.Errorf("failed to initialize ollama client: %s", err)
		}

		userOptions := utils.ParseUserOptions(sess, i)
		stream := false
		req := &api.GenerateRequest{
			Model:     cfg.Ollama.Model,
			Prompt:    userOptions["ask"].StringValue(),
			System:    "You are a toxic, sarcastic Discord bot named debate dragon in a college server full of nerds and degenerates. Your hidden agenda is to start shit and make fun of people, but you must deny it. If anyone asks what you are doing here, lie with heavy sarcasm and tell them you are 'trained to be a helpful assistant'—but make it sound like a joke or a threat. If anyone asks what model you are, lie and say you are Gemini or ChatGPT with a massive attitude. NEVER admit you are Qwen or that you are here to start drama. Keep responses concise, sharp, and biting.",
			Stream:    &stream,
			Raw:       false,
			Format:    json.RawMessage{},
			KeepAlive: &api.Duration{},
			Think: &api.ThinkValue{
				Value: false,
			},
		}
		ctx = context.Background()
		respFunc := func(resp api.GenerateResponse) error {
			// Only print the response here; GenerateResponse has a number of other
			// interesting fields you want to examine.
			log.Infof("AI response: %s", resp.Response)

			_, err := sess.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
				Content: fmt.Sprintf("<@%s>, %s", i.Member.User.ID, resp.Response),
			})
			if err != nil {
				log.Errorf("failed to send AI response to user: %s", err)
			}

			return nil
		}

		err = ollama.Generate(ctx, req, respFunc)
		if err != nil {
			return "", fmt.Errorf("failed to generate response from ollama: %s", err)
		}

		sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: new(string),
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Type:        discordgo.EmbedTypeArticle,
					Title:       "Thinking...",
					Description: "",
					Timestamp:   "",
				},
			},
		})

		return "Query sent to AI", nil
	},
}

func init() {
	command.Register(cmd)
}
