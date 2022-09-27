package ivan

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

var CommandObj = commands.CommandStruct{
	Obj:     obj,
	Handler: handler,
}

func obj() *discordgo.ApplicationCommand {
	maxLength := float64(1000)
	return &discordgo.ApplicationCommand{
		Version:     "1.0",
		Name:        "ivan",
		Description: "A command for all things Ivan",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "list",
				Description: "Length of the list of Ivan acocounts that has been banned",
				MaxValue:    maxLength,
				Required:    true,
			},
		},
	}
}

func handler(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	utils.DeferReply(sess, i.Interaction)
	optionMap := utils.ParseUserOptions(sess, i)
	listLength := optionMap["list"]

	bans, err := sess.GuildBans(i.GuildID, int(listLength.IntValue()), "", "")
	if err != nil {
		log.Fatal(err)
	}
	_, err = sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		// Content: new(string),
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "All banned ivan users",
				Description: formatList(bans),
				Color:       0,
			},
		},
	},
	)
	if err != nil {
		log.Fatal(err)
	}
}

// formatList formats an array of GuildBans into a bullet
func formatList(list []*discordgo.GuildBan) string {
	str := ""
	for _, item := range list {
		str += fmt.Sprintf("- %s\n", item.User.Username)
	}

	if str == "" {
		return "No banned Ivan users"
	}
	return str
}
