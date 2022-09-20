package commands

import "github.com/bwmarrin/discordgo"

type CommandStruct struct {
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
	Obj     func() *discordgo.ApplicationCommand
	Name    string
}

type Command interface {
	// command handler function
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
	// return a command object
	Obj() *discordgo.ApplicationCommandOption
}
