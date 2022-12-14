package newmember

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

const greeterLabel = "greeters"

// map of guild id to list of greeters discord users
type GuildGreeters map[string][]discordgo.User

// type Greeter struct {
// GuildID string
// User    []discordgo.User
// }

// greet new members
type NewMember struct{}

// Components implements commands.Command
func (NewMember) Components() []commands.Component {
	return nil
}

// Def implements commands.Command
func (NewMember) Def() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Version:     "1.0.0",
		Name:        "newmembergreeting",
		Description: "Opt in to greet new members",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "opt_in",
				Description: "Whether you would like to greet the new comers to the server",
				Required:    true,
			},
		},
	}
}

// Handler implements commands.Command
func (NewMember) Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	db, err := openDB()
	if err != nil {
		return "", err
	}

	defer db.Close()

	userOptions := utils.ParseUserOptions(sess, i)
	// add user to db
	if userOptions["opt_in"].BoolValue() {
		guildGreeters, err := GetGreeters(db, i.GuildID)
		if err != nil {
			return "", err
		}

		AddGreeterToGuild(guildGreeters, i.Member.User, i.GuildID)
		err = SaveGreeters(db, guildGreeters)
		if err != nil {
			return "", err
		}

		// get current guild information
		guild, err := sess.Guild(i.GuildID)
		if err != nil {
			return "", err
		}

		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Type:        discordgo.EmbedTypeArticle,
						Title:       "Congrats!!!",
						Description: fmt.Sprintf("Congratulations, you are now a \"server owner\", in the great democracy of **%s**", guild.Name),
						Timestamp:   "",
						Color:       0,
						Footer:      &discordgo.MessageEmbedFooter{},
						Image:       &discordgo.MessageEmbedImage{},
						Thumbnail:   &discordgo.MessageEmbedThumbnail{},
						Video:       &discordgo.MessageEmbedVideo{},
						Provider:    &discordgo.MessageEmbedProvider{},
						Author:      &discordgo.MessageEmbedAuthor{},
					},
				},
			},
		})
		return fmt.Sprintf("%s saved as a new greeter", i.Member.User.Username), err
	} else {
		guildGreeters, err := GetGreeters(db, i.GuildID)
		if err != nil {
			return "", err
		}
		guildGreeters[i.GuildID] = remove(guildGreeters[i.GuildID], i.Member.User.ID)

		err = SaveGreeters(db, guildGreeters)
		if err != nil {
			return "", err
		}

		roleID, err := getGuildRole(i.GuildID)
		if err != nil {
			return "", err
		}
		err = sess.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, roleID)
		if err != nil {
			return "", err
		}

		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Type:        discordgo.EmbedTypeArticle,
						Title:       "Sorry to see you go...",
						Description: "Your role as a \"server owner\" has been removed...",
						Timestamp:   "",
						Color:       0,
						Footer:      &discordgo.MessageEmbedFooter{},
						Image:       &discordgo.MessageEmbedImage{},
						Thumbnail:   &discordgo.MessageEmbedThumbnail{},
						Video:       &discordgo.MessageEmbedVideo{},
						Provider:    &discordgo.MessageEmbedProvider{},
						Author:      &discordgo.MessageEmbedAuthor{},
					},
				},
			},
		})

		return fmt.Sprintf("%s removed as a greeter", i.Member.User.Username), err
	}
}

