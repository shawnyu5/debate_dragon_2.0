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

type State struct {
	// time of the last notification
	LastNotificationTime time.Time
	// time of the last Carmen message
	LastMessageTime time.Time
	// message counter
	Counter int
}

var CarmenState = State{}

// updates the global state value
// func (s *State) UpdateState() {
// CarmenState = *s
// }

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

// CheckMessageAuthor Checks a discord message to see if it's SubForCarmen author
// mess    : the message to check for
// authorID: the author ID to check against
// Return  : true if the message is from SubForCarmen author. False other wise
func CheckMessageAuthor(mess *discordgo.Message, authorID string) bool {
	if mess.Author.ID == authorID {
		CarmenState.LastMessageTime = mess.Timestamp
		return true
	}
	return false
}

// IsCoolDown checks if a message is within a cool down period
// mess          : the current message to check
// coolDownPeriod: the length of cool down period in minues
// return        : true if the message is within cool down period, false other wise
func IsCoolDown(mess *discordgo.Message, coolDownPeriod float64) bool {
	// get time difference between last message time and current message time
	timeDiff := mess.Timestamp.Sub(CarmenState.LastMessageTime)

	// if time difference is within cool down period, update last message time stamp
	if timeDiff.Minutes() <= coolDownPeriod {
		CarmenState.LastMessageTime = mess.Timestamp
		return true
	}
	return false
}

// IncreaseCounter increases the message counter of the current message's time is within 5 mins of the last message. Else resets counter to 0
// return: true if the counter is increased. False other wise
func IncreaseCounter(mess *discordgo.Message) bool {
	// if the message is from SubForCarmen author, and its been longer than the cool down period
	timeDiff := mess.Timestamp.Sub(CarmenState.LastMessageTime)

	// timeDiff := time.Since(mess.Timestamp)
	// increase counter if current message is sent within 5 mins of last message
	if timeDiff.Minutes() <= float64(6) {
		CarmenState.Counter++
		CarmenState.LastMessageTime = mess.Timestamp
		return true
	}
	CarmenState.Counter = 0
	return false
}
