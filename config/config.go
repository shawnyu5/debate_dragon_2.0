package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config configuration for the bot
type Config struct {
	// Discord token used to connect to discord
	DiscordToken string `yaml:"discord_token"`
	// Configuration related to the DB
	DB struct {
		UserName string `yaml:"username"`
		Password string `yaml:"password"`
		URL      string `yaml:"url"`
		DBName   string `yaml:"db_name"`
	} `yaml:"db"`
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
	SubForCarmen struct {
		// toggle this feature on and off
		On bool `yaml:"on"`
		// id of carmen user to track messages of
		CarmenID string `yaml:"carmenId"`
		// cool down, defined in minutes
		CoolDown int `yaml:"coolDown"`
		// the guild to keep track of carmen messages
		GuildID string `yaml:"guildID"`
		// number of messages before a notification is triggered
		MessageLimit      int    `yaml:"messageLimit"`
		SubscribersRoleID string `yaml:"subscribersRoleID"`
		// channels to ignore
		IgnoredChannels []string `yaml:"ignoredChannels"`
	} `yaml:"subForCarmen"`
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
	return c
}
