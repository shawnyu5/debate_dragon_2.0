package manageIvan

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

type state struct {
	// the user to ban
	User *discordgo.User
	// amount of time in seconds before ban
	CountDownTime int
}

var config = utils.LoadConfig()

// custom ID for buttons
var (
	startBanProcessID = "start_ivan_ban"
	dontBanIvanID     = "dont_ban_ivan"
	banJumpScareID    = "ban_jump_scare"
	kickID            = "kick_ivan"
)

var ivanBanState = state{}

type ManageIvan struct{}

// Def implements commands.Command
func (ManageIvan) Def() *discordgo.ApplicationCommand {
	defaultManageMessagesPermission := int64(discordgo.PermissionManageMessages)
	minValue := float64(5)

	return &discordgo.ApplicationCommand{
		Version:                  "1.1.0",
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
}

// Components implements commands.Command
func (ManageIvan) Components() []commands.Component {
	return []commands.Component{
		{
			ComponentID:      startBanProcessID,
			ComponentHandler: handleBan,
		},
		{
			ComponentID:      dontBanIvanID,
			ComponentHandler: handleDontBan,
		},
		{
			ComponentID:      banJumpScareID,
			ComponentHandler: handleJumpScare,
		},
	}
}

// Handler implements commands.Command
func (m ManageIvan) Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
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
				CreateAllButtons(false),
			},
		},
	})

	if err != nil {
		utils.SendErrorMessage(sess, i, err.Error())
		log.Println(err)
	}
	return "Select menu sent", nil
}

// a single message to send during countdown period
type CountDownMessage struct {
	// the message to send
	message string
	// length of the count down in seconds
	countDownTime time.Duration
}

// GenerateMessages generates an array of messages for the count down, based on the length of the countDownTime.
// countDownTime: time of the countdown in seconds.
// return: an array of countDownMessage for the count down.
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
