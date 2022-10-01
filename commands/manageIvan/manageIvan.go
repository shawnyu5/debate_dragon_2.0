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

// custom ID for banning ivan
var startBanProcessID = "start_ivan_ban"
var dontBanIvanID = "dont_ban_ivan"
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
	},
}

func obj() *discordgo.ApplicationCommand {
	defaultManageServerPermission := int64(discordgo.PermissionManageServer)
	minValue := float64(5)

	return &discordgo.ApplicationCommand{
		Version:                  "1.0",
		Name:                     "manageivan",
		DefaultMemberPermissions: &defaultManageServerPermission,
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
			Content: fmt.Sprintf("Confirm if you ban <@%s> in %d seconds", ivanBanState.User.ID, ivanBanState.CountDownTime),
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
}

// startBanningIvan handles the countdown to ban ivan
func startBanningIvan(sess *discordgo.Session, i *discordgo.InteractionCreate) {
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
		log.Println(err)
		return
	}

	// keep track of all sent messages so we cant delete them later
	sentMessages := []*discordgo.Message{}
	messages := generateMessages(ivanBanState.CountDownTime)
	fmt.Println(fmt.Sprintf("startBanningIvan ivanBanState.CountDownTime: %v", ivanBanState.CountDownTime)) // __AUTO_GENERATED_PRINT_VAR__

	// start count down
	for _, message := range messages {
		mess, err := sess.ChannelMessageSend(i.ChannelID, message)
		if err != nil {
			log.Println(err)
			return
		}
		sentMessages = append(sentMessages, mess)
		time.Sleep(5 * time.Second)
	}

	// mess, err := sess.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Count down started: %ds", ivanBanState.CountDownTime))
	// if err != nil {
	// log.Println(err)
	// return
	// }

	// if err != nil {
	// log.Println(err)
	// return
	// }

	// time.Sleep(5 * time.Second)
	// for sec := ivanBanState.CountDownTime - 5; sec >= 0; sec = sec - 5 {
	// fmt.Println(fmt.Sprintf("startBanningIvan j: %v", sec)) // __AUTO_GENERATED_PRINT_VAR__

	// // change the message if there are 5 seconds or less left in the countdown and pass control to another function
	// if sec <= 5 {
	// mess, err = sess.ChannelMessageEdit(i.ChannelID, mess.ID, fmt.Sprintf("Any last words? <@%s>\nTime till ban: %ds", ivanBanState.User.ID, sec))
	// if err != nil {
	// log.Println(err)
	// return
	// }
	// break
	// } else {
	// mess, err = sess.ChannelMessageEdit(i.ChannelID, mess.ID, fmt.Sprintf("Time till ban: %vs", sec))
	// if err != nil {
	// log.Println(err)
	// return
	// }
	// }

	// time.Sleep(5 * time.Second)
	// }

	fmt.Println("startBanningIvan after for loop") // __AUTO_GENERATED_PRINTF__

	time.Sleep(5 * time.Second)
	// sess.ChannelMessageSend(i.ChannelID, "USER HAS BEEN BANNED")
	_, err = sess.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embed: &discordgo.MessageEmbed{
			URL:         "",
			Type:        "",
			Title:       "Ivan Ban",
			Description: fmt.Sprintf("<@%s> HAS BEEN BANNED", ivanBanState.User.ID),
			Timestamp:   "",
			Color:       0,
			Footer:      &discordgo.MessageEmbedFooter{},
			Image:       &discordgo.MessageEmbedImage{},
			Thumbnail:   &discordgo.MessageEmbedThumbnail{},
			Video:       &discordgo.MessageEmbedVideo{},
			Provider:    &discordgo.MessageEmbedProvider{},
			Author:      &discordgo.MessageEmbedAuthor{},
			Fields:      []*discordgo.MessageEmbedField{},
		},
	})

	if err != nil {
		log.Println(err)
	}

	time.Sleep(5 * time.Second)
	// clean up all messages
	for _, message := range sentMessages {
		go func(mess *discordgo.Message) {
			err := sess.ChannelMessageDelete(i.ChannelID, mess.ID)
			if err != nil {
				log.Println(err)
				return
			}

		}(message)
	}

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
					},
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
	}

}

// createBanButton create a ban button
func createBanButton(disable bool) discordgo.Button {
	return discordgo.Button{
		Label:    "Ban!",
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

// generateMessages generates an array of messages for the count down, based on the length of the countDownTime
// return an array of messages for the count down
func generateMessages(countDownTime int) []string {
	messages := make([]string, 0)
	for sec := countDownTime; sec > 0; sec = sec - 5 {
		// the first message should be different
		if sec == countDownTime {
			messages = append(messages, fmt.Sprintf("Count down started: %ds", sec))
		} else if sec <= 5 {
			// if there are 5 seconds or less left, ask for last words
			messages = append(messages, fmt.Sprintf("Any last words? <@%s>\nTime till ban: %ds", ivanBanState.User.ID, sec))
		} else {
			// other wise send normal count down time
			messages = append(messages, fmt.Sprintf("Time till ban: %vs", sec))
		}
	}
	return messages
}
