package commands

import "github.com/bwmarrin/discordgo"

type CommandStruct struct {
	// command object
	Obj func() *discordgo.ApplicationCommand
	// command handler
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// type Command interface {
// // command handler function
// Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
// // return a command object
// Obj() *discordgo.ApplicationCommandOption
// }
