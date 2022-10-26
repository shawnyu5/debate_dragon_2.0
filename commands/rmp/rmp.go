package rmp

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

const profSelectMenuID = "prof select menu"

type state struct {
	// all seneca profs returned from RMP
	AllSenecaProfs []ProfNode
	// user selected prof node
	SelectedProf ProfNode
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
	rmpState.AllSenecaProfs = senecaProfs

	// if not profs are found, return message
	if len(senecaProfs) == 0 {
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("No profs by the name `%s` is at Seneca...", profName),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			panic(err)
		}

	} else if len(senecaProfs) > 1 {
		// if there is more than 1 prof, respond with select menu
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 0,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							createSelectMenu(rmpState.AllSenecaProfs, false),
						},
					},
				},
			},
		})

		// disable select menu after 3 mins
		time.AfterFunc(2*time.Minute, func() {
			_, err := sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: new(string),
				Components: &[]discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							createSelectMenu(rmpState.AllSenecaProfs, true),
						},
					},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		})

		if err != nil {
			panic(err)
		}
	} else {
		// since there is only 1 prof, we just get the first element of the array
		prof := rmpState.AllSenecaProfs[0]
		// respond with prof information
		err := SendProfInformation(sess, i, prof)
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
	// get the prof node the user selected
	for _, prof := range rmpState.AllSenecaProfs {
		if prof.ID == selectedProfID {
			rmpState.SelectedProf = prof
		}
	}
	err := SendProfInformation(sess, i, rmpState.SelectedProf)
	if err != nil {
		panic(err)
	}
}

// createSelectMenu create a select menu containing the profs
func createSelectMenu(profs []ProfNode, disabled bool) discordgo.SelectMenu {
	MinValues := 1
	menu := discordgo.SelectMenu{
		CustomID:    profSelectMenuID,
		Placeholder: "Please select a prof",
		MinValues:   &MinValues,
		MaxValues:   1,
		Options:     []discordgo.SelectMenuOption{},
		Disabled:    disabled,
	}

	// add all profs as an option to the select menu
	for _, prof := range profs {
		// convert id to a string, so we can search by the id later to get the rating of a prof
		option := discordgo.SelectMenuOption{
			Label:       prof.fullName(),
			Value:       prof.ID,
			Description: fmt.Sprintf("Department: %s", prof.Department),
			Emoji:       discordgo.ComponentEmoji{},
			Default:     false,
		}
		menu.Options = append(menu.Options, option)
	}
	return menu
}

// SendProfInformation reply to an interaction with information about a professor
func SendProfInformation(sess *discordgo.Session, i *discordgo.InteractionCreate, prof ProfNode) error {
	return sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					URL:         prof.rmpURL(),
					Type:        "",
					Title:       fmt.Sprintf("%s %s", prof.FirstName, prof.LastName),
					Description: prof.profDescription(),
					Timestamp:   "",
					Color:       0,
					Footer: &discordgo.MessageEmbedFooter{
						Text:         "Information retrieved from ratemyprof.com",
						IconURL:      "https://pbs.twimg.com/profile_images/1146077191043788800/hG1lAGm9_400x400.png",
						ProxyIconURL: "",
					},
					Image:     &discordgo.MessageEmbedImage{},
					Thumbnail: &discordgo.MessageEmbedThumbnail{},
					Video:     &discordgo.MessageEmbedVideo{},
					Provider:  &discordgo.MessageEmbedProvider{},
					Author: &discordgo.MessageEmbedAuthor{
						URL:          "https://www.youtube.com/watch?v=dQw4w9WgXcQ&ab_channel=RickAstley",
						Name:         fmt.Sprintf("brought to you by @%s's mom TM", sess.State.User.Username),
						IconURL:      "",
						ProxyIconURL: "",
					},
					Fields: []*discordgo.MessageEmbedField{},
				},
			},
		},
	})

}
