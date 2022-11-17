package poll

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger"
	"github.com/olekukonko/tablewriter"
	"github.com/shawnyu5/debate_dragon_2.0/commands"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

// interview container type
type VotesContainer = map[int][]discordgo.User

// label which the data is stored in in db
const DBLabel = "interviewVotes"

var CommandObj = commands.CommandStruct{
	Name:    "interview",
	Obj:     obj,
	Handler: handler,
	Components: []struct {
		ComponentID      string
		ComponentHandler commands.HandlerFunc
	}{},
}

func obj() *discordgo.ApplicationCommand {
	var minValue = float64(0)
	return &discordgo.ApplicationCommand{
		Version:                  "1.0.0",
		Name:                     "interview",
		Description:              "a poll to keep track of the number of people that received interviews",
		DescriptionLocalizations: &map[discordgo.Locale]string{},
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionInteger,
				Name:         "vote",
				Description:  "vote for how many interviews you have received so far",
				Required:     false,
				Options:      []*discordgo.ApplicationCommandOption{},
				Autocomplete: false,
				MinValue:     &minValue,
				MaxValue:     10,
			},
			{
				Type:         discordgo.ApplicationCommandOptionBoolean,
				Name:         "getvotes",
				Description:  "view all votes",
				ChannelTypes: []discordgo.ChannelType{},
				Required:     false,
				Options:      []*discordgo.ApplicationCommandOption{},
				Autocomplete: false,
				Choices:      []*discordgo.ApplicationCommandOptionChoice{},
			},
			{
				Type:         discordgo.ApplicationCommandOptionBoolean,
				Name:         "remove",
				Description:  "remove your vote",
				ChannelTypes: []discordgo.ChannelType{},
				Required:     false,
				Options:      []*discordgo.ApplicationCommandOption{},
				Autocomplete: false,
				Choices:      []*discordgo.ApplicationCommandOptionChoice{},
			},
		},
	}
}

func handler(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	db := openDB()
	defer db.Close()

	options := utils.ParseUserOptions(sess, i)
	// record a user vote
	if val, ok := options["vote"]; ok {
		vote := val.IntValue()
		votes := getVotesFromDB(db)

		mess := formatVotes(votes, "Recording your vote... It may take a few seconds")

		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: mess,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		votes, _ = removeUserVotes(votes, *i.Member.User)
		votes = addVote(votes, int(vote), *i.Member.User)
		err := saveVotes(db, votes)
		if err == badger.ErrConflict {
			panic(err)
		}

		time.Sleep(3 * time.Second)

		mess = formatVotes(votes, "Your vote has been recorded successfully")
		sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &mess,
		})
		return "User vote recorded", nil

	} else if val, ok := options["getvotes"]; ok {
		// return all user votes
		if val.BoolValue() {
			votes := getVotesFromDB(db)
			mess := formatVotes(votes, "")

			sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: mess,
				},
			})
		} else {
			sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ok",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

		}
		return "All votes returned", nil
	} else if val, ok := options["remove"]; ok {
		// remove a user's votes
		if val.BoolValue() {
			votes := getVotesFromDB(db)
			message := formatVotes(votes, "Removing your vote... This may take a few seconds")
			sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: message,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			votes, removed := removeUserVotes(votes, *i.Member.User)
			saveVotes(db, votes)

			time.Sleep(3 * time.Second)

			// if a user's vote was removed, tell them their vote has been removed
			if removed {
				message = formatVotes(votes, "Your vote has been removed")
				sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &message,
				})
			} else {
				message = formatVotes(votes, "You have not voted... Pls vote")
				sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &message,
				})
			}
		} else {
			// if not args is passed in
			sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ok",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
		return "User vote removed", nil
	} else {
		// default response, command with no args
		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Components: []discordgo.MessageComponent{},
				Embeds: []*discordgo.MessageEmbed{
					{
						Type:  "",
						Title: "/interview",
						Description: `A poll like command to store votes internally. Allowing people to vote for how many interviews they've received for co-op so far

                  ` +
							"`/interview vote:<number>` cast your vote\n`/interview getvotes:true` gets all votes\n`/interview remove:true` removes your vote",
						Timestamp: "",
						Color:     0,
						Footer: &discordgo.MessageEmbedFooter{
							Text: "Telemetry data sent to <@799157783307092008>",
						},
					},
				},
				AllowedMentions: &discordgo.MessageAllowedMentions{},
				Files:           []*discordgo.File{},
				Flags:           0,
				Choices:         []*discordgo.ApplicationCommandOptionChoice{},
				CustomID:        "",
				Title:           "",
			},
		})
		return "Default response", nil
	}
}

// openDB opens a connection to the local database
func openDB() *badger.DB {
	// openDB opens a connection to the local database
	opts := badger.DefaultOptions("./db")
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	return db
}

// getVotesFromDB retrieve the votes from local databse
// db: the db to retrieve from
func getVotesFromDB(db *badger.DB) VotesContainer {
	// container for votes
	votes := make(VotesContainer)
	db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(DBLabel))
		if err != nil {
			return err
		}
		// return nil
		err = item.Value(func(val []byte) error {
			err := json.Unmarshal(val, &votes)
			if err != nil {
				return err
			}
			return err
		})
		return err
	})

	// if err != badger.ErrKeyNotFound {
	// panic(err)
	// }
	return votes
}

// addVote add a new vote to the existing collection of votes, and return the votes container
// votes  : a map of existing votes
// newVote: the new vote to add, in numeric value
// user   : the user who voted
// returns: the updated votes container
func addVote(votes VotesContainer, newVote int, user discordgo.User) VotesContainer {
	votes[newVote] = append(votes[newVote], user)
	return votes
}

// removeUserVotes removes all of a user's votes
// votes  : the votes container
// user   : the user's votes to remove
// returns: the updated votes container
func removeUserVotes(votes VotesContainer, user discordgo.User) (VotesContainer, bool) {
	filteredContainer := make(VotesContainer)
	// if we encountered the user we are removing
	found := false
	for voteCount, vote := range votes {
		for _, u := range vote {
			// add the users that are not the user we are removing
			if u.ID != user.ID {
				filteredContainer[voteCount] = append(filteredContainer[voteCount], u)
			} else {
				// keep track of we've encountered the user we are removing
				found = true
			}
		}
	}
	return filteredContainer, found
}

func removeIndex(s []discordgo.User, index int) []discordgo.User {
	return append(s[:index], s[index+1:]...)
}

// saveVotes save the db to local disk
// db     : the db to save
// votes  : the votes to save
// returns: error if any
func saveVotes(db *badger.DB, votes VotesContainer) error {
	err := db.Update(func(txn *badger.Txn) error {
		j, err := json.Marshal(votes)
		if err != nil {
			return err
		}
		entry := badger.NewEntry([]byte(DBLabel), j)
		err = txn.SetEntry(entry)
		return err
	})
	return err
}

// formatVotes formats a container of votes as a discord message
// votes  : the votes to format
// message: a message to append to the end of the vote tally
// retturn: the formatted message
func formatVotes(votes VotesContainer, message string) string {
	data := make([][]string, 0)
	mess := &strings.Builder{}
	mess.WriteString("```")
	for idx, vote := range votes {
		// string representation of the index
		strIdx := strconv.Itoa(idx)
		// string representation of the number of votes
		strLen := strconv.Itoa(len(vote))
		data = append(data, []string{strIdx, strLen})
	}

	table := tablewriter.NewWriter(mess)
	table.SetHeader([]string{"NO.Interviews", "Votes"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
	mess.WriteString("```")
	mess.WriteString("\n\n" + message)
	return mess.String()
}
