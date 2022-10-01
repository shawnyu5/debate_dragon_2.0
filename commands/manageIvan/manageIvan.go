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
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "countdown",
				Description: "countdown till Ivan is banned. DEFAULT 15 seconds",
				Required:    false,
				MinValue:    &minValue,
				MaxValue:    20,
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
	// TODO: fix this.
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
			Content: "Confirm if you would like to proceed",
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

	// start count down
	mess, err := sess.ChannelMessageSend(i.ChannelID, "Count down started: 15s")
	if err != nil {
		log.Println(err)
		return
	}

	if err != nil {
		log.Println(err)
		return
	}

	time.Sleep(5 * time.Second)
	for sec := ivanBanState.CountDownTime - 5; sec >= 0; sec = sec - 5 {
		fmt.Println(fmt.Sprintf("startBanningIvan j: %v", sec)) // __AUTO_GENERATED_PRINT_VAR__

		// change the message if there are 5 seconds or less left in the countdown and pass control to another function
		if sec <= 5 {
			mess, err = sess.ChannelMessageEdit(i.ChannelID, mess.ID, fmt.Sprintf("Any last words? <@%s>\nTime till ban: %ds", ivanBanState.User.ID, sec))
			if err != nil {
				log.Println(err)
				return
			}
			break
		} else {
			mess, err = sess.ChannelMessageEdit(i.ChannelID, mess.ID, fmt.Sprintf("Time till ban: %vs", sec))
			if err != nil {
				log.Println(err)
				return
			}
		}

		time.Sleep(5 * time.Second)

	}

	fmt.Println("startBanningIvan after for loop") // __AUTO_GENERATED_PRINTF__

	time.Sleep(5 * time.Second)
	sess.ChannelMessageSend(i.ChannelID, "USER HAS BEEN BANNED")

	// time.AfterFunc(5*time.Second, func() {
	// })
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
