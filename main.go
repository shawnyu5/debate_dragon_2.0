package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	subforcarmen "github.com/shawnyu5/debate_dragon_2.0/commands/subForCarmen"
	generatedocs "github.com/shawnyu5/debate_dragon_2.0/generate_docs"
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
	allCommands = []commands.CommandStruct{
		// dd.CommandObj,
		// insult.CommandObj,
		// ivan.CommandObj,
		// manageIvan.CommandObj,
		// rmp.CommandObj,
		subforcarmen.CommandObj,
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
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(sess, i)
			} else {
				utils.SendErrorMessage(sess, i, "")
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(sess, i)
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
	dg.Identify.Intents |= discordgo.IntentGuildMessages
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	removeHandler := dg.AddHandler(func(sess *discordgo.Session, mess *discordgo.MessageCreate) {
		fmt.Println(mess.Content)
		subforcarmen.Listen(sess, mess.Message)
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
	// if !c.Development {
	utils.RemoveCommands(dg, registeredCommands, c)
	// }

	log.Println("Gracefully shutting down.")
}
