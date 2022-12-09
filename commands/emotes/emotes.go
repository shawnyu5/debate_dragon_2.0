package emotes

import (
	"github.com/bwmarrin/discordgo"
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
				Autocomplete: false,
			},
		},
	}

	emotes := GetEmotes()
	for name, url := range emotes {
		def.Options[0].Choices = append(def.Options[0].Choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: url,
		})
	}

	return def
}

// Handler implements commands.Command
func (Emotes) Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	input := utils.ParseUserOptions(sess, i)
	utils.DeferReply(sess, i.Interaction)
	emoteUrl := input["name"].StringValue()
	_, err := sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &emoteUrl,
	})
	if err != nil {
		return "", err
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
