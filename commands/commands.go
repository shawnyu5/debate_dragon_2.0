package commands

import (
	"github.com/bwmarrin/discordgo"
)

// a handler function type for slash command
type HandlerFunc func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error)

type CommandStruct struct {
	// name of the slash command, as will be used in discord
	Name string
	// command object
	Obj func() *discordgo.ApplicationCommand
	// command handler to handle the slash command
	Handler HandlerFunc
	// each component can have multiple components, each with their own handler and ID
	Components []struct {
		// component custom ID
		ComponentID string
		// component handler for button clicking and such
		ComponentHandler HandlerFunc
	}
}

// A single discord component
type Component struct {
	ComponentID      string
	ComponentHandler HandlerFunc
}

// A discord slash command
type Command interface {
	Def() *discordgo.ApplicationCommand
	Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error)
	Components() []Component
}
