package command

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

// Contains all slash commands for this bot
var CmdStore []Command

// A handler function type for slash commands
type HandlerFunc func(sess *discordgo.Session, i *discordgo.InteractionCreate) (successMsg string, err error)

type Command struct {
	// // Name of the slash command
	// Name string
	// Command definition
	ApplicationCommand func() *discordgo.ApplicationCommand
	// Handler for handling slash command interactions. This function should edit an interaction response. Returns a log message, and error if any
	EditInteractionResponse HandlerFunc
	// Handler for handling slash command interactions. This function should send a direct interaction response. Returns a log message, and error if any
	InteractionRespond HandlerFunc
	// Handler for auto completion requests. Returns a log message, and error if any
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
// func RemoveCommands(sess *discordgo.Session, registeredCommands []*discordgo.ApplicationCommand) {
// 	log.Info("Removing commands...")
// 	for _, gld := range sess.State.Guilds {
// 		for _, cmd := range registeredCommands {
// 			err := sess.ApplicationCommandDelete(sess.State.User.ID, gld.ID, cmd.ID)
// 			if err != nil {
// 				log.Printf("Cannot delete '%v' command in guild '%v': %v\n", cmd.Name, gld.Name, err)
// 			} else {
// 				log.Printf("Removing command /%v in guild %v", cmd.Name, gld.Name)
// 			}
// 		}
// 	}
// }

// RegisterCommands register an array of commands to a discord session. Returns a list of all the registered commands.
//
// this function will panic if registration of a command fails.
func RegisterCommands(sess *discordgo.Session) (registeredCommands []*discordgo.ApplicationCommand) {
	registeredCommands = make([]*discordgo.ApplicationCommand, 0)

	for _, gld := range sess.State.Guilds {
		for _, c := range CmdStore {
			cmd, err := sess.ApplicationCommandCreate(sess.State.User.ID, gld.ID, c.ApplicationCommand())
			if err != nil {
				log.Fatalf("Cannot create '%v' command: %v", c.ApplicationCommand().Name, err)
			}
			log.Infof("Registering /%s in guild %v", cmd.Name, gld.Name) // __AUTO_GENERATED_PRINT_VAR__
			registeredCommands = append(registeredCommands, cmd)
		}
	}

	return registeredCommands
}

// RemoveAllSlashCommandFromAllGuilds removes all slash commands from all guilds this bot is in
func RemoveAllSlashCommandFromAllGuilds(dg *discordgo.Session) error {
	glds := dg.State.Guilds
	var wg sync.WaitGroup

	for _, guild := range glds {
		log.Infof("Removing slash commands from guild %s", guild.Name)
		log.Debug("Getting slash commands")
		cmds, err := dg.ApplicationCommands(dg.State.User.ID, guild.ID)
		if err != nil {
			log.Fatalf("failed to get slash commands: %s", err)
		}

		log.Debugf("Got %d slash commands in guild %s", len(cmds), guild.Name)

		for _, cmd := range cmds {
			wg.Add(1)
			go func(c *discordgo.ApplicationCommand, guildName string, guildID string) {
				defer wg.Done()

				log.Infof("Deleting slash command %s from guild %s", c.Name, guildName)
				err := dg.ApplicationCommandDelete(dg.State.Application.ID, guildID, c.ID)
				if err != nil {
					log.Errorf("Failed to delete slash command: %s", err)
				}
			}(cmd, guild.Name, guild.ID)
		}
	}

	wg.Wait()
	log.Info("Finished deleting all slash commands")

	return nil
}
