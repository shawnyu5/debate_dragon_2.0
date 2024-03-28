package messagetracking

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

// contains all discord messages sent since the bot startup, map of guild id to message id to discord message
var allMessagesMap = make(map[string]map[string]discordgo.Message)

// map of guild id to author ID to list of deleted messages
var deletedMessagesMap = make(map[string]map[string][]discordgo.Message)

// TrackAllSentMessage tracks all sent messages in all guilds this bot is in for the duration of the bot's life time
//
// We will only keep track of 1000 messages per guild. When we reach the 1000 message limit, delete the oldest message
func TrackAllSentMessage(mess *discordgo.MessageCreate) {
	// if the map for the current guild doesn't exist, initialize it
	if allMessagesMap[mess.GuildID] == nil {
		allMessagesMap[mess.GuildID] = map[string]discordgo.Message{}
	}

	// add message to map
	log.Debugf("Tracking message from user %s in guild %s", mess.Author.Username, mess.GuildID)
	allMessagesMap[mess.GuildID][mess.ID] = *mess.Message

	// store the messages to be removed
	oldestMessage := discordgo.Message{}
	oldestTimeDuration := time.Duration(0)

	// if we have more than 1000 messages stored for this guild, then remove the oldest message
	if len(allMessagesMap[mess.GuildID]) > 1000 {
		log.Debug("Over 1000 messages, deleting oldest message for guild", mess.GuildID)
		for _, message := range allMessagesMap[mess.GuildID] {
			timeDiff := time.Since(message.Timestamp.UTC())
			if oldestTimeDuration < timeDiff {
				oldestTimeDuration = timeDiff
				oldestMessage = message
			}
		}
		delete(allMessagesMap[mess.GuildID], oldestMessage.ID)
	}
}

// TrackDeletedMessage tracks the last 10 deleted messages metadata for a guild, excluding their content. Use allMessagesMap to get the contents of those messages
//
// When there are 10 messages, the oldest message will be deleted.
func TrackDeletedMessage(guildID string, messageID string) {
	if deletedMessagesMap[guildID] == nil {
		log.Debugf("Creating deleted messages guild map for guild %s", guildID)
		deletedMessagesMap[guildID] = make(map[string][]discordgo.Message)

	}

	// We are only able to find deleted message that was deleted when the bot is alive.
	mess := GetMessageByID(guildID, messageID)
	// No contents means the message was deleted before the bot was alive, nothing we can do about it...
	if mess.Content == "" {
		return
	}

	if len(deletedMessagesMap[mess.GuildID][mess.Author.ID]) == 10 {
		deletedMessagesMap[mess.GuildID][mess.Author.ID] = deletedMessagesMap[mess.GuildID][mess.Author.ID][1:]
	}
	deletedMessagesMap[mess.GuildID][mess.Author.ID] = append(deletedMessagesMap[mess.GuildID][mess.Author.ID], mess)
	// fmt.Printf("%+v\n", deletedMessagesMap)
}

// GetMessageByID returns the message with the given `messageID` from the given `guildID`.
func GetMessageByID(guildID, messageID string) discordgo.Message {
	return allMessagesMap[guildID][messageID]
}

// GetDeletedMessagesByAuthorID returns a list of deleted messages in a specific `guildID` for a specific `authorID`
func GetDeletedMessagesByAuthorID(guildID, authorID string) []discordgo.Message {
	usrMsg := deletedMessagesMap[guildID][authorID]
	return usrMsg
}
