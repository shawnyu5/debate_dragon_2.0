package snipe

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	messagetracking "github.com/shawnyu5/debate_dragon_2.0/commands/messageTracking"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

// List of deleted message IDs
var DeletedMessages []discordgo.MessageDelete

var snipe = command.Command{
	ApplicationCommand: func() *discordgo.ApplicationCommand {
		return &discordgo.ApplicationCommand{
			Version:     "1.0.0",
			Name:        "snipe",
			Description: "Get the contents of the last deleted message",
		}
	},
	InteractionRespond: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		utils.DeferReply(sess, i.Interaction)

		deletedMess := messagetracking.GetLastDeletedMessage()
		// deletedMess := GetMessageByID(LastDeletedMessage.GuildID, LastDeletedMessage.ID)
		// if there are no deleted messages in cache, then send error response
		if deletedMess.Content == "" && len(deletedMess.Attachments) == 0 {
			_, err := sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Type:        discordgo.EmbedTypeArticle,
						Title:       "Last Deleted Message",
						Description: "No one has deleted a message in a while...",
					},
				},
			})

			if err != nil {
				return "", err
			}
			return "No deleted message cached", nil
		}
		webHookEdit := &discordgo.WebhookEdit{}

		// if there is an image, send it
		if len(deletedMess.Attachments) > 0 {
			webHookEdit = &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Type:        discordgo.EmbedTypeArticle,
						Title:       "Sniped",
						Description: fmt.Sprintf("<@%s>", deletedMess.Author.ID),
						Image: &discordgo.MessageEmbedImage{
							URL: deletedMess.Attachments[0].URL,
						},
					},
				},
			}

		} else {
			// otherwise send the message
			webHookEdit = &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						URL:         "",
						Type:        discordgo.EmbedTypeArticle,
						Title:       "Snipe",
						Description: fmt.Sprintf("%s - <@%s>", deletedMess.Content, deletedMess.Author.ID),
					},
				},
			}

		}
		_, err := sess.InteractionResponseEdit(i.Interaction, webHookEdit)
		if err != nil {
			return "", err
		}
		return "deleted message sent", nil
	},
}

func init() {
	command.Register(snipe)
}
