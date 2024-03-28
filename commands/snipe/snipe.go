package snipe

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

// map of guild id to message id to discord message
var AllMessages = make(map[string]map[string]discordgo.Message)
var LastDeletedMessage = discordgo.MessageDelete{
	Message: &discordgo.Message{
		GuildID: "",
		ID:      "",
	},
}

// List of deleted message IDs
var DeletedMessages []discordgo.MessageDelete

var snipe = command.Command{
	Name: "snipe",
	ApplicationCommand: func() *discordgo.ApplicationCommand {
		return &discordgo.ApplicationCommand{
			Version:     "1.0.0",
			Name:        "snipe",
			Description: "Get the contents of the last deleted message",
		}
	},
	HandlerFunc: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		utils.DeferReply(sess, i.Interaction)

		deletedMess := GetMessageByID(LastDeletedMessage.GuildID, LastDeletedMessage.ID)
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

// TrackMessage adds a message to allMessages map.
//
// We will only keep track of 1000 messages per guild. When we reach the 1000 message limit, delete the oldest message
//
// Deprecated: use messagetracking.TrackAllSentMessage() instead
func TrackMessage(mess *discordgo.MessageCreate) {
	// if the map doesn't exist, initialize it
	if AllMessages[mess.GuildID] == nil {
		AllMessages[mess.GuildID] = map[string]discordgo.Message{}
	}

	// add message to map
	AllMessages[mess.GuildID][mess.ID] = *mess.Message

	// store the messages to be removed
	oldestMessage := discordgo.Message{}
	oldestTimeDuration := time.Duration(0)

	// if we have more than 1000 messages stored for this guild, then remove the oldest message
	if len(AllMessages[mess.GuildID]) > 1000 {
		for _, message := range AllMessages[mess.GuildID] {
			timeDiff := time.Since(message.Timestamp.UTC())
			if oldestTimeDuration < timeDiff {
				oldestTimeDuration = timeDiff
				oldestMessage = message
			}
		}
		delete(AllMessages[mess.GuildID], oldestMessage.ID)
	}
}

// TrackDeletedMessage tracks the last 10 deleted messages, excluding their content.
//
// When there are 10 messages, the oldest message will be deleted.
//
// Deprecated: use messagetracking.TrackDeletedMessage() instead
func TrackDeletedMessage(mess discordgo.MessageDelete) {
	if len(DeletedMessages) == 10 {
		DeletedMessages = DeletedMessages[1:]
	}
	DeletedMessages = append(DeletedMessages, mess)
}

// GetMessageByID returns the message with the given id from the given guild.
// guildID: the guild id the message is from.
// messageID: the id of the message.
// return: the discord message that was deleted.
//
// Deprecated: use messagetracking.GetMessageByID() instead
func GetMessageByID(guildID, messageID string) discordgo.Message {
	return AllMessages[guildID][messageID]
}

func init() {
	command.Register(snipe)
}
