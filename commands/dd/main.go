package dd

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
)

// Obj return a discord ApplicationCommand object defining this command
func Obj() *discordgo.ApplicationCommand {
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

// Handler a handler function for debate dragon
func Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	img, err := gg.LoadImage("./media/img/dragon_drawing.png")
	if err != nil {
		log.Fatalln(err)
	}
	// CanvasSize := 1024
	const CanvasSize = 1024
	fontSize := 70
	x := 2.3
	y := 2.1
	fontSize = shrinkFontSize(fontSize, optionMap["text"].StringValue(), 7)
	x, y = adjustTextPos(x, y, optionMap["text"].StringValue())
	fmt.Println(fmt.Sprintf("Handler x: %v", x)) // __AUTO_GENERATED_PRINT_VAR__

	ctx := gg.NewContext(CanvasSize, CanvasSize)
	ctx.SetRGB(1, 1, 1)
	ctx.Clear()
	ctx.SetRGB(0, 0, 0)

	if err := ctx.LoadFontFace("./media/font/comic_sans/comicz.ttf", float64(fontSize)); err != nil {
		log.Fatalln(err)
	}

	ctx.DrawRoundedRectangle(0, 0, 512, 512, 0)
	ctx.DrawImage(img, 0, 0)

	ctx.DrawStringAnchored(optionMap["text"].StringValue(), float64(CanvasSize/2), float64(CanvasSize/2), x, y)
	ctx.Clip()
	ctx.SavePNG("out.png")

	out, err := os.Open("out.png")
	defer out.Close()
	if err != nil {
		log.Fatalln(err)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
}

// shrinkFontSize shrink the font size passed in based on the length of user input and the maxCharacterSize
// Returns the new font size
func shrinkFontSize(fontSize int, userInput string, maxCharacterSize int) int {
	// 7 is the max character at current size
	if len(userInput) > maxCharacterSize {
		return shrinkFontSize(fontSize-5, userInput, maxCharacterSize+5)
	}
	return fontSize
}

// adjustTextPos adjust the text position based on the length of user input
// returns the adjusted x and y positions
func adjustTextPos(x, y float64, userInput string) (float64, float64) {
	if len(userInput) > 5 {
		return adjustTextPos(x-0.4, y, userInput[:len(userInput)-3])
	}

	return x, y
}
