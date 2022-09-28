package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

	Ivan struct {
		Emotes []struct {
			Name         string `json:"name"`
			FileLocation string `json:"fileLocation"`
		} `json:"emotes"`
	} `json:"ivan"`
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
			err := sess.ApplicationCommandDelete(cmd.ApplicationID, cmd.ID, gld.ID)
			// err := sess.ApplicationCommandDelete(cmd.ApplicationID, cmd.ID, "")
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

// DeferReply defers a reply
func DeferReply(sess *discordgo.Session, i *discordgo.Interaction) error {
	err := sess.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	return err
}

// LoadConfig loads the config file, and return the config in a struct
func LoadConfig() Config {
	var c Config
	// read json file
	f, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, &c)
	return c
}
