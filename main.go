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

	// "github.com/shawnyu5/debate_dragon_2.0/commands/snipe"
	generatedocs "github.com/shawnyu5/debate_dragon_2.0/generate_docs"
	"github.com/shawnyu5/debate_dragon_2.0/middware"
	utils "github.com/shawnyu5/debate_dragon_2.0/utils"
)

var cfg config.Config
var dg *discordgo.Session

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

// init reads config.json and sets global config variable
func init() {
	cfg = config.LoadConfig()
}

func init() {
	var err error
	dg, err = discordgo.New(fmt.Sprintf("Bot %s", cfg.DiscordToken))
	// dg, err = discordgo.New("Bot " + c.DiscordToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

// a handler function type for slash command and components
type handlerFunc func(ctx context.Context, sess *discordgo.Session, i *discordgo.InteractionCreate)

var (
	// array of slash command definitions
	slashCommandDefs = command.GetCmdDefs()
	// array of command handlers
	commandHandlers = command.GetCmdHandler()
	// array of component handlers
	componentsHandlers = command.GetComponentHandler()
)

func main() {
	go func() {
		log.Info("Generating docs...")
		generatedocs.Generate()
	}()

	ctx := context.Background()

	// 3. Open a connection pool to your PostgreSQL database
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.DB.UserName, cfg.DB.Password, cfg.DB.URL, cfg.DB.DBName)
	// TODO: look into tunning this configuration
	pool, err := pgxpool.New(ctx, dbUrl)
	store := db.NewStore(pool)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// Goose requires a standard standard-library database/sql *DB connection.
	// We can easily extract one out of our pgxpool without creating a new connection.
	stdlibDb := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer stdlibDb.Close()

	// Tell Goose to read migrations from our embedded Go filesystem variables
	goose.SetBaseFS(embedMigrations)

	// Set the driver type to postgres
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Goose failed to set dialect: %v", err)
	}

	log.Infof("Checking and running migrations...")
	// Automatically run "up" migrations located in the embedded "migrations" folder
	if err := goose.Up(stdlibDb, "sql/migrations"); err != nil {
		log.Fatalf("Goose migration failed: %v", err)
	}
	log.Infof("Database migration complete!")

	// // 4. Initialize sqlc's generated Queries struct using your connection pool
	// queries := db.New(pool)
	//
	// // 5. Call your type-safe query method!
	// // This matches the `-- name: ListMessages :many` annotation you wrote
	// messages, err := queries.ListMessages(ctx)
	// if err != nil {
	// 	log.Fatalf("Failed to fetch messages: %v\n", err)
	// }
	//
	// // 6. Iterate over the cleanly typed results
	// fmt.Printf("Found %d messages:\n", len(messages))
	// for _, msg := range messages {
	// 	// fields like ID, Content, CreatedAt are automatically generated inside db.Message
	// 	fmt.Printf("[%d] %s\n", msg.ID, msg.Content)
	// }

	// create database dir
	// os.Mkdir(c.DbPath, 0755)
	log.SetLevel(log.DebugLevel)
	log.Info("Starting bot...")

	dg.Identify.Intents |= discordgo.IntentGuildMessages
	dg.Identify.Intents |= discordgo.IntentGuildMembers
	dg.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
		log.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	dg.AddHandler(func(_ *discordgo.Session, mess *discordgo.MessageDelete) {
		messagetracking.TrackDeletedMessage(mess.GuildID, mess.ID)
	})

	dg.AddHandler(func(sess *discordgo.Session, mess *discordgo.MessageCreate) {
		messagetracking.TrackAllSentMessage(store, mess)
		areushawnyu.ListenForShawnYuMessages(sess, mess)
		stfu.TellUser(sess, mess)
	})

	dg.AddHandler(func(sess *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		// slash command handler
		case discordgo.InteractionApplicationCommand:
			if cmd, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				logger := middware.NewLogger(cmd)
				logger.HandleInteractionApplicationCommand(ctx, sess, i)
			} else {
				utils.SendErrorMessage(sess, i, "")
			}
		case discordgo.InteractionApplicationCommandAutocomplete:
			if cmd, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				logger := middware.NewLogger(cmd)
				logger.InteractionApplicationCommandAutocomplete(ctx, sess, i)
			} else {
				utils.SendErrorMessage(sess, i, "")
			}
		case discordgo.InteractionMessageComponent:
			if handlerFunc, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				command := command.Command{
					// This field is needed for `HandleInteractionApplicationCommand()`'s logging
					ApplicationCommand: func() *discordgo.ApplicationCommand {
						return &discordgo.ApplicationCommand{
							Name: i.MessageComponentData().CustomID,
						}
					},
					EditInteractionResponse: handlerFunc,
				}

				logger := middware.NewLogger(command)
				logger.HandleInteractionApplicationCommand(ctx, sess, i)
			} else {
				utils.SendErrorMessage(sess, i, "")
			}
		}
	})

	if err := dg.Open(); err != nil {
		log.Fatalf("Cannot open the discord session: %v", err)
	}

	// Only recreate commands in PROD
	// Dont re register slash commands on every restart in DEV
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
