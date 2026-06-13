package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/release_notes"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/rmp"
	_ "github.com/shawnyu5/debate_dragon_2.0/commands/snipe"
	"github.com/shawnyu5/debate_dragon_2.0/commands/stfu"
	"github.com/shawnyu5/debate_dragon_2.0/config"
	"github.com/shawnyu5/debate_dragon_2.0/db"

	generatedocs "github.com/shawnyu5/debate_dragon_2.0/generate_docs"
	"github.com/shawnyu5/debate_dragon_2.0/middware"
	utils "github.com/shawnyu5/debate_dragon_2.0/utils"
)

var cfg config.Config
var dg *discordgo.Session

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

func init() {
	cfg = config.LoadConfig()
}

func init() {
	var err error
	dg, err = discordgo.New(fmt.Sprintf("Bot %s", cfg.DiscordToken))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	slashCommandDefs   = command.GetCmdDefs()
	commandHandlers    = command.GetCmdHandler()
	componentsHandlers = command.GetComponentHandler()
)

func main() {
	go func() {
		log.Info("Generating docs...")
		generatedocs.Generate()
	}()

	ctx := context.Background()

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.DB.UserName, cfg.DB.Password, cfg.DB.URL, cfg.DB.DBName)
	pool, err := pgxpool.New(ctx, dbUrl)
	store := db.NewStore(pool)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	stdlibDb := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer stdlibDb.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Goose failed to set dialect: %v", err)
	}

	log.Infof("Checking and running migrations...")
	if err := goose.Up(stdlibDb, "sql/migrations"); err != nil {
		log.Fatalf("Goose migration failed: %v", err)
	}
	log.Infof("Database migration complete!")

	log.SetLevel(log.DebugLevel)
	log.Info("Starting bot...")

	dg.Identify.Intents |= discordgo.IntentGuildMessages
	dg.Identify.Intents |= discordgo.IntentGuildMembers

	dg.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
		log.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	dg.AddHandler(func(_ *discordgo.Session, mess *discordgo.MessageDelete) {
		messagetracking.TrackDeletedMessage(mess)
	})

	dg.AddHandler(func(sess *discordgo.Session, mess *discordgo.MessageCreate) {
		messagetracking.TrackAllSentMessage(store, mess)
		areushawnyu.ListenForShawnYuMessages(sess, mess)
		stfu.TellUser(sess, mess)
	})

	dg.AddHandler(handleInteraction)

	if err := dg.Open(); err != nil {
		log.Fatalf("Cannot open the discord session: %v", err)
	}

	if !cfg.DevMode {
		err := command.RemoveAllSlashCommandFromAllGuilds(dg)
		if err != nil {
			log.Fatalf("failed to remove slash command: %s", err)
		}
	}

	log.Info("Registering slash commands")
	command.RegisterCommands(dg)

	dg.AddHandler(func(_ *discordgo.Session, gld *discordgo.GuildCreate) {
		log.Printf("Bot added to new guild: %v", gld.Name)
		command.RegisterCommands(dg)
	})

	defer dg.Close()

	dg.StateEnabled = true

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Info("Press Ctrl+C to exit")
	<-stop

	log.Info("Gracefully shutting down.")
}

func handleInteraction(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if cmd, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			logger := middware.NewLogger(cmd)
			logger.HandleInteractionApplicationCommand(context.Background(), sess, i)
		} else {
			utils.SendErrorMessage(sess, i, "")
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		if cmd, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			logger := middware.NewLogger(cmd)
			logger.InteractionApplicationCommandAutocomplete(context.Background(), sess, i)
		} else {
			utils.SendErrorMessage(sess, i, "")
		}
	case discordgo.InteractionMessageComponent:
		if handlerFunc, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
			command := command.Command{
				ApplicationCommand: func() *discordgo.ApplicationCommand {
					return &discordgo.ApplicationCommand{
						Name: i.MessageComponentData().CustomID,
					}
				},
				EditInteractionResponse: handlerFunc,
			}

			logger := middware.NewLogger(command)
			logger.HandleInteractionApplicationCommand(context.Background(), sess, i)
		} else {
			utils.SendErrorMessage(sess, i, "")
		}
	default:
		log.Warnf("Unhandled interaction type: %v", i.Type)
	}
}
