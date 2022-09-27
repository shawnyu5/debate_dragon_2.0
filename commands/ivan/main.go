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
	// maxValue := float64(1000.0)
	return &discordgo.ApplicationCommand{
		Version:     "1.0",
		Name:        "ivan",
		Description: "A command for all things Ivan",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "list",
				Description: "List all the Ivan accounts that has been banned",
				// Description: "get all the ivan users that has been banned from this server. Optionally specify the length of the list to retrieve",
				// MaxValue:    maxValue,
				Required: true,
			},
		},
	}
}

func handler(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	optionMap := utils.ParseUserOptions(sess, i)
	listLength := optionMap["list"]

	bans, err := sess.GuildBans(i.GuildID, int(listLength.IntValue()), "", "")
	if err != nil {
		log.Fatal(err)
	}
	// // get all the user that has been banned cuz they were ivan
	// for _, ban := range bans {
	// // check if ban reason contains "ivan"
	// if strings.Contains(strings.ToLower(ban.User.Username), "ivan") {
	// list = append(list, ban.User.Username)
	// }
	// }

	err = sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			TTS:        false,
			Components: []discordgo.MessageComponent{},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "All banned ivan users",
					Description: formatList(bans),
					Color:       0,
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
			Files:           []*discordgo.File{},
			Flags:           0,
			Choices:         []*discordgo.ApplicationCommandOptionChoice{},
			CustomID:        "",
			Title:           "All Banned ivan users",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func formatList(list []*discordgo.GuildBan) string {
	str := ""
	fmt.Println(fmt.Sprintf("formatList list: %v", list)) // __AUTO_GENERATED_PRINT_VAR__
	for _, item := range list {
		str += fmt.Sprintf("- %s\n", item.User.Username)
	}

	if str == "" {
		return "No banned Ivan users"
	}
	return str
}
