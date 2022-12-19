package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	courseoutline "github.com/shawnyu5/debate_dragon_2.0/commands/courseOutline"
	"github.com/shawnyu5/debate_dragon_2.0/commands/dd"
	"github.com/shawnyu5/debate_dragon_2.0/commands/emotes"
	"github.com/shawnyu5/debate_dragon_2.0/commands/insult"
	"github.com/shawnyu5/debate_dragon_2.0/commands/ivan"
	"github.com/shawnyu5/debate_dragon_2.0/commands/manageIvan"
	newmember "github.com/shawnyu5/debate_dragon_2.0/commands/newMember"
	"github.com/shawnyu5/debate_dragon_2.0/commands/poll"
	"github.com/shawnyu5/debate_dragon_2.0/commands/rmp"
	"github.com/shawnyu5/debate_dragon_2.0/commands/snipe"
	subforcarmen "github.com/shawnyu5/debate_dragon_2.0/commands/subForCarmen"
	generatedocs "github.com/shawnyu5/debate_dragon_2.0/generate_docs"
	"github.com/shawnyu5/debate_dragon_2.0/middware"
	utils "github.com/shawnyu5/debate_dragon_2.0/utils"
)

var c utils.Config
var dg *discordgo.Session

// init reads config.json and sets global config variable
func init() {
	c = utils.LoadConfig()
}

func init() {
	var err error
	if c.Development {
		dg, err = discordgo.New("Bot " + c.TokenDev)
	} else {
		dg, err = discordgo.New("Bot " + c.Token)
	}
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

// a handler function type for slash command and components
type handlerFunc func(sess *discordgo.Session, i *discordgo.InteractionCreate)

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	// array of all slash commands in this bot
	allCommands = []commands.Command{
		manageIvan.ManageIvan{},
		poll.Poll{},
		dd.DD{},
		insult.Insult{},
		ivan.Ivan{},
		rmp.Rmp{},
		subforcarmen.SubForCarmen{},
		courseoutline.Outline{},
		snipe.Snipe{},
		emotes.Emotes{},
		newmember.NewMember{},
	}

	// array of slash command defs
	slashCommandDefs = utils.GetCmdDefs(allCommands)
	// array of command handlers
	commandHandlers = utils.GetCmdHandler(allCommands)
	// array of component handlers
	componentsHandlers = utils.GetComponentHandler(allCommands)
)

func init() {
	dg.AddHandler(func(sess *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		// handle slash command response and autocomplete requests the same way
		case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
			if handle, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				cmdObj := commands.CommandStruct{
					Name:    i.ApplicationCommandData().Name,
					Handler: handle,
				}
				logger := middware.NewLogger(log.New(os.Stdout, "", log.LstdFlags), cmdObj)
				logger.Handler(sess, i)
			} else {
				utils.SendErrorMessage(sess, i, "")
			}
		case discordgo.InteractionMessageComponent:
			if handle, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				cmdObj := commands.CommandStruct{
					Handler: handle,
				}

				logger := middware.NewLogger(log.New(os.Stdout, "", log.LstdFlags), cmdObj)
				logger.Handler(sess, i)
			} else {
				utils.SendErrorMessage(sess, i, "")
			}
		}
	})
}

func main() {
	go func() {
		generatedocs.Generate()
	}()
	// create database dir
	os.Mkdir(c.DbPath, 0755)

	dg.Identify.Intents |= discordgo.IntentGuildMessages
	dg.Identify.Intents |= discordgo.IntentGuildMembers
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	dg.AddHandler(func(_ *discordgo.Session, mess *discordgo.MessageDelete) {
		snipe.LastDeletedMessage = mess
	})

	removeHandler := dg.AddHandler(func(sess *discordgo.Session, mess *discordgo.MessageCreate) {
		// fmt.Println(mess.Content)
		// subforcarmen.Listen(sess, mess.Message)
		snipe.TrackMessage(mess)
	})

	dg.AddHandler(func(sess *discordgo.Session, user *discordgo.GuildMemberAdd) {
		log.Println("new user entered the guild")
		time.AfterFunc(5*time.Second, func() {
			newmember.Greet(sess, user)
		})
	})

	if !c.SubForCarmen.On {
		log.Println("removing handler")
		removeHandler()
	}

	err := dg.Open()

	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(slashCommandDefs))

	// remove old commands before adding new ones
	// utils.RemoveCommands(dg, registeredCommands)

	utils.RegisterCommands(dg, slashCommandDefs, registeredCommands)
	dg.AddHandler(func(sess *discordgo.Session, gld *discordgo.GuildCreate) {
		log.Printf("Bot added to new guild: %v", gld.Name)
		utils.RegisterCommands(dg, slashCommandDefs, registeredCommands)
	})

	defer dg.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	// TODO: commands are not being deleted in my own server
	// only remove commands in production
	if !c.Development {
		utils.RemoveCommands(dg, registeredCommands)
	}

	log.Println("Gracefully shutting down.")
}
