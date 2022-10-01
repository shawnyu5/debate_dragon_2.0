package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

type CommandStruct struct {
	// name of the slash command, as will be used in discord
	Name string
	// command object
	Obj func() *discordgo.ApplicationCommand
	// command handler to handle the slash command
	CommandHandler func(sess *discordgo.Session, i *discordgo.InteractionCreate)
	// each component can have multiple components, each with their own handler and ID
	Components []struct {
		// component custom ID
		ComponentID string
		// component handler for button clicking and such
		ComponentHandler utils.HandlerFunc
	}
}
