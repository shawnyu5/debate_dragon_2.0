package blackmail

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	messagetracking "github.com/shawnyu5/debate_dragon_2.0/commands/messageTracking"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

var blackmail = command.Command{
	Name: "blackmail",
	ApplicationCommand: func() *discordgo.ApplicationCommand {
		return &discordgo.ApplicationCommand{
			Version:     "1.0.0",
			Name:        "blackmail",
			Description: "Retrieve the last 10 deleted message from a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user",
					Description: "The user who's deleted messages will be retrieved",
					Type:        discordgo.ApplicationCommandOptionUser,
					Required:    true,
				},
			},
		}
	},
	HandlerFunc: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		utils.DeferReply(sess, i.Interaction)
		input := utils.ParseUserOptions(sess, i)
		user := input["user"].UserValue(sess)
		msgs := messagetracking.GetDeletedMessagesByAuthorID(i.GuildID, user.ID)

		if len(msgs) == 0 {
			content := fmt.Sprintf("<@%s> hasnt deleted any messages recently", user.ID)
			sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})
			return fmt.Sprintf("User %s has not deleted any messages recently", user.Username), nil
		}

		content := fmt.Sprintf("<@%s>'s last %d deleted messages", user.ID, len(msgs))
		embed := constructEmbed(msgs)
		sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
			Embeds: &[]*discordgo.MessageEmbed{
				embed,
			},
		})
		return fmt.Sprintf("Send deleted message for user %s", user), nil
	},
}

// constructEmbed given a list of messages, construct a single embed containing all the messages and their send date
func constructEmbed(msgs []discordgo.Message) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Description: "",
	}
	for _, msg := range msgs {
		timeStamp, err := discordgo.SnowflakeTimestamp(msg.ID)
		if err != nil {
			log.Printf("Failed to calculate time stamp from message ID %s for message %s", msg.ID, msg.Content)
			timeStamp = msg.Timestamp
			// continue
		}

		embed.Description = fmt.Sprintf("%s\n- %s (sent on %s)", embed.Description, msg.Content, timeStamp.Format("2006-01-02"))
	}
	return &embed
}

func init() {
	command.Register(blackmail)
}
