package releasenotes

import (
	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/command"
)

var releaseNotes = command.Command{
	ApplicationCommand: func() *discordgo.ApplicationCommand {
		return &discordgo.ApplicationCommand{
			Version:     "1.0.0",
			Name:        "release-notes",
			Description: "Get the release notes for this release",
		}
	},
	InteractionRespond: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "Discord Bot: Now 200% Smarter (We Promise)",
						Description: `This release marks the beginning of a smarter debate dragon. Now you can talk to him!

											Use the ` + "`/ask`" + ` command to ask debate dragon a question, or engage in a furious debate!

											**_Disclaimer_: I, the author, do not take responsibility for anything this bot says or does—it's basically an unsupervised toddler with internet access. Any opinions it expresses are its own chaotic creations**`,
						Author: &discordgo.MessageEmbedAuthor{
							Name: "themagicguy",
							// TODO: consider adding the icon here
							IconURL: "",
						},
					},
				},
			},
		})
		return "Release notes sent", err
	},
}

func init() {
	command.Register(releaseNotes)
}
