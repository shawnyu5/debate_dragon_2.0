package compliment

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

type Compliment struct {
	PickupLine string `json:"random_cheesy_pickup_line"`
}

var compliment = command.Command{
	ApplicationCommand: func() *discordgo.ApplicationCommand {
		return &discordgo.ApplicationCommand{
			Name:        "caramel-bot-compliment",
			Description: "Give another user a compliment",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user",
					Description: "The user you would like to compliment",
					Type:        discordgo.ApplicationCommandOptionUser,
					Required:    true,
				},
			},
		}
	},
	HandlerFunc: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		userOptions := utils.ParseUserOptions(sess, i)
		user := userOptions["user"].UserValue(sess)

		// // Compliment struct to hold compliment information
		// var compliment Compliment

		// // Get options from the application data
		// options := i.ApplicationCommandData().Options
		// var message = ""

		// // If the option exists, add the result of the user option to the message as a ping
		// if options[0] != nil {
		//    message = fmt.Sprintf("<@%s> %s", user.ID, compliment.Compliment)
		// }
		compliment, err := getCompliment()
		if err != nil {
			return "", err
		}

		// Build out the interaction response
		err = sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("<@%s> %s", user.ID, compliment),
			},
		})

		if err != nil {
			return "", err
		}
		return fmt.Sprintf("User %s complimented", user.Username), nil
	},
}

func getCompliment() (string, error) {
	url := "https://nerdy-pickup-lines1.p.rapidapi.com/cheesy-pickup-lines/random"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", utils.LoadConfig().RapidAPIKey)
	req.Header.Add("X-RapidAPI-Host", "nerdy-pickup-lines1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	var compliment Compliment
	err := json.Unmarshal(body, &compliment)
	if err != nil {
		return "", err
	}

	return compliment.PickupLine, nil
}
func init() {
	command.Register(compliment)
}
