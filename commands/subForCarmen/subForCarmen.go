package subforcarmen

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

var CommandObj = commands.CommandStruct{
	Name: "subforcarmen",
	Obj:  obj,
	CommandHandler: func(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	},
	Components: []struct {
		ComponentID      string
		ComponentHandler commands.HandlerFunc
	}{},
}

type state struct {
	// time of the last notification
	LastNotificationTime time.Time
	// time of the last Carmen message
	LastMessageTime time.Time
}

var carmenState = state{}

func obj() *discordgo.ApplicationCommand {
	c := utils.LoadConfig()
	return &discordgo.ApplicationCommand{
		GuildID:     c.SubForCarmen.GuildID,
		Version:     "1.0.0",
		Name:        "subforcarmen",
		Description: "subscribe to get notified of Carmen drama in discord",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "subscribe",
				Type:         discordgo.ApplicationCommandOptionBoolean,
				Description:  "subscribe for Carmen",
				ChannelTypes: []discordgo.ChannelType{},
				Required:     true,
				Options:      []*discordgo.ApplicationCommandOption{},
				Autocomplete: false,
				Choices:      []*discordgo.ApplicationCommandOptionChoice{},
			},
		},
	}
}

// CheckMessage Checks a discord message to see if it's SubForCarmen author
// Returns true if the message is from SubForCarmen author. False other wise
func CheckMessage(mess *discordgo.Message) bool {
	c := utils.LoadConfig()
	if mess.Author.ID == c.SubForCarmen.CarmenID {
		carmenState.LastMessageTime = mess.Timestamp
		return true
	}
	return false
}

// IsCoolDown checks if a message is within a cool down period
// mess: the current message to check
// return true if it within cool down period, false other wise
func IsCoolDown(mess *discordgo.Message) bool {
	c := utils.LoadConfig()
	// get time difference between last message time and current message time
	timeDiff := carmenState.LastMessageTime.Sub(mess.Timestamp)
	// if time difference is outside the range of cool down
	if timeDiff.Minutes() > float64(c.SubForCarmen.CoolDown) {
		return false
	}
	return true
}

// IncreaseCounter increases the message counter of the current message's time is within 5 mins of the last message
// return true if the counter is increased. False other wise
func IncreaseCounter(mess *discordgo.Message) bool {
	// if the message is from SubForCarmen author, and its been longer than the cool down period
	if CheckMessage(mess) && !IsCoolDown(mess) {
		// send message
	}
	return true
}
