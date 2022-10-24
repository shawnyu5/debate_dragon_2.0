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

// custom ID for banning ivan
var startBanProcessID = "start_ivan_ban"
var dontBanIvanID = "dont_ban_ivan"
var banJumpScareID = "ban_jump_scare"

var ivanBanState = state{}

var CommandObj = commands.CommandStruct{
	Name:           "manageivan",
	Obj:            obj,
	CommandHandler: commandHandler,
	Components: []struct {
		ComponentID      string
		ComponentHandler utils.HandlerFunc
	}{
		{
			ComponentID:      startBanProcessID,
			ComponentHandler: startBanningIvan,
		},
		{
			ComponentID:      dontBanIvanID,
			ComponentHandler: dontBanIvan,
		},
		{
			ComponentID:      banJumpScareID,
			ComponentHandler: jumpScareBan,
		},
	},
}

// obj returns the command object for `/manageivan` command
func obj() *discordgo.ApplicationCommand {
	defaultManageMessagesPermission := int64(discordgo.PermissionManageMessages)
	minValue := float64(5)

	return &discordgo.ApplicationCommand{
		Version:                  "1.1",
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

// commandHandler the commandHandler for `/manageivan` command
func commandHandler(sess *discordgo.Session, i *discordgo.InteractionCreate) {
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
						createJumpScareButton(false),
					},
				},
			},
		},
	})

	if err != nil {
		utils.SendErrorMessage(sess, i, err.Error())
		log.Println(err)
	}
}

// startBanningIvan handles the interaction countdown to ban a user
func startBanningIvan(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	// change original ephemeral message to disable all buttons and signal the start of sequence
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
						createJumpScareButton(true),
					},
				},
			},
		},
	})

	if err != nil {
		log.Println(err)
		return
	}

	// keep track of all sent messages so we cant delete them later
	sentMessages := []*discordgo.Message{}
	messages := GenerateMessages(ivanBanState.CountDownTime)

	// start count down
	for _, message := range messages {
		mess, err := sess.ChannelMessageSend(i.ChannelID, message.message)
		if err != nil {
			log.Println(err)
			return
		}
		sentMessages = append(sentMessages, mess)
		time.Sleep(message.countDownTime * time.Second)
	}

	if !config.Development {
		fmt.Println("Banning user for real")
		sess.GuildBanCreateWithReason(i.GuildID, ivanBanState.User.ID, "Ivan", 0)
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
}

// dontBanIvan handle with the dont ban button is pushed
func dontBanIvan(sess *discordgo.Session, i *discordgo.InteractionCreate) {
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
						createJumpScareButton(true),
					},
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
	}

}

// jumpScareBan handle with the jump scare button is pushed
func jumpScareBan(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	// change original ephemeral message to command executor
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "Jump scare sequence initiated...",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						createBanButton(true),
						createDontBanButton(true),
						createJumpScareButton(true),
					},
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

	// keep track of all sent messages so we cant delete them later
	sentMessages := []*discordgo.Message{}
	messages := GenerateMessages(ivanBanState.CountDownTime)

	// start count down
	for _, message := range messages {
		mess, err := sess.ChannelMessageSend(i.ChannelID, message.message)
		if err != nil {
			log.Println(err)
			return
		}
		sentMessages = append(sentMessages, mess)
		time.Sleep(message.countDownTime * time.Second)
	}
	mess, err := sess.ChannelMessageSend(i.ChannelID, fmt.Sprintf("jk, we ain't that mean, you will not be banned <@%s>", ivanBanState.User.ID))
	if err != nil {
		fmt.Println(err)
		return
	}
	sentMessages = append(sentMessages, mess)
	time.Sleep(5 * time.Second)

	mess, err = sess.ChannelMessageSend(i.ChannelID, "Good bye now...")
	if err != nil {
		fmt.Println(err)
		return
	}
	sentMessages = append(sentMessages, mess)

	time.Sleep(3 * time.Second)

	// clean up all messages
	utils.DeleteAllMessages(sess, i, sentMessages)
}

// createBanButton create a ban button
func createBanButton(disable bool) discordgo.Button {
	return discordgo.Button{
		Label:    "Ban",
		Style:    discordgo.DangerButton,
		Disabled: disable,
		Emoji:    discordgo.ComponentEmoji{},
		CustomID: startBanProcessID,
	}
}

// createDontBanButton create a dont ban button
func createDontBanButton(disable bool) discordgo.Button {
	return discordgo.Button{
		Label:    "Dont ban...",
		Style:    discordgo.PrimaryButton,
		Disabled: disable,
		Emoji:    discordgo.ComponentEmoji{},
		CustomID: dontBanIvanID,
	}
}

// createJumpScareButton create a jump scare button
func createJumpScareButton(disable bool) discordgo.Button {
	return discordgo.Button{
		Label:    "Jump scare",
		Style:    discordgo.SuccessButton,
		Disabled: disable,
		Emoji:    discordgo.ComponentEmoji{},
		CustomID: banJumpScareID,
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
