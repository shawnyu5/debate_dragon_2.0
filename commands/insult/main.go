package insult

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

func Obj() *discordgo.ApplicationCommand {
	obj := &discordgo.ApplicationCommand{
		Name:        "insult",
		Description: "Ping someone to deliver a gut wrenching insult",
		Options: []*discordgo.ApplicationCommandOption{
			// get the user to insult
			{
				Name:        "user",
				Description: "user to insult",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
			// have the option to send insult anonymously
			{
				Name:        "anonymous",
				Description: "whether or not to send the insult anonymously",
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Required:    false,
			},
		},
	}
	return obj
}

// Handler a handler function for insult command
func Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	optionsMap := utils.ParseUserOptions(sess, i)
	user := optionsMap["user"].UserValue(sess)
	if user.ID == "652511543845453855" {
		user = i.Message.Author
	}
	insult := getInsult(optionsMap["user"].UserValue(sess))
	// 652511543845453855

	// send a normal insult
	if optionsMap["anonymous"] == nil {
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: insult,
			},
		})
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		// send an anonymous insult
		_, err := sess.ChannelMessageSend(i.ChannelID, insult)
		if err != nil {
			log.Fatalln(err)
		}

		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Your insult has been send >:)",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}

// getInsult return an insult ping the user passed into the function
func getInsult(user *discordgo.User) string {
	// make http get request to insult api
	resp, err := http.Get("https://insult.mattbas.org/api/insult")
	if err != nil {
		log.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("<@%v> %v", user.ID, string(body))
}
