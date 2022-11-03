package subforcarmen

import (
	"fmt"
	"log"
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

var CarmenState = State{
	LastMessageTime: time.Now(),
}

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
	c := utils.LoadConfig()

	// no cool down period in development
	if c.Development {
		return false
	}
	// TODO: I dont think this logic is right...
	// get time difference between last message time and current message time
	timeDiff := mess.Timestamp.Sub(CarmenState.LastMessageTime)
	fmt.Printf("IsCoolDown timeDiff: %v\n", timeDiff) // __AUTO_GENERATED_PRINT_VAR__

	// if time difference is within cool down period, update last message time stamp
	if timeDiff.Minutes() <= coolDownPeriod {
		CarmenState.LastMessageTime = mess.Timestamp
		fmt.Println("IsCoolDown") // __AUTO_GENERATED_PRINTF__
		return true
	}
	return false
}

// IncreaseCounter increases the message counter if the current message's time is within 5 mins of the last message. Else resets counter to 0
// mess  : the current message
// return: true if the counter is increased. False other wise
func IncreaseCounter(mess *discordgo.Message) bool {
	c := utils.LoadConfig()
	// time difference between last message time and current message time
	timeDiff := mess.Timestamp.Sub(CarmenState.LastMessageTime)
	fmt.Printf("IncreaseCounter timeDiff: %v\n", timeDiff) // __AUTO_GENERATED_PRINT_VAR__

	// increase counter if current message is sent within 5 mins of last message
	if timeDiff.Minutes() <= float64(6) {
		CarmenState.Counter++
		CarmenState.LastMessageTime = mess.Timestamp
	}
	// reset counter if counter has reached message limit
	if CarmenState.Counter == c.SubForCarmen.MessageLimit {
		CarmenState.Counter = 0
		return true
	}
	return false
}

// ShouldTriggerNotification Checks if state counter has reached message limit
// messageLimit: the message limit to trigger a notification
// return      : true if state counter has reached messageLimit. False other wise
func ShouldTriggerNotification(messageLimit int) bool {
	return CarmenState.Counter >= messageLimit
}

// SendNotification sends a notification to a channel
// sess     : the discord session
// channelID: the channel ID to send the notification to
// subRoleID: the role to ping
func SendNotification(sess *discordgo.Session, channelID string, subRoleID string) error {
	_, err := sess.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Ayo <@&%s> Caramel is rambling again", subRoleID),
	})

	return err
}

