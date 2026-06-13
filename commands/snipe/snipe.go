package snipe

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	messagetracking "github.com/shawnyu5/debate_dragon_2.0/commands/messageTracking"
	"github.com/shawnyu5/debate_dragon_2.0/db"
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
	InteractionRespond: func(ctx context.Context, sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		utils.DeferReply(sess, i.Interaction)
		db, err := db.StoreFromContext(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to find DB in context: %s. This is a bug", err)
		}

		deletedMsg, err := messagetracking.GetDeletedMessageByGuildID(db, i.GuildID)
		if err != nil {
			return "", fmt.Errorf("failed to get deleted message: %s", err)
		}

		log.Infof("Deleted message from dB: %+v", deletedMsg)

		if deletedMsg.Content == "" {
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

		if len(deletedMsg.Attachments) > 0 {
			webHookEdit = &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Type:        discordgo.EmbedTypeArticle,
						Title:       "Sniped",
						Description: fmt.Sprintf("<@%s>", deletedMsg.AuthorID),
						Image: &discordgo.MessageEmbedImage{
							URL: deletedMsg.Attachments[0].URL,
						},
					},
				},
			}
		} else {
			webHookEdit = &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						URL:         "",
						Type:        discordgo.EmbedTypeArticle,
						Title:       "Snipe",
						Description: fmt.Sprintf("%s - <@%s>", deletedMsg.Content, deletedMsg.AuthorID),
					},
				},
			}
		}

		_, err = sess.InteractionResponseEdit(i.Interaction, webHookEdit)
		if err != nil {
			return "", err
		}
		return "deleted message sent", nil
	},
}

func init() {
	command.Register(snipe)
}
