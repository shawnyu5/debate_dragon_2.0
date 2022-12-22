package mock

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
)

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
	err := GenMeme()
	if err != nil {
		log.Fatal(err)
	}
	return "mock", nil
}

func GenMeme() error {
	img, err := gg.LoadImage("./media/img/mocking_spongebob.jpg")
	if err != nil {
		return err
	}
	const CanvasWidth = 502
	const CanvasHeight = 353
	const fontSize = 30
	ctx := gg.NewContext(CanvasWidth, CanvasHeight)
	ctx.SetRGB(0, 0, 0)

	if err := ctx.LoadFontFace("./media/font/comic_sans/comicz.ttf", float64(fontSize)); err != nil {
		return err
	}

	// make a rectangle and put the image on to it
	ctx.DrawRoundedRectangle(0, 0, 400, 400, 0)
	ctx.DrawImage(img, 0, 0)
	// center "Literally no one" at the top of the image
	ctx.DrawStringAnchored("Literally no one:", CanvasWidth/2, 12, 0.5, 0.5)

	ctx.Clip()
	ctx.SavePNG("out.png")
	return nil
}
