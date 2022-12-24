package mock

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"os"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

// make fun of a user's last message
type Mock struct{}

// Components implements commands.Command
func (Mock) Components() []commands.Component {
	return nil
}

// Def implements commands.Command
func (Mock) Def() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Version:     "1.0.0",
		Name:        "mock",
		Description: "Mock a user using the sponge bob meme",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to mock",
				Required:    true,
			},
		},
	}
}

// Handler implements commands.Command
func (Mock) Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	userOptions := utils.ParseUserOptions(sess, i)
	user := userOptions["user"].UserValue(sess)
	mess, err := GetUserLastMessage(sess, userOptions["user"].UserValue(sess), i.ChannelID)
	if err != nil {
		return "", err
	}

	err = GenMeme(fmt.Sprintf("@%s: %s", user.Username, MockText(mess.Content)))
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Open("out.png")
	if err != nil {
		return "", err
	}

	defer out.Close()
	_, err = sess.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:        "out.png",
				ContentType: "image/png",
				Reader:      out,
			},
		},
	})

	if err != nil {
		return "", err
	}

	err = sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("<@%s> mocked", user.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s mocked", user.Username), nil
}

func GenMeme(text string) error {
	img, err := gg.LoadImage("./media/img/mocking_spongebob.jpg")
	if err != nil {
		return err
	}
	const CanvasWidth = 502
	const CanvasHeight = 353
	fontSize := 30
	fontSize = utils.ShrinkFontSize(fontSize, text, CanvasWidth-10)

	ctx := gg.NewContext(CanvasWidth, CanvasHeight)
	ctx.SetRGB(0, 0, 0)

	if err := ctx.LoadFontFace("./media/font/comic_sans/comicz.ttf", float64(fontSize)); err != nil {
		return err
	}

	// make a rectangle and put the image on to it
	ctx.DrawRoundedRectangle(0, 0, 400, 400, 0)
	ctx.DrawImage(img, 0, 0)

	ctx.SetColor(color.White)
	// center "Literally no one" at the top of the image
	ctx.DrawStringAnchored("Literally no one:", CanvasWidth/2, 12, 0.5, 0.5)

	ctx.DrawStringWrapped(text, CanvasWidth-500, CanvasHeight-70, 0, 0, CanvasWidth, 2, gg.AlignCenter)
	ctx.Clip()
	ctx.SavePNG("out.png")
	return nil
}

// GetUserLastMessage gets the last message sent by a user in a channel.
// sess: discord session.
// user: the user to get the message of.
// channel: the channel to get the message from.
// return: a discord message and any errors
func GetUserLastMessage(sess *discordgo.Session, user *discordgo.User, channelID string) (*discordgo.Message, error) {
	// get all messages in channel
	messages, err := sess.ChannelMessages(channelID, 100, "", "", "")
	if err != nil {
		return nil, err
	}

	// find the last message sent by user
	for _, message := range messages {
		if message.Author.ID == user.ID {
			return message, nil
		}
	}
	return nil, errors.New("no user message in channel")
}

// MockText mock a string by turning over other letter upper case.
// text: the text to mock.
// return: the mocked text.
func MockText(text string) string {
	runes := make([]rune, 0, len(text))
	var upper bool
	for _, c := range text {
		if unicode.IsLetter(c) {
			upper = !upper
			if upper {
				c = unicode.ToUpper(c)
			}
		}
		runes = append(runes, c)
	}
	return string(runes)
}
