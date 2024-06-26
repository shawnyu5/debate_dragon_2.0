package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/spf13/afero"
)

var AppFs = afero.NewOsFs()

// config object as defined in config.json.
// do not json martial sensitive fields such as discord token
type Config struct {
	Token          string `json:"-"`
	TokenDev       string `json:"-"`
	RedditUserName string `json:"-"`
	RedditClientId string `json:"-"`
	RedditSecret   string `json:"-"`
	RedditPassword string `json:"-"`
	LogLevel       string `json:"logLevel"`
	Development    bool   `json:"-"`
	RapidAPIKey    string `json:"-"`
	// ID of the bot owner
	BotOwner string `json:"botOwner"`
	// path to the local db
	DbPath string `json:"dbPath"`
	Emotes []struct {
		// name of emote
		Name string `json:"name"`
		// url to emote
		URL string `json:"url"`
	} `json:"emotes"`
	// config for new member greetings
	NewMemberGreeting struct {
		Config []struct {
			ServerName string `json:"serverName"`
			RoleID     string `json:"roleID"`
			ServerID   string `json:"serverID"`
			ChannelID  string `json:"channelID"`
			Enable     bool   `json:"enable"`
		} `json:"config"`
	} `json:"newMemberGreeting"`
	Ivan struct {
		Emotes []struct {
			Name         string `json:"name"`
			FileLocation string `json:"fileLocation"`
		} `json:"emotes"`
	} `json:"ivan"`
	SubForCarmen struct {
		// toggle this feature on and off
		On bool `json:"on"`
		// id of carmen user to track messages of
		CarmenID string `json:"carmenId"`
		// cool down, defined in minutes
		CoolDown int `json:"coolDown"`
		// the guild to keep track of carmen messages
		GuildID string `json:"guildID"`
		// number of messages before a notification is triggered
		MessageLimit      int    `json:"messageLimit"`
		SubscribersRoleID string `json:"subscribersRoleID"`
		// channels to ignore
		IgnoredChannels []string `json:"ignoredChannels"`
	} `json:"subForCarmen"`
}

// LoadConfig loads config.json and .env.
func LoadConfig() Config {
	var c Config
	// read json file
	f, err := AppFs.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, &c)

	godotenv.Load()
	c.Token = os.Getenv("TOKEN")
	c.RapidAPIKey = os.Getenv("RAPIDAPI_KEY")
	c.TokenDev = os.Getenv("TOKEN_DEV")
	c.RedditUserName = os.Getenv("REDDIT_USERNAME")
	c.RedditClientId = os.Getenv("REDDIT_CLIENT_ID")
	c.RedditSecret = os.Getenv("REDDIT_SECRET")
	c.RedditPassword = os.Getenv("REDDIT_PASSWORD")

	dev := os.Getenv("DEVELOPMENT")
	if dev == "true" {
		c.Development = true
	} else {
		c.Development = false
	}
	return c
}

// RegisterCommands register an array of commands to a discord session.
// sess: discord session.
// commands: array of discord commands to register.
// registeredCommands: array of commands to keep track of registered commands.
// this function will panic if registration of a command fails.
func RegisterCommands(sess *discordgo.Session, commands []*discordgo.ApplicationCommand, registeredCommands []*discordgo.ApplicationCommand) {
	// c := LoadConfig()
	ignoreGuilds := make([]discordgo.Guild, 0)
	log.Info("Adding commands...")
	for _, gld := range sess.State.Guilds {
		// roles := gld.Roles
		// for _, role := range roles {
		// // if bot role is not at the top 3, dont register commmands here
		// if role.Name == "debate_dragon" && role.Position > 4 {
		// log.Printf("Bot role is not at the top 3. Bot position: %d, not registering commands in guild %v", role.Position, gld.Name)
		// ignoreGuilds = append(ignoreGuilds, *gld)
		// }

		// }
		for i, v := range commands {
			if ignore := Contains(ignoreGuilds, gld.ID); ignore {
				log.Printf("Ignoring guild %v", gld.Name)
				continue
			}
			cmd, err := sess.ApplicationCommandCreate(sess.State.User.ID, gld.ID, v)
			if err != nil {
				log.Info("Cannot create '%v' command: %v", v.Name, err)
			}
			fmt.Printf("Registering /%s in guild %v\n", cmd.Name, gld.Name) // __AUTO_GENERATED_PRINT_VAR__
			registeredCommands[i] = cmd
		}
	}
}

