package stfu

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

// tell a user to stfu for a selected period of time
type Stfu struct{}

type StfuState struct {
	// length of time to stfu for
	Length time.Duration
	// if the command is in cool down state
	IsCoolDown bool
	// length of cool down
	CoolDownLength time.Duration
	// the user to stfu
	User *discordgo.User
	// if the command is enabled
	Enable bool
}

// NewState creates a new stfu state with default values
func NewState() StfuState {
	return StfuState{
		Length:         10 * time.Second,
		IsCoolDown:     false,
		User:           &discordgo.User{},
		Enable:         false,
		CoolDownLength: 30 * time.Second,
	}
}

var State = NewState()

// Components implements commands.Command
func (Stfu) Components() []commands.Component {
	return nil
}

// Def implements commands.Command
func (Stfu) Def() *discordgo.ApplicationCommand {
	minLengthValue := float64(5)
	return &discordgo.ApplicationCommand{
		Version:     "1.0.0",
		Type:        0,
		Name:        "stfu",
		Description: "tell a user to stfu for a selected period of time",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "the user to tell stfu to",
				Required:    true,
			},
			{
				Type:         discordgo.ApplicationCommandOptionInteger,
				Name:         "length",
				Description:  "the length of time to do this for",
				ChannelTypes: []discordgo.ChannelType{},
				Autocomplete: false,
				Choices:      []*discordgo.ApplicationCommandOptionChoice{},
				MinValue:     &minLengthValue,
				MaxValue:     30,
				MinLength:    new(int),
				MaxLength:    0,
			},
		},
	}
}

// Handler implements commands.Command
func (Stfu) Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	userOptions := utils.ParseUserOptions(sess, i)
	// if a user set a custom length, use that
	if val, ok := userOptions["length"]; ok {
		length := val.IntValue()
		// parse duration into seconds
		duration, err := time.ParseDuration(strconv.Itoa(int(length)) + "s")
		if err != nil {
			return "", err
		}

		State.Length = duration
	}

	State.User = userOptions["user"].UserValue(sess)
	State.Enable = true

	return "stfu sequence initiated", nil
}

// TellUser tell a user to stfu on every message they send.
// sess: discord session.
// mess: the message to check.
func TellUser(sess *discordgo.Session, mess *discordgo.MessageCreate) {
	if !State.Enable {
		log.Println("stfu is not enabled")
		return
	} else if mess.Author.ID != State.User.ID {
		log.Println("not the user to stfu to")
		return
	}

	sess.ChannelMessageSendComplex(mess.ChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("<@%s> stfu", mess.Author.ID),
	})

	// after the stfu length, reset the state and disable telling user to stfu
	time.AfterFunc(State.Length, func() {
		State = NewState()
	})
}
