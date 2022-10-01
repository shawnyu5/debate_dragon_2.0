package commands

import "github.com/bwmarrin/discordgo"

type CommandStruct struct {
	// name of the slash command, as will be used in discord
	Name string
	// command object
	Obj func() *discordgo.ApplicationCommand
	// command handler
	CommandHandler func(sess *discordgo.Session, i *discordgo.InteractionCreate)
	// component handler
	ComponentHandler func(sess *discordgo.Session, i *discordgo.InteractionCreate)
}