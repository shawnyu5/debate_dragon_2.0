package emotes

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

type Emotes struct{}

type emotesMapType map[string]string

// map of emote name to emote url
var emoteCache = make(emotesMapType)

// Components implements commands.Command
func (Emotes) Components() []commands.Component {
	return nil
}

// Def implements commands.Command
func (Emotes) Def() *discordgo.ApplicationCommand {
	def := &discordgo.ApplicationCommand{
		Version:     "1.0.0",
		Name:        "emote",
		Description: "Send custom emotes",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "name",
				Description:  "name of emote",
				Required:     true,
				Autocomplete: true,
			},
		},
	}

	// emotes := GetEmotes()
	// for name, url := range emotes {
	// def.Options[0].Choices = append(def.Options[0].Choices, &discordgo.ApplicationCommandOptionChoice{
	// Name:  name,
	// Value: url,
	// })
	// }

	return def
}

// Handler implements commands.Command
func (Emotes) Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	input := utils.ParseUserOptions(sess, i)
	sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Preview window",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	// err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	// Data: &discordgo.InteractionResponseData{
	// Flags: discordgo.MessageFlagsEphemeral,
	// },
	// })

	switch i.Type {
	case discordgo.InteractionApplicationCommandAutocomplete:
		fmt.Println("autocomplete Handler") // __AUTO_GENERATED_PRINTF__
		userInput := input["name"]

		emotes := GetEmotes()
		choices := []*discordgo.ApplicationCommandOptionChoice{}

		// collect names of all emotes
		emoteNames := make([]string, len(emotes))
		idx := 0
		for name := range emotes {
			emoteNames[idx] = name
			idx++
		}

		if userInput.StringValue() != "" {

			matches := fuzzy.Find(userInput.StringValue(), emoteNames)

			for _, match := range matches {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  match,
					Value: emotes[match],
				})

			}
		} else {
			for _, name := range emoteNames {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  name,
					Value: emotes[name],
				})
			}
		}

		// if _, found := emotes[userInput.StringValue()]; found {
		// fmt.Println("Handler") // __AUTO_GENERATED_PRINTF__
		// content := emotes[userInput.StringValue()]
		// _, err := sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		// Content: &content,
		// })
		// fmt.Printf("Handler err: %v\n", err) // __AUTO_GENERATED_PRINT_VAR__
		// }

		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})

	default:
		fmt.Println(input["name"].Focused)
		emoteUrl := input["name"].StringValue()
		_, err := sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &emoteUrl,
		})
		// err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Type: discordgo.InteractionResponseChannelMessageWithSource,
		// Data: &discordgo.InteractionResponseData{
		// Content: input["name"].StringValue(),
		// },
		// })

		if err != nil {
			return "", err
		}
		return "emote sent", nil
	}
	return "emote sent", nil
}

// GetEmotes get emotes from config.json.
// return: map of emote name to emote url.
func GetEmotes() emotesMapType {
	if len(emoteCache) != 0 {
		return emoteCache
	}
	c := utils.LoadConfig()
	for _, emote := range c.Emotes {
		emoteCache[emote.Name] = emote.URL
	}
	return emoteCache
}
