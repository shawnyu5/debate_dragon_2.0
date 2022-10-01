package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands/dd"
	"github.com/shawnyu5/debate_dragon_2.0/commands/insult"
	"github.com/shawnyu5/debate_dragon_2.0/commands/ivan"
	"github.com/shawnyu5/debate_dragon_2.0/commands/manageIvan"
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

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	// map of command names
	// commandNames = map[string]string{
	// "dd":     "dd",
	// "insult": "insult",
	// "ivan":   "ivan",
	// }

	commands = []*discordgo.ApplicationCommand{
		dd.CommandObj.Obj(),
		insult.CommandObj.Obj(),
		ivan.CommandObj.Obj(),
		manageIvan.CommandObj.Obj(),
	}

	commandHandlers = map[string]func(sess *discordgo.Session, i *discordgo.InteractionCreate){
		dd.CommandObj.Name:         dd.CommandObj.CommandHandler,
		insult.CommandObj.Name:     insult.CommandObj.CommandHandler,
		ivan.CommandObj.Name:       ivan.CommandObj.CommandHandler,
		manageIvan.CommandObj.Name: manageIvan.CommandObj.CommandHandler,
	}

	// "gobanIvan": func(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	// err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// Type: discordgo.InteractionResponseChannelMessageWithSource,
	// Data: &discordgo.InteractionResponseData{
	// Content:         "HELLO BUTTON",
	// Components:      []discordgo.MessageComponent{},
	// Embeds:          []*discordgo.MessageEmbed{},
	// AllowedMentions: &discordgo.MessageAllowedMentions{},
	// Files:           []*discordgo.File{},
	// Flags:           0,
	// Choices:         []*discordgo.ApplicationCommandOptionChoice{},
	// CustomID:        "",
	// Title:           "",
	// },
	// })
	// if err != nil {
	// log.Println(err)
	// }
	// },
)

func init() {
	dg.AddHandler(func(sess *discordgo.Session, i *discordgo.InteractionCreate) {
		// switch i.Type {
		// case discordgo.InteractionApplicationCommand:
		// if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
		// h(s, i)
		// }
		// case discordgo.InteractionMessageComponent:

		// if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
		// h(s, i)
		// }
		// }

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
			}
		}
	})
}

func main() {
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := dg.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	// remove old commands before adding new ones
	// utils.RemoveCommands(dg, registeredCommands)

	utils.RegisterCommands(dg, commands, registeredCommands)
	dg.AddHandler(func(sess *discordgo.Session, gld *discordgo.GuildCreate) {
		log.Printf("Bot added to new guild: %v", gld.Name)
		utils.RegisterCommands(dg, commands, registeredCommands)
	})

	defer dg.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	// TODO: commands are not being deleted in my own server
	// only remove commands in production
	if !c.Development {
		utils.RemoveCommands(dg, registeredCommands, c)
	}

	log.Println("Gracefully shutting down.")
}
