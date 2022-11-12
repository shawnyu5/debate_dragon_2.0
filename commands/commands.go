package commands

import (
	"github.com/bwmarrin/discordgo"
)

// a handler function type for slash command
type HandlerFunc func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error)

// defines the interface for a slash command component
type CommandInter interface {
	// return the command handler to handle the slash command
	GetHandler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error)
	// return the name of the command
	GetName() string
}

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

// GetHandler slash command handler
func (c CommandStruct) GetHandler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	return c.Handler(sess, i)
}

// GetName return the name of the command
func (c CommandStruct) GetName() string {
	return c.Name
}

type ComponentHandler struct {
	// component custom ID
	ComponentID string
	// component handler for button clicking and such
	ComponentHandler HandlerFunc
}
