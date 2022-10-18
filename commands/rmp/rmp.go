package rmp

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

var CommandObj = commands.CommandStruct{
	Name:           "rmp",
	Obj:            obj,
	CommandHandler: handler,
	// Components: []struct {
	// ComponentID      string
	// ComponentHandler commands.HandlerFunc
	// }{},
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
	j, _ := json.MarshalIndent(searchResult, "", "   ")
	fmt.Printf("handler prof: %+v\n", string(j)) // __AUTO_GENERATED_PRINT_VAR__

	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "hello",
			Flags:   0,
		},
	})
	if err != nil {
		panic(err)
	}
}
