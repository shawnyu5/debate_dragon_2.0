package dd

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

type DD struct{}

// Components implements commands.Command
// this command does not have components
func (DD) Components() []commands.Component {
	return []commands.Component{}
}

// Def implements commands.Command
func (DD) Def() *discordgo.ApplicationCommand {
	obj := &discordgo.ApplicationCommand{
		Name:        "dd",
		Description: "summon a dragon to burn your debate floes to the ground.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "text",
				Description: "text to burn down your floe",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	}
	return obj
}

// Handler implements commands.Command
func (DD) Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	userOptions := utils.ParseUserOptions(sess, i)

	img, err := gg.LoadImage("./media/img/dragon_drawing.png")
	if err != nil {
		log.Fatalln(err)
	}
	const CanvasSize = 1024
	fontSize := 70
	fontSize = ShrinkFontSize(fontSize, userOptions["text"].StringValue(), 7)

	ctx := gg.NewContext(CanvasSize, CanvasSize)
	ctx.SetRGB(1, 1, 1)
	ctx.Clear()
	ctx.SetRGB(0, 0, 0)

	if err := ctx.LoadFontFace("./media/font/comic_sans/comicz.ttf", float64(fontSize)); err != nil {
		log.Fatalln(err)
	}

	ctx.DrawRoundedRectangle(0, 0, 512, 512, 0)
	ctx.DrawImage(img, 0, 0)

	// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
	// text. Use ax=0.5, ay=0.5 to center the text at the specified point
	x := 20
	y := 0.1
	ctx.DrawStringWrapped(userOptions["text"].StringValue(), float64(CanvasSize/2), float64(CanvasSize/2), float64(x), float64(y), 15, 1, gg.AlignCenter)
	ctx.Clip()
	ctx.SavePNG("out.png")

	out, err := os.Open("out.png")
	defer out.Close()
	if err != nil {
		log.Fatalln(err)
	}

	err = sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Files: []*discordgo.File{
				{
					Name:        "out.png",
					ContentType: "image/png",
					Reader:      out,
				},
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	return userOptions["text"].StringValue(), nil
}

// ShrinkFontSize shrink the font size passed in based on the length of user input and the maxCharacterSize
// Returns the new font size
func ShrinkFontSize(fontSize int, userInput string, maxCharacterSize int) int {
	// 7 is the max character at current size
	if len(userInput) > maxCharacterSize {
		return ShrinkFontSize(fontSize-5, userInput, maxCharacterSize+5)
	}
	return fontSize
}
