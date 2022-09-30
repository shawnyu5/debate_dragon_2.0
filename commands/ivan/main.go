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

	// dynamically add new custom ivan emotes option
	newOption := discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "emote",
		Description: "Send custom Ivan emotes",
		Required:    false,
		Choices:     []*discordgo.ApplicationCommandOptionChoice{},
	}
	// add the emote names based on the names defined in config.json
	for name := range emotes {
		newOption.Choices = append(newOption.Choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: name,
		})
	}
	obj.Options = append(obj.Options, &newOption)
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
		chosenEmote := optionMap["emote"].StringValue()

		// if the chosen emote does not exist, send error
		if emotes[chosenEmote] == "" {
			utils.SendErrorMessage(sess, i, "Your chosen emote does not exist...")
		}

		f, err := os.Open(emotes[chosenEmote])

		if err != nil {
			utils.SendErrorMessage(sess, i, err.Error())
			log.Fatal(err)
		}

		_, err = sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Files: []*discordgo.File{
				{
					Name:        chosenEmote + ".png",
					ContentType: "image/png",
					Reader:      f,
				},
			},
		})

		if err != nil {
			utils.SendErrorMessage(sess, i, err.Error())
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

// getAllEmotes returns a list of all the emotes in a map, as specified in config.json. Of name: fileLocation key value pairs
func getAllEmotes() map[string]string {
	// map of name: fileLocation for all emotes
	emotes := make(map[string]string)
	config := utils.LoadConfig()

	for _, emote := range config.Ivan.Emotes {
		emotes[emote.Name] = emote.FileLocation
	}
	return emotes
}
