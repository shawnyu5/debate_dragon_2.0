package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands/dd"
	"github.com/shawnyu5/debate_dragon_2.0/commands/insult"
	utils "github.com/shawnyu5/debate_dragon_2.0/utils"
)

type config struct {
	Token       string `json:"token"`
	TokenDev    string `json:"token_dev"`
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

var c config

var dg *discordgo.Session

// init reads config.json and sets global config variable
func init() {
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
	commandNames = map[string]string{
		"dd":     "dd",
		"insult": "insult",
	}

	commands = []*discordgo.ApplicationCommand{
		dd.Obj(),
		insult.Obj(),
	}

	commandHandlers = map[string]func(sess *discordgo.Session, i *discordgo.InteractionCreate){
		commandNames["dd"]:     dd.Handler,
		commandNames["insult"]: insult.Handler,
	}
)

func init() {
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
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
		utils.RemoveCommands(dg, registeredCommands)
	}

	log.Println("Gracefully shutting down.")
}
