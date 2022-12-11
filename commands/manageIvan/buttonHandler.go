package manageIvan

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

// handleDontBan handle when the dont ban button is pushed
func handleDontBan(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Okay, <@%s> will not be banned... :(", ivanBanState.User.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						createBanButton(true),
						createDontBanButton(true),
						createJumpScareButton(true),
					},
				},
			},
		},
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Okay, <@%s> will not be banned... :(", ivanBanState.User.ID), nil
}

// handleJumpScare handle with the jump scare button is pushed
func handleJumpScare(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	// change original ephemeral message to command executor
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "Jump scare sequence initiated...",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						createBanButton(true),
						createDontBanButton(true),
						createJumpScareButton(true),
					},
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	// keep track of all sent messages so we cant delete them later
	sentMessages := []*discordgo.Message{}
	messages := GenerateMessages(ivanBanState.CountDownTime)

	// start count down
	for _, message := range messages {
		mess, err := sess.ChannelMessageSend(i.ChannelID, message.message)
		if err != nil {
			return "", err
		}
		sentMessages = append(sentMessages, mess)
		time.Sleep(message.countDownTime * time.Second)
	}
	mess, err := sess.ChannelMessageSend(i.ChannelID, fmt.Sprintf("jk, we ain't that mean, you will not be banned <@%s>", ivanBanState.User.ID))
	if err != nil {
		return "", err
	}
	sentMessages = append(sentMessages, mess)
	time.Sleep(5 * time.Second)

	mess, err = sess.ChannelMessageSend(i.ChannelID, "Good bye now...")
	if err != nil {
		return "", err
	}
	sentMessages = append(sentMessages, mess)

	time.Sleep(3 * time.Second)

	// clean up all messages
	utils.DeleteAllMessages(sess, i, sentMessages)
	return "Finished jumpscare ban", nil
}

// startBanningIvan handles when the ban button is pushed
func handleBan(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	// change original ephemeral message to command executor
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "Sequence initiated...",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				CreateAllButtons(true),
			},
		},
	})

	if err != nil {
		return "", err
	}

	// keep track of all sent messages so we cant delete them later
	sentMessages := []*discordgo.Message{}
	messages := GenerateMessages(ivanBanState.CountDownTime)

	// start count down
	for _, message := range messages {
		mess, err := sess.ChannelMessageSend(i.ChannelID, message.message)
		if err != nil {
			return "", err
		}
		sentMessages = append(sentMessages, mess)
		time.Sleep(message.countDownTime * time.Second)
	}

	if !config.Development {
		err = sess.GuildBanCreateWithReason(i.GuildID, ivanBanState.User.ID, "Ivan", 0)
		if err != nil {
			return "", err
		}
	}

	// time.Sleep(5 * time.Second)
	// send embed that user has been banned
	_, err = sess.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embed: &discordgo.MessageEmbed{
			URL:         "",
			Type:        "",
			Title:       "Ivan Ban",
			Description: fmt.Sprintf("<@%s> HAS BEEN BANNED", ivanBanState.User.ID),
			Timestamp:   "",
			Color:       0,
		},
	})

	if err != nil {
		log.Println(err)
	}

	time.Sleep(5 * time.Second)
	// clean up all messages
	utils.DeleteAllMessages(sess, i, sentMessages)
	return fmt.Sprintf("<@%s> HAS BEEN BANNED", ivanBanState.User.ID), nil
}

// handleKick handles when the kick button is pushed.
func handleKick(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Sequence initiated... Kicking <@%s> in progress", ivanBanState.User.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				CreateAllButtons(true),
			},
		},
	})
	if err != nil {
		return "", err
	}
	// inviteURL URL
	inviteURL, err := utils.CreateInvite(sess, i.GuildID)
	if err != nil {
		return "", err
	}

	channel, err := sess.UserChannelCreate(ivanBanState.User.ID)
	message := fmt.Sprintf("<@%s>You will be kicked in a wee bit. Here is an invite link if you wantta come back. Pls come back %s", ivanBanState.User.ID, inviteURL)
	// if unable to create a channel directly to the user, then send a message in chat before the user is kicked
	if err != nil {
		_, err := sess.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Content: message,
		})
		if err != nil {
			return "", err
		}
	}

	_, err = sess.ChannelMessageSend(channel.ID, message)
	// if unable to send message directly to user, then send a message in chat before the user is kicked
	if err != nil {
		_, err := sess.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Content: message,
		})
		if err != nil {
			return "", err
		}
	}

	var sentMessages []*discordgo.Message
	messages := GenerateMessages(ivanBanState.CountDownTime)
	for _, message := range messages {
		mess, err := sess.ChannelMessageSend(i.ChannelID, message.message)
		if err != nil {
			return "", err
		}
		sentMessages = append(sentMessages, mess)
		time.Sleep(message.countDownTime * time.Second)
	}

	if !config.Development {
		err = sess.GuildBanCreateWithReason(i.GuildID, ivanBanState.User.ID, "Ivan", 0)
		if err != nil {
			return "", err
		}
	}

	// send embed that user has been banned
	_, err = sess.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embed: &discordgo.MessageEmbed{
			URL:         "",
			Type:        "",
			Title:       "Ivan Ban",
			Description: fmt.Sprintf("<@%s> HAS BEEN KICKED", ivanBanState.User.ID),
			Timestamp:   "",
			Color:       0,
		},
	})
	if err != nil {
		return "", err
	}

	utils.DeleteAllMessages(sess, i, sentMessages)

	return fmt.Sprintf("<@%s> HAS BEEN KICKED", ivanBanState.User.ID), nil
}
