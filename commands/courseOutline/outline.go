package courseoutline

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly"
	"github.com/shawnyu5/debate_dragon_2.0/command"
	"github.com/shawnyu5/debate_dragon_2.0/utils"
)

var outline = command.Command{
	Name: "outline",
	ApplicationCommand: func() *discordgo.ApplicationCommand {
		return &discordgo.ApplicationCommand{
			Version:     "1.0.1",
			Type:        0,
			Name:        "outline",
			Description: "Find a course outline",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "course_code",
					Description:  "course code",
					Required:     true,
					Autocomplete: true,
				},
			},
		}
	},
	EditInteractionResponse: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		utils.DeferReply(sess, i.Interaction)
		url := GeneratewebPageURL(sess, i)
		courseInfo := GetCourseInfo(url)
		courseInfo.URL = url
		sess.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content:    new(string),
			Components: &[]discordgo.MessageComponent{},
			Embeds: &[]*discordgo.MessageEmbed{
				{
					URL:         url,
					Type:        discordgo.EmbedTypeArticle,
					Title:       fmt.Sprintf("**%s**", courseInfo.Title),
					Description: courseInfo.Description,
					Timestamp:   "",
					Color:       0,
				},
			},
			Files:           []*discordgo.File{},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		})

		return fmt.Sprintf("Send course outline for `%s`", courseInfo.Title), nil
	},
	InteractionApplicationCommandAutocomplete: func(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
		return GenerateSubjectCodeCompletion(sess, i)
	},
}

type CourseInfo struct {
	// title of course
	Title string
	// course description
	Description string
	URL         string
}

// GenerateSubjectCodeCompletion get the subject codes to fill in for autocompletion based on current user input
func GenerateSubjectCodeCompletion(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	input := utils.ParseUserOptions(sess, i)
	type findSubjectCodesCompletion struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}

	req, err := http.NewRequest(http.MethodGet, "https://www.senecacollege.ca/ssos/jsp/ajaxCalls/findSubjectCodes.jsp?isLoggedIn=", nil)
	// req, err := http.NewRequest(http.MethodGet, "https://www.senecacollege.ca/ssos/jsp/ajaxCalls/findSubjectCodes.jsp?isLoggedIn=&term=ECN230", nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("term", input["course_code"].StringValue())
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	var subjectCodes []findSubjectCodesCompletion
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	json.Unmarshal(b, &subjectCodes)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: []*discordgo.ApplicationCommandOptionChoice{},
		},
	}
	for _, code := range subjectCodes {
		response.Data.Choices = append(response.Data.Choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  code.Label,
			Value: code.Value,
		})
	}
	sess.InteractionRespond(i.Interaction, response)
	return "Filled out subject code for Autocomplete", nil

}

// GeneratewebPageURL constructs a url that leads to the course outline page for the selected course.
// return: url to the webpage containing course outline
func GeneratewebPageURL(sess *discordgo.Session, i *discordgo.InteractionCreate) string {
	input := utils.ParseUserOptions(sess, i)
	// https://www.senecacollege.ca/ssos/findOutline.do?isLoggedIn=&subjectOrAndTitle=%5BECN230%5D+Making+Sense+of+Our+Economy&schoolCode=
	req, err := http.NewRequest(http.MethodGet, "https://www.senecacollege.ca/ssos/findOutline.do?", nil)
	if err != nil {
		log.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("subjectOrAndTitle", input["course_code"].StringValue())
	req.URL.RawQuery = q.Encode()
	return req.URL.String()

}

// GetCourseInfo scrapes the url passed in to get course information
// url: the course outline to scrape.
// return: CourseInfo struct containing the title and description of the course.
func GetCourseInfo(url string) CourseInfo {
	c := colly.NewCollector()
	info := CourseInfo{}

	idx := 0
	c.OnHTML("div.sectionModification", func(e *colly.HTMLElement) {
		e.ForEach("div", func(i int, el *colly.HTMLElement) {
			switch idx {
			case 1:
				info.Description = strings.ReplaceAll(el.Text, "Subject Description", "")
				info.Description = strings.TrimSpace(info.Description)
				idx++
			default:
				idx++
			}
		})
	})

	// the only h1 in the page is the title
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		info.Title = strings.TrimSpace(e.Text)

	})

	c.Visit(url)
	return info
}

// GetStringInBetween Returns empty string if no start string found
func GetStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	e += s + e - 1
	return str[s:e]
}

// TODO: unable to get school information. Skipping for now
func CreateSchoolMenu(sess *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	// response from the api for schools that teaches the subject
	type SchoolCompletion struct {
		SName string `json:"sName"`
		Scode string `json:"scode"`
	}
	input := utils.ParseUserOptions(sess, i)
	// https://www.senecacollege.ca/ssos/jsp/ajaxCalls/pullSchoolsThatTeachSubject.jsp
	req, err := http.NewRequest(http.MethodPost, "https://www.senecacollege.ca/ssos/jsp/ajaxCalls/pullSchoolsThatTeachSubject.jsp", nil)
	// req, err := http.NewRequest(http.MethodGet, "https://www.senecacollege.ca/ssos/findOutline.do", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json, text/javascript, */*; q=0.01")

	q := req.URL.Query()
	q.Add("subjectTitleCode", input["course_code"].StringValue())
	// q.Add("subjectOrAndTitle", input["course_code"].StringValue())
	req.URL.RawQuery = q.Encode()
	// %5BECN230%5D+Making+Sense+of+Our+Economy

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	var schools []SchoolCompletion
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	// fmt.Printf("CreateSchoolMenu b: %v\n", string(b)) // __AUTO_GENERATED_PRINT_VAR__

	json.Unmarshal(b, &schools)
	fmt.Printf("CreateSchoolMenu schools: %v\n", schools) // __AUTO_GENERATED_PRINT_VAR__
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			TTS:             false,
			Content:         "",
			Components:      []discordgo.MessageComponent{},
			Embeds:          []*discordgo.MessageEmbed{},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
			Files:           []*discordgo.File{},
			Flags:           0,
			Choices:         []*discordgo.ApplicationCommandOptionChoice{},
			CustomID:        "",
			Title:           "",
		},
	}
	for _, school := range schools {
		response.Data.Choices = append(response.Data.Choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  school.SName,
			Value: school.Scode,
		})
	}
	// sess.InteractionRespond(i.Interaction, response)
	return "Filled out school for Autocomplete", nil
}

func init() {
	command.Register(outline)
}
