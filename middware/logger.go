package middware

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
)

type Logger struct {
	Logger *log.Logger
	Next   commands.CommandStruct
}

func (l Logger) Obj() *discordgo.ApplicationCommand {
	return l.Next.Obj()
}

// handler discord slash command handler with logging
func (l Logger) Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func(begin time.Time) {
		l.Logger.Printf("command=%s response=%s took=%s", l.Next.Obj().Name, "idk yet", time.Since(begin))
	}(time.Now())
	l.Next.Handler(sess, i)
}