// RemoveCommands will delete all registered commands in all servers the discord bot is currently in.
// sess: discord session.
// registeredCommands: array of commands to remove.
func RemoveCommands(sess *discordgo.Session, registeredCommands []*discordgo.ApplicationCommand) {
	log.Info("Removing commands...")
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

// ParseUserOptions parses the user option passed to a command.
// sess: discord session.
// return: a map of input name : input value
func ParseUserOptions(sess *discordgo.Session, i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

// DeferReply defers a reply.
func DeferReply(sess *discordgo.Session, i *discordgo.Interaction) error {
	err := sess.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	return err
}

// SendErrorMessage send an empheral interaction response message notifying the user something went wrong with the slash command. With an optional error message.
func SendErrorMessage(sess *discordgo.Session, i *discordgo.InteractionCreate, err string) {
	e := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Something went wrong... " + err,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if e != nil {
		log.Printf("Error sending error response: %v", e)
	}
}

// EditErrorMessage edits the previous interaction response with an error message notifying the user something went wrong with the slash command. With an optional error message.
//
// Since the error message is not empheral, it will be deleted after 5 seconds
func EditErrorMessage(sess *discordgo.Session, i *discordgo.InteractionCreate, err string) {
	content := "Something went wrong... " + err
	msg, e := sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	},
	)
	if e != nil {
		log.Printf("Error editing error response: %v", e)
	}

	time.AfterFunc(5*time.Second, func() {
		err := sess.ChannelMessageDelete(msg.ChannelID, msg.ID)
		if err != nil {
			log.Warnf("Failed to delete command failure message: %s", err)
		}
	})
}

// DeleteAllMessages Delete all messages in a channel.
// sess    : the discord session.
// i       : discord interaction.
// messages: array of discord messages to delete.
func DeleteAllMessages(sess *discordgo.Session, i *discordgo.InteractionCreate, messages []*discordgo.Message) {
	for _, message := range messages {
		go func(mess *discordgo.Message) {
			err := sess.ChannelMessageDelete(i.ChannelID, mess.ID)
			if err != nil {
				log.Info(err)
				return
			}

		}(message)
	}

}

// Contains checks if an array contains an element.
// arr : the array to check.
// elem: the element to check for.
// returns true if the array contains the element, false otherwise.
func Contains(arr []discordgo.Guild, id string) bool {
	for _, v := range arr {
		if v.ID == id {
			return true
		}
	}
	return false
}

// GetMembersWithRole get all members from a guild with a specif role.
// guildID: the guild ID to fetch members from.
// roleID : the role ID to fetch.
// return: an array of members with the role, and any errors that arise
func GetMembersWithRole(sess *discordgo.Session, guildID, roleID string) ([]*discordgo.Member, error) {
	// checks if an array of members contains a member.
	contains := func(arr []string, id string) bool {
		for _, v := range arr {
			if v == id {
				return true
			}
		}
		return false
	}
	// all members in the guild
	membersCollection := make([]*discordgo.Member, 0)
	for {
		members := make([]*discordgo.Member, 0)
		var err error
		// the first request should get as many members as possible
		if len(membersCollection) == 0 {
			members, err = sess.GuildMembers(guildID, "", 1000)
			if err != nil {
				return nil, err
			}
		} else {
			// get members after the last member in the collection
			members, err = sess.GuildMembers(guildID, membersCollection[len(membersCollection)-1].User.ID, 1000)
			if err != nil {
				return nil, err
			}
		}
		// keep requesting members from the API till we get nothing back
		if len(members) == 0 {
			break
		}
		membersCollection = append(membersCollection, members...)
	}

	membersWithRole := make([]*discordgo.Member, 0)
	for _, member := range membersCollection {
		// if a member does have the role we are looking for, then keep track of it
		if contains(member.Roles, roleID) {
			membersWithRole = append(membersWithRole, member)
		}
	}
	return membersWithRole, nil
}

// ShrinkFontSize shrink the font size passed in based on the length of user input and the max length of character.
// fontSize: the initial font size.
// maxCharacterSize: the max length of character.
// Returns the new font size.
func ShrinkFontSize(fontSize int, userInput string, maxCharacterSize int) int {
	// 7 is the max character at current size
	if len(userInput) > maxCharacterSize {
		return ShrinkFontSize(fontSize-5, userInput, maxCharacterSize+5)
	}
	return fontSize
}
