package rmp

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

const profSelectMenuID = "prof select menu"

var CommandObj = commands.CommandStruct{
	Name:           "rmp",
	Obj:            obj,
	CommandHandler: handler,
	Components: []struct {
		ComponentID      string
		ComponentHandler commands.HandlerFunc
	}{
		{
			ComponentID:      profSelectMenuID,
			ComponentHandler: menuHandler,
		},
	},
}

func obj() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Version:     "1.0.0",
		Name:        "rmp",
		Description: "Get reviews from rate my prof",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "profname",
				Description:  "name of the professor to look up",
				Required:     true,
				Autocomplete: false,
			},
		},
	}
}

func handler(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	options := utils.ParseUserOptions(sess, i)
	profName := options["profname"].StringValue()
	searchResult := SearchRmpProfByName(profName)
	senecaProfs := FilterSenecaProfs(searchResult)
	j, _ := json.Marshal(senecaProfs)
	fmt.Printf("handler prof: %+v\n", string(j)) // __AUTO_GENERATED_PRINT_VAR__

	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Content: "hello",
			Flags: 0,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						createSelectMenu(senecaProfs, false),
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
}

// menuHandler handles when an option is selected in the select menu
func menuHandler(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("menu handler")
	data := i.MessageComponentData()
	fmt.Printf("menuHandler data: %+v\n", data) // __AUTO_GENERATED_PRINT_VAR__
	// options := utils.ParseUserOptions(sess, i)
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: data.Values[0],
		},
	})
	if err != nil {
		panic(err)
	}
}

// createSelectMenu create a select menu containing the profs
func createSelectMenu(profs []ProfNode, disabled bool) discordgo.SelectMenu {
	menu := discordgo.SelectMenu{
		CustomID:    profSelectMenuID,
		Placeholder: "Please select a prof",
		Options:     []discordgo.SelectMenuOption{},
		Disabled:    disabled,
	}

	// add all profs as an option to the select menu
	for _, prof := range profs {
		// convert id to a string, so we can search by the id later to get the rating of a prof
		profID := strconv.Itoa(int(prof.LegacyID))
		option := discordgo.SelectMenuOption{
			Label:       fmt.Sprintf("%s %s", prof.FirstName, prof.LastName),
			Value:       profID,
			Description: fmt.Sprintf("Department: %s", prof.Department),
			Default:     false,
		}
		menu.Options = append(menu.Options, option)
	}
	fmt.Printf("createSelectMenu menu: %+v\n", menu) // __AUTO_GENERATED_PRINT_VAR__
	return menu
}
