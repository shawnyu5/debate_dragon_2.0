package isIvan

import (
	"fmt"
	"log"

	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/go-cmd/cmd"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

var CommandObj = commands.CommandStruct{
	Obj:     obj,
	Handler: handler,
}

func obj() *discordgo.ApplicationCommand {
	obj := &discordgo.ApplicationCommand{
		Version:     "1.0",
		Name:        "isivan",
		Description: "Use machine learning to predict if a user is Ivan",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "user",
				Description:  "the user to check",
				Type:         discordgo.ApplicationCommandOptionUser,
				Required:     true,
				Autocomplete: false,
			},
			{
				Name:         "messagelimit",
				Description:  "number of user messages to check. Default: 50. NOTE more messages will take longer to predict",
				Required:     false,
				Autocomplete: false,
				MaxValue:     100,
				Type:         discordgo.ApplicationCommandOptionInteger,
			},
		},
	}
	return obj
}

func handler(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	optionsMap := utils.ParseUserOptions(sess, i)

	// if a user was passed in
	if optionsMap["user"] != nil {
		err := utils.DeferReply(sess, i.Interaction)
		messageLimit := 50

		if optionsMap["messageLimit"] != nil {
			messageLimit = int(optionsMap["messageLimit"].IntValue())
		}

		// get all users messages in a channel
		possibleIvanMessages := utils.GetAllUserMessageFromChannel(sess, i.ChannelID, messageLimit, optionsMap["user"].UserValue(sess).ID)

		isIvanPossiblities := make([]bool, 0)
		var wg sync.WaitGroup

		for i, message := range possibleIvanMessages {
			wg.Add(1)
			// check if the message is from ivan, and store the results in a slice
			go func(i int, message string) {
				defer wg.Done()
				isIvan := checkIsIvan(message)
				isIvanPossiblities = append(isIvanPossiblities, isIvan)
				fmt.Println(fmt.Sprintf("isIvanPossiblities: %v", isIvanPossiblities)) // __AUTO_GENERATED_PRINT_VAR__
			}(i, message)
		}

		wg.Wait()

		// calculate the number of times true and false appear
		isIvan, isNotIvan := averageTrueFalse(isIvanPossiblities)

		title := fmt.Sprintf("<@%v>", i.Member.User.ID)
		_, err = sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &title,
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title: fmt.Sprintf("Chances <@%v> is ivan are...", optionsMap["user"].UserValue(sess).Username),
					// "Chances this person is ivan are...",
					Description: fmt.Sprintf("**Is Ivan**: %v%%\n**Is not Ivan**: %v%%", isIvan, isNotIvan),
					Timestamp:   "",
					Color:       0,
					Fields:      []*discordgo.MessageEmbedField{},
				},
			},
		})
		if err != nil {
			log.Println(err)
		}
	}
}

// checkIsIvan checks if a message is Ivan
// returns true if it is an Ivan message
func checkIsIvan(message string) bool {
	c := cmd.NewCmd("python3", "train.py", message)
	c.Dir = "./ivan_detector/"
	outStatus := <-c.Start()
	out := outStatus.Stdout[0]

	if out == "True" {
		return true
	} else {
		return false
	}
}

// averageTrueFalse calculates the percentage of true and false in an array
// return percentage appearance of true, and false
func averageTrueFalse(possiblities []bool) (float64, float64) {
	length := float64(len(possiblities))
	trueCount, falseCount := 0.0, 0.0

	for _, v := range possiblities {
		if v == true {
			trueCount++
		} else {
			falseCount++
		}
	}

	return (trueCount / length) * 100, (falseCount / length) * 100
}
