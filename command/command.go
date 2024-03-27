package command

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Contains all slash commands for this bot
var CmdStore []Command

// A handler function type for slash commands
type HandlerFunc func(sess *discordgo.Session, i *discordgo.InteractionCreate) (successMsg string, err error)

type Command struct {
	// Name of the slash command
	Name string
	// Command definition
	ApplicationCommand func() *discordgo.ApplicationCommand
	// Handler for handling slash command interactions. This function should edit an interaction response. Returns a log message, and error
	EditInteractionResponse HandlerFunc
	// Handler for handling slash command interactions. This function should send a direct interaction response. Returns a log message, and error
	HandlerFunc HandlerFunc
	// Handler for auto completion requests
	InteractionApplicationCommandAutocomplete HandlerFunc
	// Command components, with their handler
	Components []struct {
		// Custom component ID
		ComponentID string
		// Handler for component interactions
		ComponentHandler HandlerFunc
	}
}

// Register registers a new slash command.
func Register(c Command) {
	CmdStore = append(CmdStore, c)
}

// // GetCmdDefs return all slash command definitions.
func GetCmdDefs() []*discordgo.ApplicationCommand {
	slashCmds := make([]*discordgo.ApplicationCommand, 0)
	for _, cmd := range CmdStore {
		slashCmds = append(slashCmds, cmd.ApplicationCommand())
	}
	return slashCmds
}

// GetCmdHandler returns a map of command name to all their handler function.
func GetCmdHandler() map[string]Command {
	cmdHandlers := map[string]Command{}
	for _, cmd := range CmdStore {
		cmdHandlers[cmd.ApplicationCommand().Name] = cmd
		// if cmd.EditInteractionResponse != nil {
		//    fmt.Printf("command %s, using edit handler\n", cmd.ApplicationCommand().Name)
		//    cmdHandlers[cmd.ApplicationCommand().Name] = cmd.EditInteractionResponse
		// } else if cmd.HandlerFunc != nil {
		//    fmt.Printf("command %s, using handler func\n", cmd.ApplicationCommand().Name)
		//    cmdHandlers[cmd.ApplicationCommand().Name] = cmd.HandlerFunc
		// } else {
		//    panic("A command must have either `EditInteractionResponse` or `HandlerFunc` defined")
		// }
	}
	return cmdHandlers

}

// GetComponentHandler returns a map of component ID and the handler function.
func GetComponentHandler() map[string]HandlerFunc {
	componentHandlers := map[string]HandlerFunc{}
	for _, cmd := range CmdStore {
		for _, component := range cmd.Components {
			componentHandlers[component.ComponentID] = component.ComponentHandler
		}
	}
	return componentHandlers
}

// RemoveCommands removes all registered slash commands from all servers the bot is in
func RemoveCommands(sess *discordgo.Session, registeredCommands []*discordgo.ApplicationCommand) {
	log.Println("Removing commands...")
	for _, gld := range sess.State.Guilds {
		for _, cmd := range registeredCommands {
			err := sess.ApplicationCommandDelete(sess.State.User.ID, gld.ID, cmd.ID)
			if err != nil {
				log.Printf("Cannot delete '%v' command in guild '%v': %v\n", cmd.Name, gld.Name, err)
			} else {
				log.Printf("Removing command /%v in guild %v", cmd.Name, gld.Name)
			}
		}
	}
}

// RegisterCommands register an array of commands to a discord session. Returns a list of all the registered commands.
//
// this function will panic if registration of a command fails.
func RegisterCommands(sess *discordgo.Session) (registeredCommands []*discordgo.ApplicationCommand) {
	registeredCommands = make([]*discordgo.ApplicationCommand, 0)
	// c := LoadConfig()
	// ignoreGuilds := make([]discordgo.Guild, 0)

	for _, gld := range sess.State.Guilds {
		// NOTE: should we do this?
		// roles := gld.Roles
		// for _, role := range roles {
		// // if bot role is not at the top 3, dont register commmands here
		// if role.Name == "debate_dragon" && role.Position > 4 {
		// log.Printf("Bot role is not at the top 3. Bot position: %d, not registering commands in guild %v", role.Position, gld.Name)
		// ignoreGuilds = append(ignoreGuilds, *gld)
		// }
		// }

		for _, v := range CmdStore {
			cmd, err := sess.ApplicationCommandCreate(sess.State.User.ID, gld.ID, v.ApplicationCommand())
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.ApplicationCommand().Name, err)
			}
			fmt.Printf("Registering /%s in guild %v\n", cmd.Name, gld.Name) // __AUTO_GENERATED_PRINT_VAR__
			registeredCommands = append(registeredCommands, cmd)
		}
	}
	return registeredCommands
}
