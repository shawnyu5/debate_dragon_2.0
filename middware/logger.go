package middware

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
)

// Slash command handler logger
type Logger struct {
	Logger *log.Logger
	Next   commands.CommandStruct
}

// handler calls discord slash command handler with logging
func (l Logger) Handle(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	output, err := l.Next.Handler(sess, i)
	defer func(begin time.Time, output string) {
		l.Logger.Printf("command=%s response='%s' err=%v took=%s", l.Next.Name, output, err, time.Since(begin))
	}(time.Now(), output)

	return output, err
}
