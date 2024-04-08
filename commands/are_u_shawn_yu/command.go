package areushawnyu

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

type state struct {
	Active        bool
	UserID        *string
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
			Version:     "1.0.0",
			Name:        "are-u-shawn-yu",
			Description: "Accuses a user of being Shawn Yu",
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
	EditInteractionResponse: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		// If the command is currently active, cancel the previous iteration before continuing
		if s.Active {
			s.CancelHandler.Stop()
			s.Active = false
		}

		userOptions := utils.ParseUserOptions(sess, i)
		user := userOptions["user"].UserValue(sess)

		content := fmt.Sprintf("<@%s> are you Shawn Yu?", user.ID)
		_, err := sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		if err != nil {
			return "", err
		}

		s.Active = true
		s.UserID = &user.ID
		return fmt.Sprintf("Asked @%s if they are Shawn Yu", user.Username), nil
	},
}

// ListenForShawnYuMessages checks if `s.Active` is true. If it is, and the sent message matches `s.UserID`, tell `s.UserID` they are Shawn Yu.
// And set `s.Active` to false.
func ListenForShawnYuMessages(sess *discordgo.Session, mess *discordgo.MessageCreate) {
	if !s.Active {
		return
	}
	s.Active = false

	sess.ChannelTyping(mess.ChannelID)
	s.CancelHandler = time.AfterFunc(3*time.Second, func() {
		if mess.Author.ID == *s.UserID {
			log.Debugf("Accusing user %s of being Shawn Yu", mess.Author.Username)
			_, err := sess.ChannelMessageSend(mess.ChannelID, fmt.Sprintf("<@%s> You are Shawn Yu", *s.UserID))
			if err != nil {
				log.Error(err)
			}

			sess.ChannelTyping(mess.ChannelID)
			s.CancelHandler = time.AfterFunc(7*time.Second, func() {
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
