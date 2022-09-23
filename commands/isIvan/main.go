package isIvan

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
	obj := &discordgo.ApplicationCommand{
		Version:     "1.0",
		Name:        "isivan",
		Description: "Use machine learning to predict if a user or a message is Ivan",
		Options: []*discordgo.ApplicationCommandOption{
			// {
			// Name:         "user",
			// Description:  "the user to check",
			// Type:         discordgo.ApplicationCommandOptionUser,
			// Required:     false,
			// Autocomplete: false,
			// },
			{
				Name:        "message",
				Description: "A message to check",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	}
	return obj
}

func handler(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	optionsMap := utils.ParseUserOptions(sess, i)
	fmt.Println(fmt.Sprintf("handler optionsMap: %v", optionsMap)) // __AUTO_GENERATED_PRINT_VAR__

	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "HELLO???",
		},
	})

	res := "HELLO????"
	mess, err := sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:         &res,
		Components:      &[]discordgo.MessageComponent{},
		Embeds:          &[]*discordgo.MessageEmbed{},
		Files:           []*discordgo.File{},
		AllowedMentions: &discordgo.MessageAllowedMentions{},
	})
	if err != nil {
		log.Println(err)
	}
	fmt.Println(fmt.Sprintf("handler mess: %v", mess)) // __AUTO_GENERATED_PRINT_VAR__
	// err = sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// Type: discordgo.InteractionResponseChannelMessageWithSource,
	// Data: &discordgo.InteractionResponseData{
	// Content:         "HELLO???",
	// Components:      []discordgo.MessageComponent{},
	// Embeds:          []*discordgo.MessageEmbed{},
	// AllowedMentions: &discordgo.MessageAllowedMentions{},
	// Files:           []*discordgo.File{},
	// Title:           "FOLLOW UP",
	// },
	// })
	if err != nil {
		log.Println(err)
	}

	// gld, err := sess.Guild(i.GuildID)
	// if err != nil {
	// panic(err)
	// }

	// c := gld.Channels
	// for _, v := range c {
	// if v.Type == discordgo.ChannelTypeGuildText {
	// res, err := sess.Client.Get(v.ID)
	// if err != nil {
	// panic(err)
	// }
	// fmt.Println(fmt.Sprintf("handler res: %v", res)) // __AUTO_GENERATED_PRINT_VAR__
	// }
	// }
}
