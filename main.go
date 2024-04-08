package main

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	areushawnyu "github.com/shawnyu5/debate_dragon_2.0/commands/are_u_shawn_yu"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/caramel_bot/bitch"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/caramel_bot/compliment"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/courseOutline"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/dd"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/emotes"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/insult"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/ivan"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/manageIvan"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/memes/mock"
	messagetracking "github.com/shawnyu5/debate_dragon_2.0/commands/messageTracking"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/reddit"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/rmp"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/snipe"
	"github.com/shawnyu5/debate_dragon_2.0/commands/stfu"

	// "github.com/shawnyu5/debate_dragon_2.0/commands/snipe"
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
	// // array of all slash commands in this bot
	// allCommands = []commands.Command{
	//    snipe.Snipe{},
	// }

	// array of slash command definitions
	slashCommandDefs = command.GetCmdDefs()
	// array of command handlers
	commandHandlers = command.GetCmdHandler()
	// array of component handlers
	componentsHandlers = command.GetComponentHandler()
)

func init() {
	dg.AddHandler(func(sess *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		// handle slash command response and autocomplete requests the same way
		case discordgo.InteractionApplicationCommand:
			if cmd, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				logger := middware.NewLogger(cmd)
				logger.HandleInteractionApplicationCommand(sess, i)
				// logger := middware.NewLogger(cmd)
				// if cmd.EditInteractionResponse != nil {
				//    logger.EditIteractionResponse(sess, i)
				// } else if cmd.HandlerFunc != nil {
				//    logger.HandlerFunc(sess, i)
				// }
			} else {
				utils.SendErrorMessage(sess, i, "")
			}
		case discordgo.InteractionApplicationCommandAutocomplete:
			if cmd, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				logger := middware.NewLogger(cmd)
				logger.InteractionApplicationCommandAutocomplete(sess, i)
			} else {
				utils.SendErrorMessage(sess, i, "")
			}
		case discordgo.InteractionMessageComponent:
			if handlerFunc, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				command := command.Command{
					EditInteractionResponse: handlerFunc,
				}

				logger := middware.NewLogger(command)
				logger.EditIteractionResponse(sess, i)
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
	// os.Mkdir(c.DbPath, 0755)
	log.SetLevel(log.DebugLevel)
	log.Info("Starting bot...")

	dg.Identify.Intents |= discordgo.IntentGuildMessages
	dg.Identify.Intents |= discordgo.IntentGuildMembers
	dg.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
		log.Info("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	dg.AddHandler(func(_ *discordgo.Session, mess *discordgo.MessageDelete) {
		messagetracking.TrackDeletedMessage(mess.GuildID, mess.ID)
	})

	dg.AddHandler(func(sess *discordgo.Session, mess *discordgo.MessageCreate) {
		messagetracking.TrackAllSentMessage(mess)
		areushawnyu.ListenForShawnYuMessages(sess, mess)
		stfu.TellUser(sess, mess)
	})

	if err := dg.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// utils.RegisterCommands(dg, slashCommandDefs, registeredCommands)
	registeredCommands := command.RegisterCommands(dg)
	dg.AddHandler(func(_ *discordgo.Session, gld *discordgo.GuildCreate) {
		log.Printf("Bot added to new guild: %v", gld.Name)
		command.RegisterCommands(dg)
	})

	// command.DiscoverCommands()

	defer dg.Close()

	dg.StateEnabled = true

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Info("Press Ctrl+C to exit")
	<-stop

	// TODO: commands are not being deleted in my own server
	// only remove commands in production
	if !c.Development {
		command.RemoveCommands(dg, registeredCommands)
	}

	log.Info("Gracefully shutting down.")
}
