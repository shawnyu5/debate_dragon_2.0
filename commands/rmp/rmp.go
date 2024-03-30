package rmp

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

const profSelectMenuID = "prof select menu"

var rmp = command.Command{
	ApplicationCommand: func() *discordgo.ApplicationCommand {
		return &discordgo.ApplicationCommand{
			Name:        "rmp",
			Version:     "2.0.0",
			Description: "Get reviews from rate my prof",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "profname",
					Description:  "name of the professor to look up",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		}
	},
	InteractionRespond: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		options := utils.ParseUserOptions(sess, i)
		profName := options["profname"].StringValue()
		searchResult := SearchRmpProfByName(profName)
		profs := searchResult.Data.Search.Teachers.Edges
		rmpState.AllSenecaProfs = profs

		// if not profs are found, return message
		if len(profs) == 0 {
			err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("No profs by the name `%s` is at Seneca...", profName),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("No profs by the name `%s` is at Seneca...", profName), nil

		} else if len(profs) > 1 {
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

			// disable select menu after 2 mins
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
				return "", err
			}
			return "Multi profs found, select menu sent", nil
		} else {
			// since there is only 1 prof, we just get the first element of the array
			prof := rmpState.AllSenecaProfs[0]
			rmpState.SelectedProf = prof
			// respond with prof information
			err := SendProfInformation(sess, i, prof)
			if err != nil {
				return "", err
			}
			return "end of case", nil
		}

	},
	InteractionApplicationCommandAutocomplete: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		err := autoCompleteHandler(sess, i)
		if err != nil {
			return "", err
		}
		return "Filled in auto complete menu", nil
	},

	Components: []struct {
		ComponentID      string
		ComponentHandler command.HandlerFunc
	}{
		{
			ComponentID:      profSelectMenuID,
			ComponentHandler: menuHandler,
		},
	},
}

type state struct {
	// all seneca profs returned from RMP
	AllSenecaProfs []ProfNode
	// user selected prof node
	SelectedProf ProfNode
}

var rmpState = state{}

// menuHandler handles when an option is selected in the select menu
func menuHandler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	data := i.MessageComponentData()
	// id of the prof selected by the user
	selectedProfID := data.Values[0]
	// get the prof node the user selected
	for _, prof := range rmpState.AllSenecaProfs {
		if prof.Node.ID == selectedProfID {
			rmpState.SelectedProf = prof
		}
	}
	err := SendProfInformation(sess, i, rmpState.SelectedProf)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Prof %s information sent", rmpState.SelectedProf.fullName()), nil
}

// createSelectMenu create a select menu containing the profs.
// profs  : list of profs to display in the select menu.
// disable: if true, the select menu will be disabled.
// return: a select menu component.
func createSelectMenu(profs []ProfNode, disable bool) discordgo.SelectMenu {
	MinValues := 1
	menu := discordgo.SelectMenu{
		CustomID:    profSelectMenuID,
		Placeholder: "Please select a prof",
		MinValues:   &MinValues,
		MaxValues:   1,
		Options:     []discordgo.SelectMenuOption{},
		Disabled:    disable,
	}

	// add all profs as an option to the select menu
	for _, prof := range profs {
		option := discordgo.SelectMenuOption{
			Label:       prof.fullName(),
			Value:       prof.Node.ID,
			Description: fmt.Sprintf("Department: %s", prof.Node.Department),
			Emoji:       discordgo.ComponentEmoji{},
			Default:     false,
		}
		menu.Options = append(menu.Options, option)
	}
	return menu
}

// SendProfInformation reply to an interaction with information about a professor.
// sess  : discord session.
// i     : discord interaction.
// prof  : professor information to send.
// return: error if any.
func SendProfInformation(sess *discordgo.Session, i *discordgo.InteractionCreate, prof ProfNode) error {
	return sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					URL:         prof.rmpURL(),
					Type:        "",
					Title:       fmt.Sprintf("%s %s", prof.Node.FirstName, prof.Node.LastName),
					Description: prof.profDescription(),
					Timestamp:   "",
					Color:       0,
					Footer: &discordgo.MessageEmbedFooter{
						Text:         "Information retrieved from ratemyprof.com",
						IconURL:      "https://pbs.twimg.com/profile_images/1146077191043788800/hG1lAGm9_400x400.png",
						ProxyIconURL: "",
					},
					Author: &discordgo.MessageEmbedAuthor{
						URL:          "https://www.youtube.com/watch?v=dQw4w9WgXcQ&ab_channel=RickAstley",
						Name:         fmt.Sprintf("brought to you by @%s's mom TM", i.Member.User.Username),
						IconURL:      "",
						ProxyIconURL: "",
					},
				},
			},
		},
	})

}

// autoCompleteHandler handles filling the auto complete results
// sess  : discord session.
// i     : discord interaction.
func autoCompleteHandler(sess *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := utils.ParseUserOptions(sess, i)
	profName := options["profname"].StringValue()
	res := SearchRmpProfByName(profName)

	choices := []*discordgo.ApplicationCommandOptionChoice{}

	// add the choice the user is currently typing list as well
	if profName != "" {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  profName,
			Value: profName,
		})
	}

	for _, prof := range res.Data.Search.Teachers.Edges {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  prof.fullName(),
			Value: prof.fullName(),
		})
	}

	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func init() {
	command.Register(rmp)
}
