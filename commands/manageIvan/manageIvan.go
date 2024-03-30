package manageIvan

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

type state struct {
	// the user to ban
	User *discordgo.User
	// amount of time in seconds before ban
	CountDownTime int
}

var manageIvan = command.Command{
	ApplicationCommand: func() *discordgo.ApplicationCommand {
		defaultManageMessagesPermission := int64(discordgo.PermissionManageMessages)
		minValue := float64(5)

		return &discordgo.ApplicationCommand{
			Version:                  "2.0.0",
			Name:                     "manageivan",
			DefaultMemberPermissions: &defaultManageMessagesPermission,
			Description:              "Command to help the management of Ivan",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionUser,
					Name:         "user",
					Description:  "Ivan account to ban",
					Required:     true,
					Autocomplete: false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "countdown",
					Description: "countdown till Ivan is banned. DEFAULT 15 seconds",
					Required:    false,
					MinValue:    &minValue,
					MaxValue:    30,
				},
			},
		}
	},
	InteractionRespond: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		optionsMap := utils.ParseUserOptions(sess, i)
		countDown := optionsMap["countdown"]

		// set default value to 15, is not set by user
		// also keep track of countdown time in state
		if countDown == nil {
			ivanBanState.CountDownTime = int(15)
		} else {
			ivanBanState.CountDownTime = int(countDown.IntValue())
		}

		// store the user to ban in state
		ivanBanState.User = optionsMap["user"].UserValue(sess)

		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Confirm if you want to ban <@%s> in %d seconds", ivanBanState.User.ID, ivanBanState.CountDownTime),
				Flags:   discordgo.MessageFlagsEphemeral,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							createBanButton(false),
							createDontBanButton(false),
						},
					},
				},
			},
		})

		if err != nil {
			utils.SendErrorMessage(sess, i, err.Error())
			log.Println(err)
		}
		return "Select menu sent", nil
	},
	Components: []struct {
		ComponentID      string
		ComponentHandler command.HandlerFunc
	}{
		{
			ComponentID:      startBanProcessID,
			ComponentHandler: startBanningIvan,
		},
		{
			ComponentID:      dontBanIvanID,
			ComponentHandler: dontBanButtonHandler,
		},
	},
}

var config = utils.LoadConfig()

// custom ID for banning ivan
var startBanProcessID = "start_ivan_ban"
var dontBanIvanID = "dont_ban_ivan"

var ivanBanState = state{}

// bans ivan accounts with a countdown
type ManageIvan struct{}

// startBanningIvan handles the interaction countdown to ban a user
func startBanningIvan(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	// change original ephemeral message to command executor
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "Sequence initiated...",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						createBanButton(true),
						createDontBanButton(true),
					},
				},
			},
		},
	})

	if err != nil {
		return "", err
	}

	// keep track of all sent messages so we cant delete them later
	sentMessages := []*discordgo.Message{}
	messages := GenerateMessages(ivanBanState.CountDownTime)

	// start count down
	for _, message := range messages {
		mess, err := sess.ChannelMessageSend(i.ChannelID, message.message)
		if err != nil {
			return "", err
		}
		sentMessages = append(sentMessages, mess)
		time.Sleep(message.countDownTime * time.Second)
	}

	if !config.Development {
		err = sess.GuildBanCreateWithReason(i.GuildID, ivanBanState.User.ID, "Ivan", 0)
		if err != nil {
			return "", err
		}
	}

	// time.Sleep(5 * time.Second)
	// send embed that user has been banned
	_, err = sess.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embed: &discordgo.MessageEmbed{
			URL:         "",
			Type:        "",
			Title:       "Ivan Ban",
			Description: fmt.Sprintf("<@%s> HAS BEEN BANNED", ivanBanState.User.ID),
			Timestamp:   "",
			Color:       0,
		},
	})

	if err != nil {
		log.Println(err)
	}

	time.Sleep(5 * time.Second)
	// clean up all messages
	utils.DeleteAllMessages(sess, i, sentMessages)
	return fmt.Sprintf("<@%s> HAS BEEN BANNED", ivanBanState.User.ID), nil
}

// dontBanButtonHandler handle when the dont ban button is pushed
func dontBanButtonHandler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Okay, <@%s> will not be banned... :(", ivanBanState.User.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						createBanButton(true),
						createDontBanButton(true),
					},
				},
			},
		},
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Okay, <@%s> will not be banned... :(", ivanBanState.User.String()), nil
}

// createBanButton create a ban button
func createBanButton(disable bool) discordgo.Button {
	return discordgo.Button{
		Label:    "Ban",
		Style:    discordgo.DangerButton,
		Disabled: disable,
		Emoji: discordgo.ComponentEmoji{
			Name:     "✅",
			Animated: false,
		},
		CustomID: startBanProcessID,
	}
}

// createDontBanButton create a dont ban button
func createDontBanButton(disable bool) discordgo.Button {
	return discordgo.Button{
		Label:    "Dont ban...",
		Style:    discordgo.PrimaryButton,
		Disabled: disable,
		Emoji: discordgo.ComponentEmoji{
			Name:     "❌",
			Animated: false,
		},
		CustomID: dontBanIvanID,
	}
}

type CountDownMessage struct {
	// the message to send
	message string
	// length of the count down in seconds
	countDownTime time.Duration
}

// GenerateMessages generates an array of messages for the count down, based on the length of the countDownTime
// return an array of countDownMessage for the count down
func GenerateMessages(countDownTime int) []CountDownMessage {
	messages := make([]CountDownMessage, 0)
	// if the user picks 5 secs as the count down time, then dont bother counting down
	for sec := countDownTime; sec > 0; sec = sec - 5 {
		// if there are 5 seconds or less left, ask for last words
		if sec <= 5 {
			messages = append(messages, CountDownMessage{
				message:       fmt.Sprintf("Any last words? <@%s>\nTime till ban: %ds", ivanBanState.User.ID, sec),
				countDownTime: time.Duration(sec),
			})
		} else if sec == countDownTime {
			// the very first message should be different
			messages = append(messages, CountDownMessage{
				message:       fmt.Sprintf("Count down started to ban <@%s>: %ds", ivanBanState.User.ID, sec),
				countDownTime: time.Duration(5),
			})
		} else {
			// other wise send normal count down time
			messages = append(messages, CountDownMessage{
				message:       fmt.Sprintf("Time till ban: %vs", sec),
				countDownTime: time.Duration(5),
			})
		}
	}
	return messages
}

func init() {
	command.Register(manageIvan)
}
