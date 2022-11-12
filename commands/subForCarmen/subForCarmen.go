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
	Name:    "subforcarmen",
	Obj:     obj,
	Handler: handler,
	// Components: []struct {
	// ComponentID      string
	// ComponentHandler commands.HandlerFunc
	// }{},
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
	// start off with a message in the past, so it will trigger a notification
	LastNotificationTime: time.Now().Add(time.Duration(-19) * time.Hour),
	LastMessageTime:      time.Now(),
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

func handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	userOptions := utils.ParseUserOptions(sess, i)
	c := utils.LoadConfig()
	// if subscribe, give user sub role
	if userOptions["subscribe"].BoolValue() {
		fmt.Printf("handler i.GuildID: %v\n", i.GuildID) // __AUTO_GENERATED_PRINT_VAR__
		err := sess.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, c.SubForCarmen.SubscribersRoleID)
		if err != nil {
			return "", err
		}
		content := "You have subscribed for Carmen. Congrats!!!"
		err = sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				TTS:     false,
				Content: content,
			},
		})
		if err != nil {
			return "", err
		}
	} else { // else remove sub role
		err := sess.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, c.SubForCarmen.SubscribersRoleID)
		if err != nil {
			return "", err
		}
		content := "You have unsubscribed for Carmen. Sorry to see you go..."
		err = sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				TTS:     false,
				Content: content,
			},
		})
		if err != nil {
			return "", err
		}
		return content, nil
	}
	return "", nil
}

// Listen Checks a discord message to see if it's SubForCarmen author. And does the needed actions if it is
// sess    : the discord session
// mess    : the discord message
// guildID : the guild ID to check
// Return  : true if a notification is sent. False other wise
func Listen(sess *discordgo.Session, mess *discordgo.Message) bool {
	c := utils.LoadConfig()
	if !IsValidMessage(mess) { // If the message is not valid
		log.Println("(subForCarmen)not a valid message")
		return false
	} else if IsIgnoredChannel(mess.ChannelID) {
		log.Println("(subForCarmen)Channel in ignore list, ignoring")
		return false
	} else if IsCoolDown(mess) { // if we are within cool down period
		log.Println("(subForCarmen)Within cool down period")
		return false
	}

	IncreaseCounter(mess) // increase counter if all above condition is met

	if !ShouldTriggerNotification(c.SubForCarmen.MessageLimit) {
		log.Println("Not enough messages to trigger a notification")
		return false
	} else {
		err := SendNotification(sess, mess.ChannelID, c.SubForCarmen.SubscribersRoleID)
		// update last notification time
		CarmenState.LastNotificationTime = time.Now()
		// reset counter
		CarmenState.Counter = 0
		CarmenState.LastMessageTime = mess.Timestamp
		if err != nil {
			log.Println(err)
			return false
		}
	}
	return true

}

// IsValidMessage checks if this is a message we should be parsing
// return: true if this is a message we should be parsing. False other wise
func IsValidMessage(mess *discordgo.Message) bool {
	c := utils.LoadConfig()
	// checks if an element is in a slice
	// arr: the slice to check
	// e  : the element to check
	contains := func(arr []string, e string) bool {
		for _, v := range arr {
			if v == e {
				return true
			}
		}
		return false
	}

	// if the author is carmen, message is in correct guild, and the channel is not being ignored, return true
	return mess.Author.ID == c.SubForCarmen.CarmenID &&
		mess.GuildID == c.SubForCarmen.GuildID &&
		!contains(c.SubForCarmen.IgnoredChannels, mess.ChannelID)
}

// IsCoolDown checks if a message is within a cool down period
// mess          : the current message to check
// coolDownPeriod: the length of cool down period in minues
// return        : true if the message is within cool down period, false other wise
func IsCoolDown(mess *discordgo.Message) bool {
	c := utils.LoadConfig()
	// no cool down period in development
	if c.Development {
		return false
	}
	// get time difference between last notification time and current message time
	timeDiff := mess.Timestamp.Sub(CarmenState.LastNotificationTime)

	// return if time difference is within cool down period
	return timeDiff.Minutes() <= float64(c.SubForCarmen.CoolDown)
}

// IncreaseCounter increases the message counter if the current message's time is within 5 mins of the last message
// mess  : the current message
// return: true if the counter is increased. False if counter is reset to 0
func IncreaseCounter(mess *discordgo.Message) bool {
	c := utils.LoadConfig()
	// time difference between last message time and current message time
	timeDiff := mess.Timestamp.Sub(CarmenState.LastMessageTime)
	fmt.Printf("IncreaseCounter timeDiff: %v\n", timeDiff) // __AUTO_GENERATED_PRINT_VAR__

	// increase counter if current message is sent within 6 mins of last message
	if timeDiff.Minutes() <= float64(6) {
		CarmenState.Counter++
	}

	// reset counter if counter has reached message limit
	if CarmenState.Counter >= c.SubForCarmen.MessageLimit {
		CarmenState.Counter = 0
		log.Println("Resetting counter")
		return false
	}
	return true
}

// ShouldTriggerNotification Checks if state counter has reached message limit
// messageLimit: the message limit to trigger a notification
// return      : true if state counter has reached messageLimit. False other wise
func ShouldTriggerNotification(messageLimit int) bool {
	fmt.Printf("ShouldTriggerNotification CarmenState.Counter: %v\n", CarmenState.Counter) // __AUTO_GENERATED_PRINT_VAR__
	return CarmenState.Counter == messageLimit
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

// IsIgnoredChannel check if the message is from a channel that should be ignored
// chann: the channel ID to look for
// return: true if the channel is in the ignore list. False other wise
func IsIgnoredChannel(chanID string) bool {
	c := utils.LoadConfig()
	for _, channel := range c.SubForCarmen.IgnoredChannels {
		if chanID == channel {
			return true
		}
	}
	return false
}
