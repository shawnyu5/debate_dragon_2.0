package utils

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Token       string `json:"token"`
	TokenDev    string `json:"token_dev"`
	AppID       string `json:"app_id"`
	AppIDDev    string `json:"app_id_dev"`
	ClientID    string `json:"clientID"`
	GuildID     string `json:"guildID"`
	LogLevel    string `json:"logLevel"`
	Development bool   `json:"development"`

	CarmenRambles struct {
		CarmenID          string `json:"carmenId"`
		ChannelID         string `json:"channelId"`
		CoolDown          int64  `json:"coolDown"`
		GuildID           string `json:"guildID"`
		MessageLimit      int64  `json:"messageLimit"`
		SubscribersRoleID string `json:"subscribersRoleID"`
	} `json:"carmenRambles"`
}

// RegisterCommands register an array of commands to a discord session.
// Receives an instance of discord session to store commands in. An array of discord application commands to keep track of the stored commands. And an array of commands to register
// Will panic if registration of a command fails.
func RegisterCommands(dg *discordgo.Session, commands []*discordgo.ApplicationCommand, registeredCommands []*discordgo.ApplicationCommand) {
	log.Println("Adding commands...")
	for _, gld := range dg.State.Guilds {
		for i, v := range commands {
			cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, gld.ID, v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
			fmt.Println(fmt.Sprintf("Registering /%s in guild %v", cmd.Name, gld.Name)) // __AUTO_GENERATED_PRINT_VAR__
			registeredCommands[i] = cmd
		}
	}
}

// RemoveCommands will delete all registered commands in all servers the discord bot is currently in
// receives an instance of discord session to remove commands from. An array of discord application commands to remove
func RemoveCommands(sess *discordgo.Session, registeredCommands []*discordgo.ApplicationCommand, c Config) {
	log.Println("Removing commands...")
	for _, gld := range sess.State.Guilds {
		for _, cmd := range registeredCommands {
			// TODO: fix this
			err := sess.ApplicationCommandDelete(sess.State.User.ID, cmd.ID, gld.ID)
			// err := sess.ApplicationCommandDelete(sess.State.User.ID, gld.ID, cmd.ID)
			if err != nil {
				log.Printf("Cannot delete '%v' command in guild '%v': %v\n", cmd.Name, gld.Name, err)
			} else {
				log.Printf("Removing command /%v in guild %v", cmd.Name, gld.Name)
			}
		}
	}
}

// ParseUserOptions parses the user option passed to a command, and returns a map of data options
func ParseUserOptions(sess *discordgo.Session, i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}
