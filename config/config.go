package config

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// Config configuration for the bot
type Config struct {
	// Discord token used to connect to discord
	DiscordToken string `yaml:"discord_token"`
	// If the bot should generate `README` docs on startup
	GenerateDocs bool `yaml:"generate_docs"`
	// Configuration related to the DB
	DB struct {
		UserName string `yaml:"username"`
		Password string `yaml:"password"`
		URL      string `yaml:"url"`
		DBName   string `yaml:"db_name"`
	} `yaml:"db"`
	Ollama struct {
		Host string `yaml:"host" validate:"omitempty,url"`
		// Name of the model to use
		Model string `yaml:"model" validate:"required"`
		// Request to ollama timeout duration
		Timeout string `yaml:"timeout"`
	} `yaml:"ollama"`
	LogLevel string `yaml:"logLevel"`
	DevMode  bool   `yaml:"dev_mode"`
	// RapidAPIKey string `yaml:"-"`
	// ID of the bot owner
	BotOwner string `yaml:"botOwner"`
	Emotes   []struct {
		// name of emote
		Name string `yaml:"name"`
		// url to emote
		URL string `yaml:"url"`
	} `yaml:"emotes"`
	// config for new member greetings
	NewMemberGreeting struct {
		Config []struct {
			ServerName string `yaml:"serverName"`
			RoleID     string `yaml:"roleID"`
			ServerID   string `yaml:"serverID"`
			ChannelID  string `yaml:"channelID"`
			Enable     bool   `yaml:"enable"`
		} `yaml:"config"`
	} `yaml:"newMemberGreeting"`
	Ivan struct {
		Emotes []struct {
			Name         string `yaml:"name"`
			FileLocation string `yaml:"fileLocation"`
		} `yaml:"emotes"`
	} `yaml:"ivan"`
}

// LoadConfig loads configuration from config.yml
func LoadConfig() Config {
	var c Config
	// read json file
	f, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	yaml.Unmarshal(f, &c)
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(c)
	if err != nil {
		log.Fatalf("Invalid config file format: %s", err)
	}

	return c
}