// Greet greets a new members to a discord server.
// sess: the discord session.
// user: the user that joined the server.
// return: string, and error for logging
func Greet(sess *discordgo.Session, user *discordgo.GuildMemberAdd) (string, error) {
	c := utils.LoadConfig()
	// check if greeting is enabled for this server
	for _, guild := range c.NewMemberGreeting.Config {
		if guild.ServerID == user.GuildID && !guild.Enable {
			return "greeting not enabled for this server...", nil
		}
	}

	db, err := openDB()
	if err != nil {
		return "", err
	}

	greeters, err := GetGreeters(db, user.GuildID)
	if err != nil {
		return "", err
	}

	guildGreeters := greeters[user.GuildID]
	if len(guildGreeters) == 0 {
		return "", errors.New("no greeters found for this server")
	}

	// generate random number
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := len(guildGreeters) - 1
	randomInt := (rand.Intn(max-min+1) + min)

	channelID := ""
	for _, server := range c.NewMemberGreeting.Config {
		if server.ServerID == user.GuildID {
			channelID = server.ChannelID
		}
	}

	// add roleID to user
	roleID, err := getGuildRole(user.GuildID)
	if err != nil {
		return "", err
	}

	err = removeAllRole(sess, user.GuildID, roleID, guildGreeters)
	if err != nil {
		return "", err
	}

	err = sess.GuildMemberRoleAdd(user.GuildID, guildGreeters[randomInt].ID, roleID)
	if err != nil {
		return "", err
	}
	_, err = sess.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Hi <@%s>, welcome to <@%s>'s server!", user.User.ID, guildGreeters[randomInt].ID),
	})
	if err != nil {
		return "", err
	}

	return "Greeting sent", nil
}

// openDB opens the database.
// return: an instance of the database.
func openDB() (*badger.DB, error) {
	c := utils.LoadConfig()
	// openDB opens a connection to the local database
	opts := badger.DefaultOptions(c.DbPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return db, err

}

// GetGreeters gets the greeters for a guild.
// db: the database to get from.
// guildID: the guild to get the greeters for.
// return: the greeters for the guild. And errors if any
func GetGreeters(db *badger.DB, guildID string) (GuildGreeters, error) {
	greeters := make(GuildGreeters, 0)
	err := db.View(func(txn *badger.Txn) error {
		item, _ := txn.Get([]byte(greeterLabel))

		if item == nil {
			return nil
		}
		err := item.Value(func(val []byte) error {
			err := json.Unmarshal(val, &greeters)
			return err
		})
		return err
	})
	return greeters, err
}

// SaveGreeters saves the greeters for a guild.
// db: the database to save to.
// data: the greeters to save.
// return: an error if any.
func SaveGreeters(db *badger.DB, data GuildGreeters) error {
	err := db.Update(func(txn *badger.Txn) error {
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}

		err = txn.Set([]byte(greeterLabel), b)
		return err
	})

	return err
}

// remove removes a user from an array of users based on the ID
// users: the array of users to remove from.
// id: the id of the user to remove.
// return: the array of users with the user removed.
func remove(users []discordgo.User, id string) []discordgo.User {
	// handle when array is empty
	if len(users) == 0 {
		return users
	}
	for i, user := range users {
		if user.ID == id {
			users[i] = users[len(users)-1]
		}
	}
	return users[:len(users)-1]
}

// AddGreeterToGuild adds a user to the greeter list for their array. Makes sure there are no duplicate users per guild
// guildGreeters: the greeters for the guild.
// user: the user to add.
// guildID: the guild id to add the user to.
func AddGreeterToGuild(guildGreeters GuildGreeters, user *discordgo.User, guildID string) {
	// check if the user is already in the array
	for _, greeter := range guildGreeters[guildID] {
		if greeter.ID == user.ID {
			return
		}
	}
	guildGreeters[guildID] = append(guildGreeters[guildID], *user)
}

// getGuildRole gets server owner role for a guild
// guildID: the guild to retrieve role for
// return: the role id of the server owner role for a guild. An error if the role is not set in config.json
func getGuildRole(guildID string) (string, error) {
	c := utils.LoadConfig()
	for _, guild := range c.NewMemberGreeting.Config {
		if guild.ServerID == guildID {
			return guild.RoleID, nil
		}
	}
	return "", errors.New("no role set for this server")
}

// removeAllRole removes all server owner role from all greeters in a server.
// sess: the discord session.
// guildID: the guild to remove user roles from.
// roleID: the role to remove.
// users: the users to remove the role from.
// return: an error if any.
func removeAllRole(sess *discordgo.Session, guildID, roleID string, users []discordgo.User) error {
	for _, user := range users {
		err := sess.GuildMemberRoleRemove(guildID, user.ID, roleID)
		if err != nil {
			return err
		}
	}
	return nil
}
