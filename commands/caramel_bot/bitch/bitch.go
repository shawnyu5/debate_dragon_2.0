package bitch

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

var bitch = command.Command{
	ApplicationCommand: func() *discordgo.ApplicationCommand {
		return &discordgo.ApplicationCommand{
			Name:        "caramel-bot-bitch",
			Description: "Call another user a bitch",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user",
					Description: "The user you would like to call a bitch",
					Type:        discordgo.ApplicationCommandOptionUser,
					Required:    true,
				},
			},
		}
	},
	HandlerFunc: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		var message = ""
		options := utils.ParseUserOptions(sess, i)
		user := options["user"].UserValue(sess)

		if user != nil && user.ID != "246732655373189120" {
			message = "<@" + user.ID + "> is a bitch."
		} else if user != nil {
			message = "<@" + i.Member.User.ID + "> nice try, you're a bitch."
		}

		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})

		if err != nil {
			return "", err
		}

		return fmt.Sprintf("User %s has been called a bitch", user.Username), nil
	},
}

func init() {
	command.Register(bitch)
}
