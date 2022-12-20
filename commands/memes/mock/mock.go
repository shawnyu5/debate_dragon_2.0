package mock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
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
	getMeme()
	return "mock", nil
}

func getMeme() {
	// TODO: change this
	type Body struct {
		TemplateID string `json:"template_id"`
		UserName   string `json:"username"`
		Password   string `json:"password"`
		Text0      string `json:"text0"`
		Text1      string `json:"text1"`
	}
	// body := Body{
	// TemplateID: "102156234",
	// UserName:   "shawnyu",
	// Password:   "vvfx!v6tbnn2b!q",
	// Text0:      "mock",
	// Text1:      "world",
	// }

	// body := []byte(`{
	// "template_id": "102156234",
	// "username": "ShawnYu",
	// "password": "Vvfx!v6TbNN2B!q",
	// "text0": "mock",
	// "text1": "world"
	// }`)
	payload := url.Values{
		"template_id": {"102156234"},
		"username":    {"ShawnYu"},
		"password":    {"Vvfx!v6TbNN2B!q"},
		"boxes[text]": {"hello", "world"},
	}
	// "text0":       {"mock"},
	// "text1":       {"world"},
	// j, err := json.Marshal(body)
	// if err != nil {
	// log.Fatal(err)
	// }
	// fmt.Printf("getMeme payload.Encode(): %v\n", payload.Get("text0")) // __AUTO_GENERATED_PRINT_VAR__
	req, _ := http.NewRequest(http.MethodPost, "https://api.imgflip.com/caption_image", strings.NewReader(payload.Encode()))
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 ")
	res, _ := http.DefaultClient.Do(req)
	b, _ := ioutil.ReadAll(res.Body)

	fmt.Printf("getMeme b: %v\n", string(b)) // __AUTO_GENERATED_PRINT_VAR__
}
