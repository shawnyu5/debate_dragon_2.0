package rmp

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

const profSelectMenuID = "prof select menu"

type state struct {
	// all seneca profs returned from RMP
	allProfs []ProfNode
	// user selected prof node
	selectedProf ProfNode
}

var rmpState = state{}

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

	// if not profs are found, return message
	if len(senecaProfs) == 0 {
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No profs by that name is at seneca...",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			panic(err)
		}

		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
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
	} else {
		// since there is only 1 prof, we just get the first element of the array
		prof := rmpState.allProfs[0]
		// respond with prof information
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:    "",
				Components: []discordgo.MessageComponent{},
				Embeds: []*discordgo.MessageEmbed{
					{
						Type:        "",
						Title:       fmt.Sprintf("%s %s", prof.FirstName, prof.LastName),
						Description: profDescription(prof),
						Color:       0,
						Footer: &discordgo.MessageEmbedFooter{
							Text:         "Information retrieved from ratemyprof.com",
							IconURL:      "https://pbs.twimg.com/profile_images/1146077191043788800/hG1lAGm9_400x400.png",
							ProxyIconURL: "",
						},
						// Image:     &discordgo.MessageEmbedImage{},
						// Thumbnail: &discordgo.MessageEmbedThumbnail{},
						// Video:     &discordgo.MessageEmbedVideo{},
						// Provider:  &discordgo.MessageEmbedProvider{},
						Author: &discordgo.MessageEmbedAuthor{
							URL:          "",
							Name:         "brought to you by your mom TM",
							IconURL:      "",
							ProxyIconURL: "",
						},
						Fields: []*discordgo.MessageEmbedField{},
					},
				},
			},
		})
		if err != nil {
			panic(err)
		}
	}
}

// menuHandler handles when an option is selected in the select menu
func menuHandler(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	// id of the prof selected by the user
	selectedProfID := data.Values[0]
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: selectedProfID,
		},
	})
	if err != nil {
		panic(err)
	}

	// find the prof that was selected by the user, and store it in the global state
	for _, prof := range rmpState.allProfs {
		if prof.ID == selectedProfID {
			rmpState.selectedProf = prof
		}
	}
}

// createSelectMenu create a select menu containing the profs
func createSelectMenu(profs []ProfNode, disabled bool) discordgo.SelectMenu {
	MinValues := 1
	menu := discordgo.SelectMenu{
		CustomID:    profSelectMenuID,
		Placeholder: "Please select a prof",
		Options:     []discordgo.SelectMenuOption{},
		Disabled:    disabled,
		MinValues:   &MinValues,
	}

	// add all profs as an option to the select menu
	for _, prof := range profs {
		// convert id to a string, so we can search by the id later to get the rating of a prof
		option := discordgo.SelectMenuOption{
			Label:       fmt.Sprintf("%s %s", prof.FirstName, prof.LastName),
			Value:       prof.ID,
			Description: fmt.Sprintf("Department: %s", prof.Department),
			Default:     false,
		}
		menu.Options = append(menu.Options, option)
	}
	return menu
}

// profDescription generate a description about a professor
func profDescription(prof ProfNode) string {
	return fmt.Sprintf(`- **Average rating**: %f
- **Average difficulty**: %f
- **Would take again**: %f%%`, prof.AvgRating, prof.AvgDifficulty, prof.WouldTakeAgainPercent)
}
