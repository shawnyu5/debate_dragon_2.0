package commands

import (
	"github.com/bwmarrin/discordgo"
)

// a handler function type for slash command
type HandlerFunc func(sess *discordgo.Session, i *discordgo.InteractionCreate)

// defines the interface for a slash command component
type CommandInter interface {
	Obj() *discordgo.ApplicationCommand
	// command handler to handle the slash command
	Handler(sess *discordgo.Session, i *discordgo.InteractionCreate)
}

type CommandStruct struct {
	// name of the slash command, as will be used in discord
	Name string
	// command object
	Obj func() *discordgo.ApplicationCommand
	// command handler to handle the slash command
	Handler func(sess *discordgo.Session, i *discordgo.InteractionCreate)
	// each component can have multiple components, each with their own handler and ID
	Components []struct {
		// component custom ID
		ComponentID string
		// component handler for button clicking and such
		ComponentHandler HandlerFunc
	}
}

type ComponentHandler struct {
	// component custom ID
	ComponentID string
	// component handler for button clicking and such
	ComponentHandler HandlerFunc
}
