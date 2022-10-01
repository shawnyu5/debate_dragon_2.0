package manageIvan

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

var CommandObj = commands.CommandStruct{
	Name:    "manageivan",
	Obj:     obj,
	CommandHandler: handler,
}

func obj() *discordgo.ApplicationCommand {
	defaultMemberPermissions := int64(discordgo.PermissionManageServer)

	return &discordgo.ApplicationCommand{
		Version:                  "1.0",
		Name:                     "manageivan",
		DefaultMemberPermissions: &defaultMemberPermissions,
		Description:              "Command to help the management of Ivan",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionUser,
				Name:         "user",
				Description:  "Ivan account to ban",
				Required:     true,
				Autocomplete: false,
			},
		},
	}
}

// handler the handler for `/manageivan` command
func handler(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	utils.DeferReply(sess, i.Interaction)
	optionsMap := utils.ParseUserOptions(sess, i)
	fmt.Println(fmt.Sprintf("handler optionsMap: %v", optionsMap)) // __AUTO_GENERATED_PRINT_VAR__

	content := "YOOO"
	_, err := sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Confirmation",
						Style:    discordgo.DangerButton,
						Disabled: false,
						Emoji:    discordgo.ComponentEmoji{},
						// URL:      "https://google.com",
						CustomID: "gobanIvan",
					},
				},
			},
		},
		// Embeds:          &[]*discordgo.MessageEmbed{},
		// Files:           []*discordgo.File{},
		// AllowedMentions: &discordgo.MessageAllowedMentions{},
	})

	if err != nil {
		utils.SendErrorMessage(sess, i, err.Error())
		log.Println(err)
	}
}
