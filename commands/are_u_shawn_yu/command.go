package areushawnyu

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/shawnyu5/debate_dragon_2.0/command"
)

type state struct {
	Active        bool
	UserID        *string
	GuildID       *string
	CancelHandler *time.Timer
}

var s = state{
	Active:        false,
	UserID:        nil,
	CancelHandler: &time.Timer{},
}

var cmd = command.Command{
	ApplicationCommand: func() *discordgo.ApplicationCommand {
		return &discordgo.ApplicationCommand{
			Version:     "2.0.0",
			Name:        "are-u-shawn-yu",
			Description: "DEPRECATED: Accuses a user of being Shawn Yu",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user",
					Description: "The user to accuse of being Shawn Yu",
					Type:        discordgo.ApplicationCommandOptionUser,
					Required:    true,
				},
			},
		}
	},
	InteractionRespond: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (successMsg string, err error) {
		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "DEPRECATED",
			},
		})
		return "Deprecation message sent", nil
	},
}

// ListenForShawnYuMessages checks if `s.Active` is true. If it is, and the sent message matches `s.UserID`, tell `s.UserID` they are Shawn Yu.
// And set `s.Active` to false.
func ListenForShawnYuMessages(sess *discordgo.Session, mess *discordgo.MessageCreate) {
	// If the command is not active, or if the message is in a different guild, do nothing
	if !s.Active || mess.GuildID != *s.GuildID {
		return
	}
	s.Active = false

	sess.ChannelTyping(mess.ChannelID)
	s.CancelHandler = time.AfterFunc(3*time.Second, func() {
		if mess.Author.ID == *s.UserID {
			log.Debugf("Accusing user @%s of being Shawn Yu", mess.Author.Username)
			_, err := sess.ChannelMessageSend(mess.ChannelID, fmt.Sprintf("<@%s> You are Shawn Yu", *s.UserID))
			if err != nil {
				log.Error(err)
			}

			sess.ChannelTyping(mess.ChannelID)
			s.CancelHandler = time.AfterFunc(7*time.Second, func() {
				log.Debugf("Telling user @%s that the bot will come after them", mess.Author.Username)
				_, err = sess.ChannelMessageSend(mess.ChannelID, fmt.Sprintf("<@%s> I will find you https://www.youtube.com/watch?v=jZOywn1qArI&ab_channel=GreatestMovieClips", *s.UserID))
			})

			if err != nil {
				log.Error(err)
			}
		}

	})
}

func init() {
	command.Register(cmd)
}
