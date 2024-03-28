package middware

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/shawnyu5/debate_dragon_2.0/command"
)

// Slash command handler logger
type Logger struct {
	Logger *log.Logger
	Next   command.Command
}

// handler calls discord slash command handler with logging
func (l Logger) EditIteractionResponse(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	output, err := l.Next.EditInteractionResponse(sess, i)
	defer func(begin time.Time, output string) {
		l.Logger.Infof("command=%s edited interaction response='%s' err=%s took=%s", l.Next.ApplicationCommand().Name, output, err, time.Since(begin))
	}(time.Now(), output)

	return output, err
}

// handler calls discord slash command handler with logging
func (l Logger) HandlerFunc(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	output, err := l.Next.HandlerFunc(sess, i)
	defer func(begin time.Time, output string) {
		l.Logger.Infof("command=%s response='%s' err=%s took=%s", l.Next.ApplicationCommand().Name, output, err, time.Since(begin))
	}(time.Now(), output)

	return output, err
}

// handler calls discord slash command handler with logging
func (l Logger) InteractionApplicationCommandAutocomplete(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	output, err := l.Next.InteractionApplicationCommandAutocomplete(sess, i)
	defer func(begin time.Time, output string) {
		l.Logger.Infof("command=%s auto complete response='%s' err=%s took=%s", l.Next.ApplicationCommand().Name, output, err, time.Since(begin))
	}(time.Now(), output)

	return output, err
}

// NewLogger creates a new logger middware.
//
// logger: logger to use.
//
// next: next middware in chain.
func NewLogger(logger *log.Logger, next command.Command) Logger {
	return Logger{Logger: logger, Next: next}
}
