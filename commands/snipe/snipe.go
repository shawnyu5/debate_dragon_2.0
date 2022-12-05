package snipe

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
)

// command to get the last deleted message
type Snipe struct{}

// map of guild id: { message id : discord message }
var AllMessages = make(map[string]map[string]discordgo.Message)

// Components implements commands.Command
func (Snipe) Components() []commands.Component {
	return nil
}

// Def implements commands.Command
func (Snipe) Def() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Version:                  "1.0.0",
		Type:                     0,
		Name:                     "snipe",
		Description:              "Get the contents of the last deleted message",
		DescriptionLocalizations: &map[discordgo.Locale]string{},
		Options:                  []*discordgo.ApplicationCommandOption{},
	}
}

// Handler implements commands.Command
func (Snipe) Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	panic("unimplemented")
}

// trackMessage adds a message to allMessages map
// mess: the message to track
func TrackMessage(mess *discordgo.MessageCreate) {
	// if the map doesn't exist, initialize it
	if AllMessages[mess.GuildID] == nil {
		AllMessages[mess.GuildID] = map[string]discordgo.Message{}
	}
	// add message to map
	AllMessages[mess.GuildID][mess.ID] = *mess.Message

	// store the messages to be removed
	oldestMessage := discordgo.Message{}
	oldestTimeDuration := time.Duration(0)

	// if we have more than 100 messages stored for this guild, then remove the oldest message
	if len(AllMessages[mess.GuildID]) > 100 {
		for _, message := range AllMessages[mess.GuildID] {
			timeDiff := time.Since(message.Timestamp.UTC())
			if oldestTimeDuration < timeDiff {
				oldestTimeDuration = timeDiff
				oldestMessage = message
			}
		}
		delete(AllMessages[mess.GuildID], oldestMessage.ID)
	}
}
