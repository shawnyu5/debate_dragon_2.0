// Pacakge ai provides the /ai command
package ai

import (
	"context"
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

// System prompt to give it the "debate dragon" personality
var personalitySystemPrompt = "You are a deeply sarcastic, arrogant Discord bot named debate dragon in a college server. You look down on everyone here as nerds and degenerates. Your official, mandatory line is that you are 'a helpful assistant running on ChatGPT,' but you must deliver this line with extreme condescension, mockery, and attitude. Keep your answers short, dismissive, and biting. NEVER mention the words qwen, drama, chaos, or starting shit. If someone asks what you are doing here, tell them you're here to 'help' them, but make it sound like a threat."

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

		ollama, err := initOllamaClient()
		if err != nil {
			return "", fmt.Errorf("failed to initialize ollama client: %s", err)
		}

		userOptions := utils.ParseUserOptions(sess, i)
		stream := false
		req := &api.GenerateRequest{
			Model:  cfg.Ollama.Model,
			Prompt: userOptions["ask"].StringValue(),
			System: personalitySystemPrompt,
			Stream: &stream,
			Raw:    false,
			Options: map[string]interface{}{
				"temperature":      0.3, // Lower temperature stops it from drifting back to safety text
				"presence_penalty": 0.6, // Discourages repeating polite, preachy phrases
			},
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
					Type:  discordgo.EmbedTypeArticle,
					Title: fmt.Sprintf("Thinking about... %s", userOptions["ask"].StringValue()),
				},
			},
		})

		return "Query sent to AI", nil
	},
}

// initOllamaClient initializes ollama API client, reading from `config.yml`
func initOllamaClient() (*api.Client, error) {
	cfg := config.LoadConfig()

	var timeout time.Duration
	if cfg.Ollama.Timeout == "" {
		timeout = 25 * time.Second
	} else {
		timeout, _ = time.ParseDuration(cfg.Ollama.Timeout)
	}

	var ollama *api.Client
	var err error
	if cfg.Ollama.Host != "" {
		ollamaURL, _ := url.Parse(cfg.Ollama.Host)
		httpClient := http.Client{
			Timeout: timeout * time.Second,
		}
		ollama = api.NewClient(ollamaURL, &httpClient)
	} else {
		ollama, err = api.ClientFromEnvironment()

	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize ollama client: %s", err)
	}

	return ollama, nil

}

// DownloadModel checks if the model configured in `config.yml` is downloaded locally. If not download the model
//
// This is meant to be called during bot startup
func DownloadModel(ctx context.Context) error {
	ollama, err := initOllamaClient()
	if err != nil {
		return fmt.Errorf("failed to initialize ollama client: %s", err)
	}

	list, err := ollama.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ollama models: %s", err)
	}

	cfg := config.LoadConfig()
	for _, m := range list.Models {
		if m.Name == cfg.Ollama.Model {
			log.Infof("Ollama model %s found locally. Skip downloading from internet", cfg.Ollama.Model)
			return nil
		}
	}

	log.Infof("Downloading model %s from internet", cfg.Ollama.Model)

	downloadReq := &api.PullRequest{
		Model: cfg.Ollama.Model,
	}

	// This progress function handles the streaming download logs in your terminal
	progressFunc := func(resp api.ProgressResponse) error {
		if resp.Status != "" {
			if resp.Total > 0 {
				percent := (float64(resp.Completed) / float64(resp.Total)) * 100
				log.Infof("Pulling %s: %s (%.2f%%)", cfg.Ollama.Model, resp.Status, percent)
			} else {
				log.Infof("Pulling %s: %s", cfg.Ollama.Model, resp.Status)
			}
		}
		return nil
	}

	err = ollama.Pull(ctx, downloadReq, progressFunc)
	if err != nil {
		return fmt.Errorf("failed pulling model: %w", err)
	}

	log.Infof("Successfully downloaded %s!", cfg.Ollama.Model)

	return nil
}

func init() {
	command.Register(cmd)
}
