package commands

import "github.com/bwmarrin/discordgo"

// type Command interface {
// // command handler function
// Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
// // return a command object
// Obj() CommandObj
// }

type Command interface {
	// command handler function
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
	// return a command object
	Obj() *discordgo.ApplicationCommandOption
}
