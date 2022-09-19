package dd

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
)

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
	const Size = 1024
	ctx := gg.NewContext(Size, Size)
	ctx.SetRGB(1, 1, 1)
	ctx.Clear()
	ctx.SetRGB(0, 0, 0)

	if err := ctx.LoadFontFace("./media/font/comic_sans/comicz.ttf", 96); err != nil {
		log.Fatalln(err)
	}

	// NOTE: idk why this is necessary
	// ctx.DrawStringAnchored(optionMap["text"].StringValue(), Size/6, Size/6, 0.50, 0.10)
	ctx.DrawRoundedRectangle(0, 0, 512, 512, 0)
	ctx.DrawImage(img, 0, 0)
	// TODO: change values here to put the text where it should be
	ctx.DrawStringAnchored(optionMap["text"].StringValue(), Size/6, Size/6, 0.50, 6.0)
	ctx.Clip()
	ctx.SavePNG("out.png")

	out, err := os.Open("out.png")
	defer out.Close()
	if err != nil {
		log.Fatalln(err)
	}

	// outImg, outImageType, err := image.Decode(out)

	// dg.ChannelFileSend(i.ChannelID, "out.png", out)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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

	// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// Type: discordgo.InteractionResponseChannelMessageWithSource,
	// Data: &discordgo.InteractionResponseData{
	// Content: "Here is your dragon!",
	// Flags:   discordgo.MessageFlagsEphemeral,
	// },
	// },
	// )
}
