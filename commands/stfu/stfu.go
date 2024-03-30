package stfu

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

var stfu = command.Command{
	ApplicationCommand: func() *discordgo.ApplicationCommand {
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
					MaxValue:     60,
				},
			},
		}
	},
	InteractionRespond: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
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
		if userOptions["user"].UserValue(sess).ID == "903372725605785761" {
			sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "stfu",
							Description: "Can not tell the bot to stfu",
						},
					},
				},
			})
			return "", errors.New("can not tell the bot to stfu")
		}

		State.User = userOptions["user"].UserValue(sess)
		State.InUse = true

		if !State.IsCoolDown {
			sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "stfu",
							Description: fmt.Sprintf("<@%s> will be told to stfu for %s seconds", State.User.ID, State.Length.String()),
						},
					},
				},
			})
			return "stfu sequence initiated", nil
		} else {
			sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "stfu",
							Description: "Within cool down period...",
							Timestamp:   "",
						},
					},
				},
			})
			return "stfu is in cool down", nil
		}
	},
}

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
	// if the command is in use right now
	InUse bool
}

// NewState creates a new stfu state with default values
func NewState() StfuState {
	return StfuState{
		Length:         10 * time.Second,
		IsCoolDown:     false,
		User:           &discordgo.User{},
		InUse:          false,
		CoolDownLength: 30 * time.Second,
	}
}

var State = NewState()

// TellUser tell a user to stfu on every message they send.
// sess: discord session.
// mess: the message to check.
func TellUser(sess *discordgo.Session, mess *discordgo.MessageCreate) {
	if !State.InUse && !State.IsCoolDown {
		log.Debug("stfu is not in use right now")
		return
	} else if mess.Author.ID != State.User.ID {
		log.Debug("not the user to stfu to")
		return
	}

	sess.ChannelMessageSendComplex(mess.ChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("<@%s> stfu", mess.Author.ID),
	})

	// after the stfu length, reset the state and disable telling user to stfu
	time.AfterFunc(State.Length, func() {
		State = NewState()
		State.IsCoolDown = true

		// after the cooldown period, re enable command
		time.AfterFunc(State.CoolDownLength, func() {
			State.IsCoolDown = false
		})
	})
}

func init() {
	command.Register(stfu)
}
