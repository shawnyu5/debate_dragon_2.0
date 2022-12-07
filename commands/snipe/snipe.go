package snipe

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

// command to get the last deleted message
type Snipe struct{}

// map of guild id: { message id : discord message }
var AllMessages = make(map[string]map[string]discordgo.Message)
var LastDeletedMessage = &discordgo.MessageDelete{}

// Components implements commands.Command
func (Snipe) Components() []commands.Component {
	return nil
}

// Def implements commands.Command
func (Snipe) Def() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Version:                  "1.0.0",
		Type:                     0,
		Name:                     "snipe",
		Description:              "Get the contents of the last deleted message",
		DescriptionLocalizations: &map[discordgo.Locale]string{},
		Options:                  []*discordgo.ApplicationCommandOption{},
	}
}

// Handler implements commands.Command
func (Snipe) Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	utils.DeferReply(sess, i.Interaction)

	deletedMess := GetMessageByID(LastDeletedMessage.GuildID, LastDeletedMessage.ID)
	// if there are no deleted messages in cache, then send error response
	if deletedMess.Content == "" && len(deletedMess.Attachments) == 0 {
		_, err := sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Type:        discordgo.EmbedTypeArticle,
					Title:       "Last Deleted Message",
					Description: "No deleted message",
				},
			},
		})

		if err != nil {
			return "", err
		}
		return "No deleted message cached", nil
	}
	_, err := sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:    new(string),
		Components: &[]discordgo.MessageComponent{},
		Embeds: &[]*discordgo.MessageEmbed{
			{
				URL:         "",
				Type:        discordgo.EmbedTypeArticle,
				Title:       "Snipe",
				Description: fmt.Sprintf("%s - <@%s>", deletedMess.Content, deletedMess.Author.ID),
				Timestamp:   "",
				Color:       0,
				Footer:      &discordgo.MessageEmbedFooter{},
				Image: &discordgo.MessageEmbedImage{
					URL: deletedMess.Attachments[0].URL,
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{},
				Video:     &discordgo.MessageEmbedVideo{},
				Provider:  &discordgo.MessageEmbedProvider{},
				Author:    &discordgo.MessageEmbedAuthor{},
				Fields:    []*discordgo.MessageEmbedField{},
			},
		},
		Files:           []*discordgo.File{},
		AllowedMentions: &discordgo.MessageAllowedMentions{},
	})
	if err != nil {
		return "", err
	}
	return "deleted message sent", nil
}

// trackMessage adds a message to allMessages map
// mess: the message to track
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

	// if we have more than 100 messages stored for this guild, then remove the oldest message
	if len(AllMessages[mess.GuildID]) > 100 {
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

// GetMessageByID returns the message with the given id from the given guild.
// guildID: the guild id the message is from.
// messageID: the id of the message.
// return: the discord message that was deleted.
func GetMessageByID(guildID, messageID string) discordgo.Message {
	return AllMessages[guildID][messageID]
}
