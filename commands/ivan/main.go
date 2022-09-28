package ivan

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

var CommandObj = commands.CommandStruct{
	Obj:     obj,
	Handler: handler,
}

type emote struct {
	Name         string `json:"name"`
	FileLocation string `json:"fileLocation"`
}

func obj() *discordgo.ApplicationCommand {
	maxLength := float64(1000)
	emotes := getAllEmotes()
	obj := &discordgo.ApplicationCommand{
		Version:     "1.0",
		Name:        "ivan",
		Description: "A command for all things Ivan",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "list",
				Description: "Length of the list of Ivan acocounts that has been banned",
				MaxValue:    maxLength,
				Required:    false,
			},
		},
	}

	// new custom ivan emotes option
	newOption := discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "emote",
		Description: "Send custom Ivan emotes",
		Required:    false,
		Choices:     []*discordgo.ApplicationCommandOptionChoice{},
	}
	for _, emote := range emotes {
		newOption.Choices = append(newOption.Choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  emote.Name,
			Value: emote.Name,
		})
	}
	obj.Options = append(obj.Options, &newOption)

	fmt.Println(fmt.Sprintf("obj obj: %+v", obj.Options[1].Choices)) // __AUTO_GENERATED_PRINT_VAR__
	return obj
}

func handler(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	utils.DeferReply(sess, i.Interaction)
	optionMap := utils.ParseUserOptions(sess, i)

	if optionMap["list"] != nil {
		listLength := optionMap["list"]
		bans, err := sess.GuildBans(i.GuildID, int(listLength.IntValue()), "", "")
		bans = filterIvanBans(bans)

		if err != nil {
			log.Fatal(err)
		}

		_, err = sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
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
	} else if optionMap["emote"] != nil {
		emotes := getAllEmotes()
		fmt.Println(fmt.Sprintf("handler emotes: %v", emotes)) // __AUTO_GENERATED_PRINT_VAR__
		f, err := os.Open("./test_img.jpg")
		if err != nil {
			utils.EditErrorResponse(sess, i, err.Error())
			log.Fatal(err)
		}

		_, err = sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Files: []*discordgo.File{
				{
					Name:        "test.png",
					ContentType: "image/jpg",
					Reader:      f,
				},
			},
		})

		if err != nil {
			log.Fatal(err)
		}
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
	str += fmt.Sprintf("**Total accounts: %d**", len(list))
	return str
}

// filterIvanBans filter out all the ivan bans, and return a new []*discordgo.GuildBan
func filterIvanBans(bans []*discordgo.GuildBan) []*discordgo.GuildBan {
	list := make([]*discordgo.GuildBan, 0)

	for _, ban := range bans {
		if strings.Contains(strings.ToLower(ban.Reason), "ivan") {
			list = append(list, ban)
		}
	}
	return list
}

func getAllEmotes() []emote {
	// map of name: fileLocation for all emotes
	emotes := make([]emote, 0)
	config := utils.LoadConfig()

	for _, emote := range config.Ivan.Emotes {
		emotes = append(emotes, emote)
	}

	return emotes

}
