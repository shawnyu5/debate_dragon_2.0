package reddit

import (
	"context"
	"log"
	"math/rand/v2"

	"github.com/bwmarrin/discordgo"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Reddit struct{}

// Components implements commands.Command.
func (Reddit) Components() []commands.Component {
	return []commands.Component{}
}

// Def implements commands.Command.
func (Reddit) Def() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Version:     "1.0.0",
		Name:        "reddit",
		Description: "Get a random Reddit post from r/Seneca",
	}
}

// Handler implements commands.Command.
func (Reddit) Handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		err := utils.DeferReply(sess, i.Interaction)
		if err != nil {
			return "", err
		}

		response, err := getRandomRedditPost()
		if err != nil {
			return "", err
		}
		_, err = sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &response.URL,
		})
		if err != nil {
			return "", err
		}
	}
	return "Send r/Seneca reddit post", nil
}

// getRandomRedditPost get a random post from r/Seneca
func getRandomRedditPost() (*reddit.Post, error) {
	client, err := reddit.NewReadonlyClient()
	if err != nil {
		log.Fatal(err)
	}

	posts, _, err := client.Subreddit.TopPosts(context.Background(), "seneca", &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 30,
		},
		Time: "all",
	})

	if err != nil {
		log.Fatal(err)
	}

	randomNumber := rand.IntN(30)
	return posts[randomNumber], nil

}
